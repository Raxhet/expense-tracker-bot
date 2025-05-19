package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/Raxhet/expense-tracker-bot/internal/config"
	"github.com/Raxhet/expense-tracker-bot/internal/model"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewPostgresDB(cfg *config.Config) *Storage {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to database")
	return &Storage{conn: conn}
}

func (s *Storage) Close() {
	s.conn.Close(context.Background())
}

func (s *Storage) AddTransaction(ctx context.Context, tx model.Transaction) error {
	query := `
		INSERT INTO transactions (user_id, amount, category, type)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.conn.Exec(ctx, query, tx.UserID, tx.Amount, tx.Category, tx.Type)
	return err
}
