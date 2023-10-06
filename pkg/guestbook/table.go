package guestbook

import "time"

type Row struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
