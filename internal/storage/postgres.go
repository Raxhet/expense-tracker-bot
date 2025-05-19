package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/Raxhet/expense-tracker-bot/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgresDB(cfg *config.Config) *pgx.Conn {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	log.Print("DSN: ", dsn)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to database")
	return conn
}
