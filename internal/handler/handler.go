package handler

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Raxhet/expense-tracker-bot/internal/model"
	"github.com/Raxhet/expense-tracker-bot/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	storage *storage.Storage
}

func NewHandler(bot *tgbotapi.BotAPI, storage *storage.Storage) *Handler {
	return &Handler{
		bot:     bot,
		storage: storage,
	}
}

func (h *Handler) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	update := h.bot.GetUpdatesChan(u)
	for update := range update {
		if update.Message != nil {
			h.handleMessage(update.Message)
		}
	}
}

// реагирует просто на start, можно добавить switch/case для обработки необходимых команд
func (h *Handler) handleMessage(msg *tgbotapi.Message) {
	switch {
	case strings.HasPrefix(msg.Text, "/расход"):
		h.handleExpenseCommand(msg)
	default:
		h.sendMessage(msg.Chat.ID, "Напиши /расход <сумма> <категория>")
	}
}

func (h *Handler) handleExpenseCommand(msg *tgbotapi.Message) {
	args := strings.Fields(msg.Text)
	log.Printf("ARGS: %v", args) // logs
	if len(args) < 3 {
		h.sendMessage(msg.Chat.ID, "Неверный формат. Пример: /расход 100 еда")
		return
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	log.Printf("Amount: %v", amount) // logs
	if err != nil {
		h.sendMessage(msg.Chat.ID, "Сумма должна быть числом")
		return
	}

	category := strings.Join(args[2:], " ")
	tx := model.Transaction{
		UserID:   msg.Chat.ID,
		Amount:   amount,
		Category: category,
		Type:     model.Expense,
	}

	err = h.storage.AddTransaction(context.Background(), tx)

	if err != nil {
		h.sendMessage(msg.Chat.ID, "Ошибка добавления :(")
		log.Printf("Ошибка записи: %v", err)
		return
	}

	h.sendMessage(msg.Chat.ID, fmt.Sprintf("Расход %.2f RUB в категории '%s' сохранен!", amount, category))

}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.bot.Send(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}
