# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.24-alpine AS builder

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы для кеширования зависимостей
COPY payment/go.mod payment/go.sum ./payment/
COPY shared/go.mod shared/go.sum ./shared/
COPY platform/go.mod platform/go.sum ./platform/

# Загружаем зависимости
WORKDIR /app/payment
RUN go mod download

# Возвращаемся в корневую директорию и копируем исходный код
WORKDIR /app
COPY payment/ ./payment/
COPY shared/ ./shared/
COPY platform/ ./platform/

# Скачиваем grpc-health-probe — это утилита для проверки состояния gRPC-сервиса
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

# Делаем скачанный файл исполняемым
RUN chmod +x grpc-health-probe

# Собираем приложение с оптимизациями
WORKDIR /app/payment
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o payment-service ./cmd/main.go

# Финальный минимальный образ
FROM scratch

# Копируем SSL сертификаты для HTTPS запросов
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Создаем пользователя без прав root для безопасности
COPY --from=builder /etc/passwd /etc/passwd

# Копируем бинарный файл
COPY --from=builder /app/payment/payment-service /payment-service

# Открываем порт для gRPC
EXPOSE 50052

# Запускаем приложение
ENTRYPOINT ["/payment-service"]
