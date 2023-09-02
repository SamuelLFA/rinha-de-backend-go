package main

import (
	"github.com/SamuelLFA/rinha-de-backend-go/internal/api"
	"github.com/SamuelLFA/rinha-de-backend-go/pkg/database"
)

func main() {
	// Initialize the database
	db := database.InitDatabase()

	// Initialize the application
	api.Init(db)
}
