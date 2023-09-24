package dbmodels

type GuestBookEntry struct {
	ID      uint   `db:"id"`
	Name    string `db:"name"`
	Message string `db:"message"`
}
