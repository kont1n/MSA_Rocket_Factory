package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
)

var _ def.AuthService = (*service)(nil)

const (
	// Длительность сессии
	sessionDuration = 24 * time.Hour
)

type service struct {
	iamRepository repository.IAMRepository
}

func NewService(iamRepository repository.IAMRepository) *service {
	return &service{
		iamRepository: iamRepository,
	}
}

func (s *service) Login(ctx context.Context, login, password string) (*model.Session, error) {
	// Валидация входных данных
	if login == "" {
		return nil, model.ErrEmptyLogin
	}
	if password == "" {
		return nil, model.ErrEmptyPassword
	}

	// Получаем пользователя по логину
	user, err := s.iamRepository.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	// Проверяем пароль
	valid, err := s.verifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, model.ErrPasswordVerification
	}
	if !valid {
		return nil, model.ErrInvalidCredentials
	}

	// Создаем новую сессию
	session := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  user.UUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	// Сохраняем сессию в БД
	session, err = s.iamRepository.CreateSession(ctx, session)
	if err != nil {
		return nil, model.ErrFailedToCreateSession
	}

	return session, nil
}

func (s *service) Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, *model.User, error) {
	// Получаем сессию
	session, err := s.iamRepository.GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		return nil, nil, err
	}

	// Проверяем, не истекла ли сессия
	if session.IsExpired() {
		return nil, nil, model.ErrSessionExpired
	}

	// Получаем пользователя
	user, err := s.iamRepository.GetUserByUUID(ctx, session.UserUUID)
	if err != nil {
		return nil, nil, err
	}

	return session, user, nil
}

func (s *service) Logout(ctx context.Context, sessionUUID uuid.UUID) error {
	return s.iamRepository.DeleteSession(ctx, sessionUUID)
}

// verifyPassword проверяет пароль против хэша
func (s *service) verifyPassword(password, encodedHash string) (bool, error) {
	// Парсим хэш
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, model.ErrPasswordVerification
	}

	var memory, time uint32
	var threads uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, model.ErrPasswordVerification
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, model.ErrPasswordVerification
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, model.ErrPasswordVerification
	}

	// Хэшируем предоставленный пароль с теми же параметрами
	keyLen := len(decodedHash)
	if keyLen < 0 || keyLen > 0xFFFFFFFF {
		return false, model.ErrPasswordVerification
	}
	passwordHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(keyLen))

	// Сравниваем хэши
	return subtle.ConstantTimeCompare(decodedHash, passwordHash) == 1, nil
}
