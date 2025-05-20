# 📊 Expense Tracker Telegram Bot

Это простой Telegram-бот для отслеживания личных доходов и расходов. Написан на Go, использует PostgreSQL и Docker.

## 🚀 Функционал

- Добавление расходов через команду `/расход`
- Подключение к PostgreSQL через Docker
- Структурированный Go-проект по best practices
- Поддержка миграций базы данных

## 📦 Стек технологий

- Go
- PostgreSQL
- Docker + docker-compose
- Telegram Bot API
- SQL миграции (`migrate/migrate`)

## 🛠 Установка

1. Клонируй репозиторий:

```bash
git clone https://github.com/yourname/expense-tracker-bot.git
cd expense-tracker-bot