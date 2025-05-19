package main

import (
	_ "context"
	_ "fmt"
	"github.com/Raxhet/expense-tracker-bot/internal/config"
	"github.com/Raxhet/expense-tracker-bot/internal/handler"
	"github.com/Raxhet/expense-tracker-bot/internal/storage"
	"log"
	"log/slog"
	"os"
	_ "os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.LoadConfig()

	db := storage.NewPostgresDB(cfg)
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		slog.Error("Failed to create bot api", err)
	}
	bot.Debug = true
	log.Printf("Authorized on account @%s", bot.Self.UserName)
	h := handler.NewHandler(bot, db)
	h.StartPolling()
}
