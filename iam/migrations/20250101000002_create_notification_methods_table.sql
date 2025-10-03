-- +goose Up
-- Создание таблицы способов уведомления
CREATE TABLE notification_methods (
    id BIGSERIAL PRIMARY KEY,
    user_uuid UUID NOT NULL REFERENCES users(user_uuid) ON DELETE CASCADE,
    provider_name VARCHAR(50) NOT NULL, -- telegram, email, push и т.д.
    target VARCHAR(255) NOT NULL, -- email адрес, telegram chat id и т.д.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индексов
CREATE INDEX idx_notification_methods_user_uuid ON notification_methods(user_uuid);
CREATE INDEX idx_notification_methods_provider ON notification_methods(provider_name);

-- +goose Down
-- Удаление индексов
DROP INDEX IF EXISTS idx_notification_methods_provider;
DROP INDEX IF EXISTS idx_notification_methods_user_uuid;

-- Удаление таблицы способов уведомления
DROP TABLE IF EXISTS notification_methods;
