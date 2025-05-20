package state

import (
	"github.com/Raxhet/expense-tracker-bot/internal/model"
	"log"
)

type Step string

const (
	Idle             Step = "idle"
	AwaitingAmount   Step = "awaiting_amount"
	AwaitingCategory Step = "awaiting_category"
)

type UserSession struct {
	Step       Step
	Type       model.TransactionType
	TempAmount float64
}

var sessions = make(map[int64]*UserSession)

func GetSession(chatID int64) *UserSession {
	if session, ok := sessions[chatID]; ok {
		return session
	}

	session := &UserSession{Step: Idle}
	sessions[chatID] = session
	log.Printf("Sessions: %v", sessions)
	return session
}

func Reset(chatID int64) {
	delete(sessions, chatID)
}
