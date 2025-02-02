package database

import (
	"fmt"
	"log"
	"os"

	"github.com/group-project/aitu-fan/aitu-fan/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database using GORM:", err)
	}
	DB = gormDB
	fmt.Println("PostgreSQL connected successfully with GORM!")

	err = DB.AutoMigrate(&models.Event{})
	if err != nil {
		log.Fatal("Error migrating the Event model:", err)
	}
	fmt.Println("Database migration completed: Event model created/updated.")

	err = DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatal("Error migrating models:", err)
	}
	fmt.Println("Database migration completed: User and Post models created/updated.")

}

func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error getting underlying SQL DB:", err)
	}
	sqlDB.Close()
	fmt.Println("Database connection closed.")
}
