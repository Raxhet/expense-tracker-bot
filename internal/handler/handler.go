package handler

import (
	"context"
	"fmt"
	"github.com/Raxhet/expense-tracker-bot/internal/state"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"

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

func (h *Handler) handleMessage(msg *tgbotapi.Message) {
	s := state.GetSession(msg.Chat.ID)
	log.Printf("Chat ID: %v", msg.Chat.ID)
	switch { // TODO: доход
	case strings.HasPrefix(msg.Text, "/расход"):
		s.Step = state.AwaitingAmount
		s.Type = model.Expense
		h.sendMessage(msg.Chat.ID, "Введите сумму расхода: ")
	case s.Step == state.AwaitingAmount:
		h.handleStepAmount(msg)
	case s.Step == state.AwaitingCategory:
		h.handleStepCategory(msg)
	default:
		h.sendMessage(msg.Chat.ID, "Напиши /расход или /доход")
	}
}

func (h *Handler) handleStepAmount(msg *tgbotapi.Message) {
	s := state.GetSession(msg.Chat.ID)
	amount, err := strconv.ParseFloat(msg.Text, 64)
	if err != nil {
		h.sendMessage(msg.Chat.ID, "ЕБЛАН, ВВЕДИ ЧИСЛО")
		return
	}
	s.TempAmount = amount
	s.Step = state.AwaitingCategory

	h.showCategoryKeyboard(msg.Chat.ID) // test
}

func (h *Handler) handleStepCategory(msg *tgbotapi.Message) {
	s := state.GetSession(msg.Chat.ID)
	_ = s
	log.Printf("GOOD: %+v", s)
	s.Step = state.Idle
}

func (h *Handler) handleExpenseCommand(msg *tgbotapi.Message) {
	args := strings.Fields(msg.Text)
	log.Printf("ARGS: %v", args) // logs

	if len(args) == 1 {

	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = h.storage.AddTransaction(ctx, tx)

	if err != nil {
		h.sendMessage(msg.Chat.ID, "Ошибка добавления :(")
		log.Printf("Ошибка записи: %v", err)
		return
	}

	h.sendMessage(msg.Chat.ID, fmt.Sprintf("Расход %.2f RUB в категории '%s' сохранен!", amount, category))

}

func (h *Handler) handleCategoryCommand(msg *tgbotapi.Message) {
	args := strings.Fields(msg.Text)
	log.Printf("ARGS category: %v", args)
	if len(args) < 2 {
		h.sendMessage(msg.Chat.ID, "Формат: /категория <название>")
		return
	}
	name := strings.Join(args[1:], " ")
	ctg := model.Category{
		UserID: msg.Chat.ID,
		Name:   name,
	}
	err := h.storage.AddCategory(context.Background(), ctg)
	if err != nil {
		h.sendMessage(msg.Chat.ID, "Не удалось добавить категорию")
		log.Printf("Ошибка сохранения категории: %v", err)
		return
	}

	h.sendMessage(msg.Chat.ID, fmt.Sprintf("Категория '%s' добавлена!", name))
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := h.bot.Send(msg)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (h *Handler) showCategoryKeyboard(chatID int64) {
	ctg, err := h.storage.GetUserCategories(context.Background(), chatID)
	log.Printf("CTG: %+v, CHAT ID: %v", ctg, chatID)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	if len(ctg) == 0 {
		h.sendMessage(chatID, "У вас пока нет категорий. Добавьте командой /команда <название>")
		return
	}

	var rows [][]tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton

	for i, category := range ctg {
		btn := tgbotapi.NewKeyboardButton(category.Name)
		row = append(row, btn)

		if (i+1)%2 == 0 {
			rows = append(rows, row)
		}
	}

	if len(row) > 0 {
		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg := tgbotapi.NewMessage(chatID, "Выберите категорию и напишите сумму")
	msg.ReplyMarkup = keyboard

	h.bot.Send(msg)
}
