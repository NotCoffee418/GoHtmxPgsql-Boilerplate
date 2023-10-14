package guestbook

type MessageDisplay struct {
	Name      string
	Message   string
	CreatedAt string
}

func NewMessageDisplay(row Row) MessageDisplay {
	return MessageDisplay{
		Name:      row.Name,
		Message:   row.Message,
		CreatedAt: row.CreatedAt.Format("02-01-2006 15:04:05"),
	}
}

func NewMessageDisplaySlice(rows []Row) []MessageDisplay {
	var displays []MessageDisplay
	for _, row := range rows {
		displays = append(displays, NewMessageDisplay(row))
	}
	return displays
}
