package database

import (
	"log"

	"github.com/group-project/aitu-fan/aitu-fan/models"
)

func Migrate() {
	if DB == nil {
		log.Fatal(" Database not initialized. Call Connect() first.")
	}

	err := DB.AutoMigrate(&models.Comment{})
	if err != nil {
		log.Fatal(" Migration error:", err)
	}

	log.Println(" Migration successful")
}
