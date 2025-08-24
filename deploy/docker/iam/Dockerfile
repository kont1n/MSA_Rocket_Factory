# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.24-alpine AS builder

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы для кеширования зависимостей
COPY iam/go.mod iam/go.sum ./iam/
COPY shared/go.mod shared/go.sum ./shared/
COPY platform/go.mod platform/go.sum ./platform/

# Загружаем зависимости
WORKDIR /app/iam
RUN go mod download

# Возвращаемся в корневую директорию и копируем исходный код
WORKDIR /app
COPY iam/ ./iam/
COPY shared/ ./shared/
COPY platform/ ./platform/

# Скачиваем grpc-health-probe — это утилита для проверки состояния gRPC-сервиса
ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

# Делаем скачанный файл исполняемым
RUN chmod +x grpc-health-probe

# Собираем приложение с оптимизациями
WORKDIR /app/iam
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o iam-service ./cmd/main.go

# Финальный образ на основе Alpine для поддержки миграций
FROM alpine:3.20

# Устанавливаем ca-certificates для HTTPS запросов и создаем пользователя
RUN apk --no-cache add ca-certificates && \
    addgroup -g 1001 appgroup && \
    adduser -u 1001 -G appgroup -s /bin/sh -D appuser

# Создаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /app/iam/iam-service .

# Копируем миграции
COPY --from=builder /app/iam/migrations ./migrations

# Копируем grpc-health-probe для health checks
COPY --from=builder /app/grpc-health-probe .

# Изменяем владельца файлов
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт для gRPC
EXPOSE 50050

# Запускаем приложение
ENTRYPOINT ["./iam-service"]