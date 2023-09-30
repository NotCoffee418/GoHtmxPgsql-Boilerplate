package db_access

import (
	"github.com/NotCoffee418/GoWebsite-Boilerplate/database"
	"github.com/jmoiron/sqlx"
	"time"
)

type GuestBookEntry struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type GuestBookEntryRepository struct{}

func (*GuestBookEntryRepository) GetRecent(db *sqlx.DB, amount int) chan database.QueryResults[GuestBookEntry] {
	return database.ExecuteQuery[GuestBookEntry](db,
		"SELECT * FROM guestbook ORDER BY created_at DESC LIMIT ?",
		amount)
}

func (*GuestBookEntryRepository) Insert(db *sqlx.DB, name string, message string) chan database.NonQueryResult {
	return database.ExecuteNonQuery(db,
		"INSERT INTO guestbook (name, message) VALUES (?, ?)",
		name, message)
}
