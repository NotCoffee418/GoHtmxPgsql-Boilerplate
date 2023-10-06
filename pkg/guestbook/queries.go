package guestbook

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Repository struct{}

func (r *Repository) GetRecent(db *sqlx.DB, amount int) ([]Row, error) {
	var entries []Row
	err := db.Select(&entries, "SELECT * FROM guestbook ORDER BY created_at DESC LIMIT ?", amount)
	if err != nil {
		log.Printf("Error getting recent guestbook entries: %s", err)
		return nil, err
	}
	return entries, nil
}

func (r *Repository) Insert(db *sqlx.DB, name string, message string) error {
	_, err := db.Exec("INSERT INTO guestbook (name, message) VALUES (?, ?)", name, message)
	if err != nil {
		log.Printf("Error inserting into guestbook: %s", err)
		return err
	}
	return nil
}
