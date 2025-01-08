package main

import (
	InitRouters "FileStorage/cmd/api/init"
	conf "FileStorage/config"
	"FileStorage/database/connection"
	"FileStorage/database/migrations"
	fileModels "FileStorage/database/models/file"
	tokenModels "FileStorage/database/models/token"
	userModels "FileStorage/database/models/user"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	config, err := conf.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	db, err := connection.Connect(config)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	migrations.Migrate(db, &userModels.User{}, &tokenModels.UserToken{}, &fileModels.File{})

	r := gin.Default()
	InitRouters.Run(r, db, config)

	r.Run(":8080")
	log.Printf("server started at localhost:8080")
}
