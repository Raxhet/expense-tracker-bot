# Dockerfile

FROM golang:1.24.3-alpine

# Устанавливаем утилиты
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o bot ./cmd/bot

# Запускаем
CMD ["./bot"]