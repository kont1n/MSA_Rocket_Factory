package converter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
)

func TestToRepoUserPostgres_Success(t *testing.T) {
	// Arrange
	userUUID := uuid.New()
	now := time.Now()

	user := &model.User{
		UUID:         userUUID,
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password_value",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Act
	result := ToRepoUserPostgres(user)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, userUUID, result.UUID)
	assert.Equal(t, "testuser", result.Login)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "hashed_password_value", result.PasswordHash)
	assert.Equal(t, now.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix())
}

func TestToRepoUserPostgres_EmptyFields(t *testing.T) {
	// Arrange
	user := &model.User{
		UUID:         uuid.New(),
		Login:        "",
		Email:        "",
		PasswordHash: "",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}

	// Act
	result := ToRepoUserPostgres(user)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "", result.Login)
	assert.Equal(t, "", result.Email)
	assert.Equal(t, "", result.PasswordHash)
}

func TestToModelUserFromPostgres_Success(t *testing.T) {
	// Arrange
	userUUID := uuid.New()
	now := time.Now()

	repoUser := &repoModel.UserPostgres{
		UUID:         userUUID,
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	notificationMethods := []model.NotificationMethod{
		{
			ProviderName: "telegram",
			Target:       "@testuser",
		},
		{
			ProviderName: "email",
			Target:       "test@example.com",
		},
	}

	// Act
	result := ToModelUserFromPostgres(repoUser, notificationMethods)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, userUUID, result.UUID)
	assert.Equal(t, "testuser", result.Login)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "hashed_password", result.PasswordHash)
	assert.Len(t, result.NotificationMethods, 2)
	assert.Equal(t, "telegram", result.NotificationMethods[0].ProviderName)
	assert.Equal(t, "@testuser", result.NotificationMethods[0].Target)
	assert.Equal(t, now.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix())
}

func TestToModelUserFromPostgres_WithoutNotificationMethods(t *testing.T) {
	// Arrange
	repoUser := &repoModel.UserPostgres{
		UUID:         uuid.New(),
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Act
	result := ToModelUserFromPostgres(repoUser, []model.NotificationMethod{})

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result.NotificationMethods)
}

func TestToModelUserFromPostgres_NilNotificationMethods(t *testing.T) {
	// Arrange
	repoUser := &repoModel.UserPostgres{
		UUID:         uuid.New(),
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Act
	result := ToModelUserFromPostgres(repoUser, nil)

	// Assert
	assert.NotNil(t, result)
	assert.Nil(t, result.NotificationMethods)
}

func TestToRepoNotificationMethodPostgres_Success(t *testing.T) {
	// Arrange
	userUUID := uuid.New().String()
	method := model.NotificationMethod{
		ProviderName: "telegram",
		Target:       "@testuser",
	}

	// Act
	result := ToRepoNotificationMethodPostgres(userUUID, method)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "telegram", result.ProviderName)
	assert.Equal(t, "@testuser", result.Target)
}

func TestToRepoNotificationMethodPostgres_DifferentProviders(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		target       string
	}{
		{
			name:         "telegram",
			providerName: "telegram",
			target:       "@user123",
		},
		{
			name:         "email",
			providerName: "email",
			target:       "user@example.com",
		},
		{
			name:         "push",
			providerName: "push",
			target:       "device-token-12345",
		},
		{
			name:         "sms",
			providerName: "sms",
			target:       "+79001234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			userUUID := uuid.New().String()
			method := model.NotificationMethod{
				ProviderName: tt.providerName,
				Target:       tt.target,
			}

			// Act
			result := ToRepoNotificationMethodPostgres(userUUID, method)

			// Assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.providerName, result.ProviderName)
			assert.Equal(t, tt.target, result.Target)
		})
	}
}

func TestToModelNotificationMethod_Success(t *testing.T) {
	// Arrange
	repoMethod := &repoModel.NotificationMethodPostgres{
		ProviderName: "telegram",
		Target:       "@testuser",
	}

	// Act
	result := ToModelNotificationMethod(repoMethod)

	// Assert
	assert.Equal(t, "telegram", result.ProviderName)
	assert.Equal(t, "@testuser", result.Target)
}

func TestToModelNotificationMethod_EmptyFields(t *testing.T) {
	// Arrange
	repoMethod := &repoModel.NotificationMethodPostgres{
		ProviderName: "",
		Target:       "",
	}

	// Act
	result := ToModelNotificationMethod(repoMethod)

	// Assert
	assert.Equal(t, "", result.ProviderName)
	assert.Equal(t, "", result.Target)
}

func TestUserConverters_RoundTrip(t *testing.T) {
	// Arrange
	originalUser := &model.User{
		UUID:         uuid.New(),
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	notificationMethods := []model.NotificationMethod{
		{
			ProviderName: "telegram",
			Target:       "@testuser",
		},
	}

	// Act - конвертируем туда и обратно
	repoUser := ToRepoUserPostgres(originalUser)
	convertedUser := ToModelUserFromPostgres(repoUser, notificationMethods)

	// Assert
	assert.Equal(t, originalUser.UUID, convertedUser.UUID)
	assert.Equal(t, originalUser.Login, convertedUser.Login)
	assert.Equal(t, originalUser.Email, convertedUser.Email)
	assert.Equal(t, originalUser.PasswordHash, convertedUser.PasswordHash)
	assert.Equal(t, originalUser.CreatedAt.Unix(), convertedUser.CreatedAt.Unix())
	assert.Equal(t, originalUser.UpdatedAt.Unix(), convertedUser.UpdatedAt.Unix())
}

func TestNotificationMethodConverters_RoundTrip(t *testing.T) {
	// Arrange
	userUUID := uuid.New().String()
	originalMethod := model.NotificationMethod{
		ProviderName: "telegram",
		Target:       "@testuser",
	}

	// Act - конвертируем туда и обратно
	repoMethod := ToRepoNotificationMethodPostgres(userUUID, originalMethod)
	convertedMethod := ToModelNotificationMethod(repoMethod)

	// Assert
	assert.Equal(t, originalMethod.ProviderName, convertedMethod.ProviderName)
	assert.Equal(t, originalMethod.Target, convertedMethod.Target)
}

func TestToModelUserFromPostgres_MultipleNotificationMethods(t *testing.T) {
	// Arrange
	repoUser := &repoModel.UserPostgres{
		UUID:         uuid.New(),
		Login:        "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	notificationMethods := []model.NotificationMethod{
		{
			ProviderName: "telegram",
			Target:       "@user",
		},
		{
			ProviderName: "email",
			Target:       "user@test.com",
		},
		{
			ProviderName: "push",
			Target:       "token123",
		},
	}

	// Act
	result := ToModelUserFromPostgres(repoUser, notificationMethods)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result.NotificationMethods, 3)
	assert.Equal(t, "telegram", result.NotificationMethods[0].ProviderName)
	assert.Equal(t, "email", result.NotificationMethods[1].ProviderName)
	assert.Equal(t, "push", result.NotificationMethods[2].ProviderName)
}

func TestToRepoUserPostgres_DifferentPasswordHashes(t *testing.T) {
	tests := []struct {
		name         string
		passwordHash string
	}{
		{
			name:         "bcrypt hash",
			passwordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		},
		{
			name:         "argon2 hash",
			passwordHash: "$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG",
		},
		{
			name:         "empty hash",
			passwordHash: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			user := &model.User{
				UUID:         uuid.New(),
				Login:        "testuser",
				Email:        "test@example.com",
				PasswordHash: tt.passwordHash,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			// Act
			result := ToRepoUserPostgres(user)

			// Assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.passwordHash, result.PasswordHash)
		})
	}
}

func TestToModelUserFromPostgres_DifferentTimes(t *testing.T) {
	tests := []struct {
		name      string
		createdAt time.Time
		updatedAt time.Time
	}{
		{
			name:      "одинаковые времена",
			createdAt: time.Now(),
			updatedAt: time.Now(),
		},
		{
			name:      "разные времена",
			createdAt: time.Now().Add(-24 * time.Hour),
			updatedAt: time.Now(),
		},
		{
			name:      "старый пользователь",
			createdAt: time.Now().Add(-365 * 24 * time.Hour),
			updatedAt: time.Now().Add(-30 * 24 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repoUser := &repoModel.UserPostgres{
				UUID:         uuid.New(),
				Login:        "testuser",
				Email:        "test@example.com",
				PasswordHash: "hash",
				CreatedAt:    tt.createdAt,
				UpdatedAt:    tt.updatedAt,
			}

			// Act
			result := ToModelUserFromPostgres(repoUser, nil)

			// Assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.createdAt.Unix(), result.CreatedAt.Unix())
			assert.Equal(t, tt.updatedAt.Unix(), result.UpdatedAt.Unix())
		})
	}
}
