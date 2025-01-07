package connection

import (
	conf "FileStorage/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Database struct {
	DB *gorm.DB
}

func getConnStr(config conf.Config) string {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	return connStr
}

func Connect() (*gorm.DB, error) {
	config, err := conf.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	connStr := getConnStr(*config)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Printf("connected to database: %s", config.Database.Name)

	return db, nil
}
