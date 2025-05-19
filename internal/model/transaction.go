package model

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID        int
	UserID    int64
	Amount    float64
	Category  string
	Type      TransactionType
	CreatedAt time.Time
}
