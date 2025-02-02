package database

import (
	"github.com/group-project/aitu-fan/aitu-fan/models"
)

func GetAllEvents() ([]models.Event, error) {
	rows, err := DB.Raw("SELECT id, title, description, date FROM events").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Date); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func CreateEvent(event models.Event) error {
	result := DB.Create(&event)
	return result.Error
}
