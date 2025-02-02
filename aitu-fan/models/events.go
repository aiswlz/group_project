package models

import (
	_ "gorm.io/gorm"
	"time"
)

type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
