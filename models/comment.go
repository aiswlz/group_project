package models

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	Author    string `gorm:"size:100;not null"`
	CreatedAt time.Time
	PostID    *uint `gorm:"index"`
	EventID   *uint `gorm:"index"`
}
