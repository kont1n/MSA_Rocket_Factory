package model

import "errors"

var (
	// Ошибки пользователей
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInvalidUserData    = errors.New("invalid user data")
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrFailedToGetUser    = errors.New("failed to get user")
	ErrFailedToUpdateUser = errors.New("failed to update user")

	// Ошибки сессий
	ErrSessionNotFound       = errors.New("session not found")
	ErrSessionExpired        = errors.New("session expired")
	ErrFailedToCreateSession = errors.New("failed to create session")
	ErrFailedToGetSession    = errors.New("failed to get session")
	ErrFailedToUpdateSession = errors.New("failed to update session")
	ErrFailedToDeleteSession = errors.New("failed to delete session")

	// Ошибки валидации
	ErrEmptyLogin    = errors.New("login cannot be empty")
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrEmptyEmail    = errors.New("email cannot be empty")
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrWeakPassword  = errors.New("password is too weak")

	// Ошибки хэширования
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrPasswordVerification = errors.New("password verification failed")
)
