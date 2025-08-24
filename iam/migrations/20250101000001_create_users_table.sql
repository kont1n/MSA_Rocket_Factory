-- +goose Up
-- Создание таблицы пользователей
CREATE TABLE users (
    user_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индексов для быстрого поиска
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_email ON users(email);

-- +goose Down
-- Удаление индексов
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_login;

-- Удаление таблицы пользователей
DROP TABLE IF EXISTS users;
