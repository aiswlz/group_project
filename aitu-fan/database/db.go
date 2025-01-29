package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	sqlDB *sql.DB
)

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

	// Initialize sql.DB connection
	sqlDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database using sql.DB:", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Database is not reachable:", err)
	}

	fmt.Println("PostgreSQL connected successfully with sql.DB!")
}

func CloseDatabase() {
	if sqlDB != nil {
		sqlDB.Close()
		fmt.Println("Database connection closed.")
	}
}
