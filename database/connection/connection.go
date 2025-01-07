package connection

import (
	conf "FileStorage/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

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

func Connect() (*pgx.Conn, error) {
	config, err := conf.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	connStr := getConnStr(*config)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Printf("connected to database: %s", config.Database.Name)

	return conn, nil
}
