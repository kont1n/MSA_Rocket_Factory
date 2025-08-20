package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
)

var _ def.UserService = (*service)(nil)

const (
	// Параметры для Argon2
	argon2Time    = 1
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLen       = 16

	// Минимальная длина пароля
	minPasswordLength = 8
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

type service struct {
	iamRepository repository.IAMRepository
}

func NewService(iamRepository repository.IAMRepository) *service {
	return &service{
		iamRepository: iamRepository,
	}
}

func (s *service) Register(ctx context.Context, registrationInfo *model.UserRegistrationInfo) (*model.User, error) {
	// Валидация входных данных
	if err := s.validateRegistrationInfo(registrationInfo); err != nil {
		return nil, err
	}

	// Проверяем, не существует ли уже пользователь с таким логином
	existingUser, err := s.iamRepository.GetUserByLogin(ctx, registrationInfo.Login)
	if err == nil && existingUser != nil {
		return nil, model.ErrUserAlreadyExists
	}

	// Хэшируем пароль
	passwordHash, err := s.hashPassword(registrationInfo.Password)
	if err != nil {
		return nil, err
	}

	// Создаем нового пользователя
	user := &model.User{
		UUID:                uuid.New(),
		Login:               registrationInfo.Login,
		Email:               registrationInfo.Email,
		PasswordHash:        passwordHash,
		NotificationMethods: registrationInfo.NotificationMethods,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Сохраняем пользователя в БД
	user, err = s.iamRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetUser(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	return s.iamRepository.GetUserByUUID(ctx, userUUID)
}

// validateRegistrationInfo валидирует данные для регистрации
func (s *service) validateRegistrationInfo(info *model.UserRegistrationInfo) error {
	if info.Login == "" {
		return model.ErrEmptyLogin
	}

	if info.Password == "" {
		return model.ErrEmptyPassword
	}

	if len(info.Password) < minPasswordLength {
		return model.ErrWeakPassword
	}

	if info.Email == "" {
		return model.ErrEmptyEmail
	}

	if !s.isValidEmail(info.Email) {
		return model.ErrInvalidEmail
	}

	return nil
}

// isValidEmail проверяет корректность email адреса
func (s *service) isValidEmail(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	return emailRegex.MatchString(email)
}

// hashPassword хэширует пароль с использованием Argon2
func (s *service) hashPassword(password string) (string, error) {
	// Генерируем соль
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", model.ErrFailedToHashPassword
	}

	// Хэшируем пароль
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Кодируем результат в base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Формат: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	encodedHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2Memory, argon2Time, argon2Threads, b64Salt, b64Hash)

	return encodedHash, nil
}
