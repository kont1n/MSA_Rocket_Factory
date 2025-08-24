package user

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	"github.com/kont1n/MSA_Rocket_Factory/platform/pkg/logger"
)

var _ def.UserService = (*service)(nil)

const (
	// Параметры для Argon2
	argon2Time    = 1
	argon2Memory  = 64 * 1024
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLen       = 16
)

type service struct {
	iamRepository repository.IAMRepository
}

func NewService(iamRepository repository.IAMRepository) *service {
	return &service{
		iamRepository: iamRepository,
	}
}

func (s *service) Register(ctx context.Context, registrationInfo *model.UserRegistrationInfo) (*model.User, error) {
	// Валидация входных данных с помощью встроенной валидации модели
	if err := registrationInfo.Validate(); err != nil {
		logger.Warn(ctx, "🚫 Регистрация отклонена: некорректные данные",
			zap.String("login", registrationInfo.Login),
			zap.Error(err))
		return nil, err
	}

	// Проверяем, не существует ли уже пользователь с таким логином
	existingUser, err := s.iamRepository.GetUserByLogin(ctx, registrationInfo.Login)
	if err == nil && existingUser != nil {
		logger.Warn(ctx, "🚫 Регистрация отклонена: пользователь уже существует",
			zap.String("login", registrationInfo.Login))
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
		logger.Error(ctx, "❌ Ошибка при создании пользователя в БД",
			zap.String("login", registrationInfo.Login),
			zap.Error(err))
		return nil, err
	}

	logger.Info(ctx, "✅ Пользователь успешно зарегистрирован",
		zap.String("login", user.Login),
		zap.String("user_uuid", user.UUID.String()),
		zap.String("email", user.Email))

	return user, nil
}

func (s *service) GetUser(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	return s.iamRepository.GetUserByUUID(ctx, userUUID)
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
