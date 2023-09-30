package database

import (
	"bufio"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/utils"
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

type MigrationState struct {
	AvailableVersion int
	InstalledVersion int
	Migrations       []migrationFileInfo
}

// MigrateUp migrates the database up to the latest version
func MigrateUp(db *sqlx.DB) {
	// Get migration state
	getLiveMigrationStateChan := make(chan MigrationState)
	go GetLiveMigrationInfo(db, getLiveMigrationStateChan)
	migrationState := <-getLiveMigrationStateChan
	close(getLiveMigrationStateChan)

	// Check if already up to date
	if migrationState.InstalledVersion == migrationState.AvailableVersion {
		log.Printf("Already up to date at version %d.\n", migrationState.InstalledVersion)
		return
	} else if migrationState.InstalledVersion > migrationState.AvailableVersion {
		log.Fatalf(
			"Installed migration version (%d) is higher than highest available migration (%d).",
			migrationState.InstalledVersion, migrationState.AvailableVersion)
	} else {
		log.Printf("Migrating from %d to %d...\n",
			migrationState.InstalledVersion, migrationState.AvailableVersion)
	}

	// Filter out new migrations to apply and grab their up/down contents
	var migrationsToApply []migrationFileInfo
	for _, migration := range migrationState.Migrations {
		if migration.version > migrationState.InstalledVersion {
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
	close(filledChannel)

	// Apply up migrations
	for _, migration := range migrationsToApply {
		// Init tx for this migration
		tx, err := db.Beginx()
		if err != nil {
			log.Fatalf("Error beginning transaction: %v", err)
		}

		// Run migration code
		log.Printf("Applying migration %d...\n", migration.version)
		_, err = tx.Exec(migration.contents.up)
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("Error applying migration (Exec) %d: %v", migration.version, err)
		}

		// Insert migration into migrations table
		_, err = tx.Exec(
			"INSERT INTO migrations (version, installed_at) VALUES ($1, $2)",
			migration.version, time.Now())
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("Error inserting migration version into migrations table %d: %v", migration.version, err)
		}

		// Commit tx
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Error committing migration %d: %v", migration.version, err)
		}
	}
	log.Println("Migration complete.")
}

// MigrateDown migrates the database down to the previous version
func MigrateDown(db *sqlx.DB) {
	// Get migration state
	getLiveMigrationStateChan := make(chan MigrationState)
	go GetLiveMigrationInfo(db, getLiveMigrationStateChan)
	liveState := <-getLiveMigrationStateChan
	close(getLiveMigrationStateChan)

	// Find index of previous migration
	prevMigrationIdx := -2
	for i, migration := range liveState.Migrations {
		if migration.version >= liveState.InstalledVersion {
			prevMigrationIdx = i - 1
			break
		}
	}
	migration := &liveState.Migrations[prevMigrationIdx]

	// Validation
	if prevMigrationIdx == -2 {
		log.Fatalf("Failed to find currently installed migration  %d", liveState.InstalledVersion)
	} else if prevMigrationIdx == -1 {
		log.Fatalf("No migrations to revert.")
	} else {
		log.Printf("Migrating from %d to %d...\n",
			liveState.InstalledVersion, liveState.Migrations[prevMigrationIdx].version)
	}

	// Get migration contents
	filledChannel := make(chan bool)
	fillMigrationContents(migration, filledChannel)
	<-filledChannel
	close(filledChannel)

	// Init tx for this migration
	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("Error beginning transaction: %v", err)
	}

	// Run migration code
	log.Printf("Applying migration %d...\n", migration.version)
	_, err = tx.Exec(migration.contents.down)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalf("Error applying migration (Exec) %d: %v", migration.version, err)
	}

	// Insert migration into migrations table
	_, err = tx.Exec(
		"DELETE FROM migrations WHERE version = $1", migration.version)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalf("Error removing version from migrations table %d: %v", migration.version, err)
	}

	// Commit tx
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error committing migration %d: %v", migration.version, err)
	}
}

// GetLiveMigrationInfo returns the latest migration version and the installed migration version
func GetLiveMigrationInfo(db *sqlx.DB, resultChan chan MigrationState) {
	log.Println("Getting migration info...")

	// Get installed migration version
	installedMigrationChan := make(chan int)
	go getInstalledMigrationVersion(db, installedMigrationChan)

	// Prepare channels for getting migration info
	allMigrationsChan := make(chan []migrationFileInfo)
	go listAllMigrations(allMigrationsChan)

	// Local migration info
	allMigrations := <-allMigrationsChan
	close(allMigrationsChan)
	totalMigrationCount := len(allMigrations)

	// Installed migration info
	installedMigration := <-installedMigrationChan
	close(installedMigrationChan)

	// Return
	if totalMigrationCount == 0 {
		log.Println("No migrations found.")
		resultChan <- MigrationState{
			AvailableVersion: 0,
			InstalledVersion: installedMigration,
			Migrations:       nil,
		}
		return
	} else {
		highestAvailableMigration := allMigrations[totalMigrationCount-1]
		resultChan <- MigrationState{
			AvailableVersion: highestAvailableMigration.version,
			InstalledVersion: installedMigration,
			Migrations:       allMigrations,
		}
	}

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
			log.Fatalf("Duplicate migration version: %d", version)
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

func getInstalledMigrationVersion(db *sqlx.DB, resultChan chan int) {
	// Ensure migrations table exists
	ensureMigrationsTableChan := make(chan bool)
	go EnsureMigrationTableExists(db, ensureMigrationsTableChan)
	<-ensureMigrationsTableChan
	close(ensureMigrationsTableChan)

	// Get installed migration version
	var version int
	err := db.Get(&version, "SELECT version FROM migrations ORDER BY version DESC LIMIT 1")
	if err != nil {
		// No migrations applied yet
		if errors.Is(err, sql.ErrNoRows) {
			resultChan <- 0
			return
		} else {
			log.Fatalf("Error getting migration version: %v", err)
		}
	}
	resultChan <- version
}

func fillMigrationContents(migration *migrationFileInfo, doneChan chan bool) {
	upRx := regexp.MustCompile(`(?i)--s*\+up(\s*)?(.+)?`)      // +up
	downRx := regexp.MustCompile(`(?i)--\s*\+down(\s*)?(.+)?`) // +down

	// Read file contents
	file, err := os.Open(migration.file)
	if err != nil {
		log.Fatalf("Error opening migration file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing migration file: %v", err)
		}
	}(file)

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
	doneChan <- true
}

// EnsureMigrationTableExists checks if the migrations table exists and creates it if it doesn't
func EnsureMigrationTableExists(db *sqlx.DB, doneChan chan bool) {
	// Exist check
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'migrations')")
	if err != nil {
		log.Fatalf("Error checking if migrations table exists: %v", err)
	}

	// Create on missing
	if !exists {
		_, err := db.Exec(`
		CREATE TABLE migrations (
			version INT NOT NULL,
			installed_at TIMESTAMP NOT NULL
			)
		`)
		if err != nil {
			log.Fatalf("Error creating migrations table: %v", err)
		}

	}
	doneChan <- true
}
