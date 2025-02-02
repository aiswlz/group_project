package models

import "time"

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint   `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
