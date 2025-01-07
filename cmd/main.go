package main

import (
	"FileStorage/database/connection"
	"FileStorage/database/migrations"
	models "FileStorage/database/models/user"
	"log"
)

func main() {
	db, err := connection.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	migrations.Migrate(db, &models.User{})
}
