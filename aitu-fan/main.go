package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/group-project/aitu-fan/config"
	"github.com/group-project/aitu-fan/database"
	"github.com/group-project/aitu-fan/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	config.InitCloudinary()

	database.ConnectMongoDB()

	router := mux.NewRouter()

	routes.SetupRoutes(router)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
