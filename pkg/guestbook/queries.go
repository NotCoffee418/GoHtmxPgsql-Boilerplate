package guestbook

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Repository struct{}

func (r *Repository) GetRecent(db *sqlx.DB, amount int) ([]Row, error) {
	var entries []Row
	err := db.Select(&entries, "SELECT * FROM demo_guestbook ORDER BY created_at DESC LIMIT ?", amount)
	if err != nil {
		log.Errorf("Error getting recent guestbook entries: %s", err)
		return nil, err
	}
	return entries, nil
}

func (r *Repository) Insert(db *sqlx.DB, name string, message string) error {
	_, err := db.Exec("INSERT INTO demo_guestbook (name, message) VALUES ($1, $2)", name, message)
	if err != nil {
		log.Errorf("Error inserting into guestbook: %s", err)
		return err
	}
	return nil
}
