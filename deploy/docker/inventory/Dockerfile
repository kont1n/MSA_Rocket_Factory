# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.24-alpine AS builder

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы для кеширования зависимостей
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY shared/go.mod shared/go.sum ./shared/
COPY platform/go.mod platform/go.sum ./platform/

# Загружаем зависимости
WORKDIR /app/inventory
RUN go mod download

# Возвращаемся в корневую директорию и копируем исходный код
WORKDIR /app
COPY inventory/ ./inventory/
COPY shared/ ./shared/
COPY platform/ ./platform/

# Собираем приложение с оптимизациями
WORKDIR /app/inventory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o inventory-service ./cmd/main.go

# Финальный минимальный образ
FROM scratch

# Копируем SSL сертификаты для HTTPS запросов
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Создаем пользователя без прав root для безопасности
COPY --from=builder /etc/passwd /etc/passwd

# Копируем бинарный файл
COPY --from=builder /app/inventory/inventory-service /inventory-service

# Открываем порт для gRPC
EXPOSE 50051

# Запускаем приложение
ENTRYPOINT ["/inventory-service"]