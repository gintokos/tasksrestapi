FROM golang:1.23.2 AS builder

# Устанавливаем рабочую директорию для сборки
WORKDIR /tasks

# Копируем только go.mod и go.sum для того, чтобы скачать зависимости
COPY go.mod go.sum ./

# Скачиваем все зависимости
RUN go mod tidy

# Копируем весь исходный код, включая папку cmd
COPY . .

# Создаем папку для бинарников и компилируем приложение в нее
RUN mkdir -p /tasks/bin && go build -o /tasks/bin/main ./cmd

# Используем более новую версию образа Debian для запуска
FROM debian:bookworm-slim

# Устанавливаем SQLite в образ
RUN apt-get update && apt-get install -y sqlite3

# Копируем скомпилированное приложение из builder в папку bin
COPY --from=builder /tasks/bin/main /main

# Копируем файл config.json в контейнер
COPY config.json /config.json

# Указываем команду для запуска приложения
CMD ["/main"]
