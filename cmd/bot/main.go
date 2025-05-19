package main

import (
	"context"
	"fmt"
	"github.com/Raxhet/expense-tracker-bot/internal/config"
	"github.com/Raxhet/expense-tracker-bot/internal/storage"
	"os"
)

func main() {
	fmt.Printf("Connecting to DB at %s:%s, user=%s db=%s\n",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"))
	cfg := config.LoadConfig()

	db := storage.NewPostgresDB(cfg)
	defer db.Close(context.Background())
}
