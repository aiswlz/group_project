package models

import (
	_ "gorm.io/gorm"
	"time"
)

type Event struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Date        time.Time `gorm:"type:date;not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   uint      `gorm:"autoCreateTime"`
	UpdatedAt   uint      `gorm:"autoUpdateTime"`
}
