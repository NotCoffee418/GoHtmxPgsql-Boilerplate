package db

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/server"
	"github.com/NotCoffee418/GoHtmxPgsql-Boilerplate/internal/utils"
	"github.com/jmoiron/sqlx"
)

type migrationFileInfo struct {
	version  int
	file     string
	contents *migrationContents // not always populated
}

type migrationContents struct {
	up   string
	down string
}

type MigrationsTable struct {
	Version     int       `db:"version"`
	InstalledAt time.Time `db:"installed_at"`
}

// MigrateUp migrates the database up to the latest version
func MigrateUp() {
	conn := server.GetDBConn()

	// Prepare channels for getting migration info
	log.Println("Getting migration info...")
	allMigrationsChan := make(chan []migrationFileInfo)
	go listAllMigrations(allMigrationsChan)

	installedMigrationChan := make(chan int)
	go getInstalledMigrationVersion(conn, installedMigrationChan)

	// Extract migration info
	allMigrations := <-allMigrationsChan
	totalMigrationCount := len(allMigrations)
	if totalMigrationCount == 0 {
		log.Println("No migrations found.")
		return
	}
	highestAvailableMigration := allMigrations[totalMigrationCount-1]
	installedMigrationVersion := <-installedMigrationChan

	// Check if already up to date
	if installedMigrationVersion == highestAvailableMigration.version {
		log.Printf("Already up to date at version %d.\n", installedMigrationVersion)
		return
	} else if installedMigrationVersion > highestAvailableMigration.version {
		log.Fatalf(
			"Installed migration version (%d) is higher than highest available migration (%d).",
			installedMigrationVersion, highestAvailableMigration.version)
	} else {
		log.Printf("Migrating from %d to %d...\n",
			installedMigrationVersion, highestAvailableMigration.version)
	}

	// Filter out new migrations to apply and grab their up/down contents
	migrationsToApply := []migrationFileInfo{}
	for i, migration := range allMigrations {
		if migration.version > installedMigrationVersion {
			migrationsToApply = append(migrationsToApply, migration)
		}
	}

	// fill up/down contents concurrently
	filledChannel := make(chan bool)
	for i := range migrationsToApply {
		idx := i
		fillMigrationContents(&migrationsToApply[idx], filledChannel)
	}
	for range migrationsToApply {
		<-filledChannel
	}

	// Apply up migrations
	for _, migration := range migrationsToApply {
		log.Printf("Applying migration %d...\n", migration.version)
		_, err := conn.Exec(migration.contents.up)
		if err != nil {
			log.Fatalf("Error applying migration %d: %v", migration.version, err)
		}
		_, err = conn.Exec(
			"INSERT INTO migrations (version, installed_at) VALUES ($1, $2)",
			migration.version, time.Now())
		if err != nil {
			log.Fatalf("Error inserting migration version %d: %v", migration.version, err)
		}
	}

}

// MigrateDown migrates the database down to the previous version
func MigrateDown() {
	//conn := server.GetDBConn()

}

func listAllMigrations(result chan []migrationFileInfo) {
	// List all valid migration files
	re := regexp.MustCompile(`.+[/|\\](\d{4})_\S+\.sql`)
	migrationFiles, err := utils.GetRecursiveFiles("./migrations", func(p string) bool {
		return re.FindStringSubmatch(p) != nil
	})
	if err != nil {
		log.Fatalf("Error reading migrations directory: %v", err)
	}

	// Create map of version per file path
	migrationMap := make(map[int]string)
	for _, file := range migrationFiles {
		matches := re.FindStringSubmatch(file)
		version, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatalf("Error parsing migration version: %v", err)
		}

		// Duplicate version check
		if _, exists := migrationMap[version]; exists {
			log.Fatalf("Duplicate migration version: %s", version)
		}
		migrationMap[version] = file
	}

	// Get sorted migration versions for iterating map
	sortedVersions := make([]int, 0, len(migrationFiles))
	for k := range migrationFiles {
		sortedVersions = append(sortedVersions, k)
	}
	sort.Ints(sortedVersions)

	// Return slice of sorted migrationFileInfo
	sortedMigrationFiles := make([]migrationFileInfo, 0, len(migrationFiles))
	for _, version := range sortedVersions {
		sortedMigrationFiles = append(sortedMigrationFiles, migrationFileInfo{
			version: version,
			file:    migrationMap[version],
		})
	}
	result <- sortedMigrationFiles
}

func getInstalledMigrationVersion(conn *sqlx.DB, result chan int) {
	var version int
	err := conn.Get(&version, "SELECT version FROM migrations ORDER BY version DESC LIMIT 1")
	if err != nil {
		if err == sql.ErrNoRows {
			result <- 0
		}
		log.Fatalf("Error getting migration version: %v", err)
	}
	result <- version
}

func fillMigrationContents(migration *migrationFileInfo, returnChan chan bool) {
	upRx := regexp.MustCompile(`(?i)\/\/\s*\+up(\s*)?(.+)?`)     // +up
	downRx := regexp.MustCompile(`(?i)\/\/\s*\+down(\s*)?(.+)?`) // +down

	// Read file contents
	file, err := os.Open(migration.file)
	if err != nil {
		log.Fatalf("Error opening migration file: %v", err)
	}
	defer file.Close()

	foundUp := false
	foundDown := false
	capturingSection := 0
	var upContents, downContents strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check for up/down section
		if upRx.MatchString(line) {
			if foundUp {
				log.Fatalf("Duplicate up section in migration %d", migration.version)
			}
			foundUp = true
			capturingSection = 1
			continue
		} else if downRx.MatchString(line) {
			if foundDown {
				log.Fatalf("Duplicate down section in migration %d", migration.version)
			}
			foundDown = true
			capturingSection = 2
			continue
		}

		// Capture up/down section contents
		if capturingSection == 1 {
			upContents.WriteString(line)
			upContents.WriteString("\n")
		} else if capturingSection == 2 {
			downContents.WriteString(line)
			downContents.WriteString("\n")
		}
	}

	// Validation
	if !foundUp {
		log.Fatalf("Missing `// +up` section in migration %d", migration.version)
	}
	if !foundDown {
		log.Fatalf("Missing `// +down` section in migration %d", migration.version)
	}

	// Return
	migration.contents = &migrationContents{
		up:   upContents.String(),
		down: downContents.String(),
	}
	returnChan <- true
}
