package main

import (
	"fmt"
	"github.com/group-project/aitu-fan/aitu-fan/routes"
	"log"
	"net/http"
)

func main() {
	r := routes.SetupRoutes()

	fmt.Println(" Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
