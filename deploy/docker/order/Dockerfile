# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.24-alpine AS builder

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы для кеширования зависимостей
COPY order/go.mod order/go.sum ./order/
COPY shared/go.mod shared/go.sum ./shared/
COPY platform/go.mod platform/go.sum ./platform/

# Загружаем зависимости
WORKDIR /app/order
RUN go mod download

# Возвращаемся в корневую директорию и копируем исходный код
WORKDIR /app
COPY order/ ./order/
COPY shared/ ./shared/
COPY platform/ ./platform/

# Собираем приложение с оптимизациями
WORKDIR /app/order
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o order-service ./cmd/main.go

# Финальный образ на основе Alpine для поддержки миграций
FROM alpine:3.20

# Устанавливаем ca-certificates для HTTPS запросов и создаем пользователя
RUN apk --no-cache add ca-certificates && \
    addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# Создаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /app/order/order-service .

# Копируем миграции
COPY --from=builder /app/order/migrations ./migrations

# Изменяем владельца файлов
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт для HTTP API
EXPOSE 8080

# Запускаем приложение
ENTRYPOINT ["./order-service"]