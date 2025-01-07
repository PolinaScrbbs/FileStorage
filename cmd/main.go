package main

import (
	"FileStorage/database/connection"
	"context"
	"log"
)

func main() {
	conn, err := connection.Connect()
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer conn.Close(context.Background())
}
