# Указываем базовый образ
FROM golang:1.23.6 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
# RUN go build -o news-app .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o news-app .

# Указываем образ для выполнения
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /usr/local/bin

# Устанавливаем необходимые пакеты (если необходимо)
RUN apk --no-cache add ca-certificates

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/news-app .

# Делаем исполняемый файл доступным для выполнения
RUN chmod +x news-app

# Копируем файл с переменными окружения
COPY config.env /usr/local/bin/config.env

COPY migrations /usr/local/bin/migrations

# Указываем команду для запуска приложения
CMD ["news-app"]
