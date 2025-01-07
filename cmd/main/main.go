package main

import (
	InitRouters "FileStorage/cmd/api/init"
	"FileStorage/database/connection"
	"FileStorage/database/migrations"
	models "FileStorage/database/models/user"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db, err := connection.Connect()
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	migrations.Migrate(db, &models.User{})

	r := gin.Default()
	InitRouters.Run(r, db)

	r.Run(":8080")
	log.Printf("server started at localhost:8080")
}
