package converter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

func TestToProtoUser_Success(t *testing.T) {
	// Arrange
	userUUID := uuid.New()
	now := time.Now()

	user := &model.User{
		UUID:  userUUID,
		Login: "testuser",
		Email: "test@example.com",
		NotificationMethods: []model.NotificationMethod{
			{
				ProviderName: "telegram",
				Target:       "@testuser",
			},
			{
				ProviderName: "email",
				Target:       "test@example.com",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	result := ToProtoUser(user)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, userUUID.String(), result.Uuid)
	assert.NotNil(t, result.Info)
	assert.Equal(t, "testuser", result.Info.Login)
	assert.Equal(t, "test@example.com", result.Info.Email)
	assert.Len(t, result.Info.NotificationMethods, 2)
	assert.Equal(t, now.Unix(), result.CreatedAt.AsTime().Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.AsTime().Unix())
}

func TestToProtoUser_Nil(t *testing.T) {
	// Act
	result := ToProtoUser(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoUser_WithoutNotificationMethods(t *testing.T) {
	// Arrange
	user := &model.User{
		UUID:                uuid.New(),
		Login:               "testuser",
		Email:               "test@example.com",
		NotificationMethods: []model.NotificationMethod{},
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Act
	result := ToProtoUser(user)

	// Assert
	assert.NotNil(t, result)
	assert.NotNil(t, result.Info)
	assert.Empty(t, result.Info.NotificationMethods)
}

func TestToProtoUserInfo_Success(t *testing.T) {
	// Arrange
	user := &model.User{
		Login: "testuser",
		Email: "test@example.com",
		NotificationMethods: []model.NotificationMethod{
			{
				ProviderName: "telegram",
				Target:       "123456",
			},
		},
	}

	// Act
	result := ToProtoUserInfo(user)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Login)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Len(t, result.NotificationMethods, 1)
}

func TestToProtoUserInfo_Nil(t *testing.T) {
	// Act
	result := ToProtoUserInfo(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoNotificationMethods_Success(t *testing.T) {
	// Arrange
	methods := []model.NotificationMethod{
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
			Target:       "device-token",
		},
	}

	// Act
	result := ToProtoNotificationMethods(methods)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "telegram", result[0].ProviderName)
	assert.Equal(t, "@user", result[0].Target)
	assert.Equal(t, "email", result[1].ProviderName)
	assert.Equal(t, "user@test.com", result[1].Target)
	assert.Equal(t, "push", result[2].ProviderName)
	assert.Equal(t, "device-token", result[2].Target)
}

func TestToProtoNotificationMethods_Nil(t *testing.T) {
	// Act
	result := ToProtoNotificationMethods(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoNotificationMethods_Empty(t *testing.T) {
	// Arrange
	methods := []model.NotificationMethod{}

	// Act
	result := ToProtoNotificationMethods(methods)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestToProtoNotificationMethod_Success(t *testing.T) {
	// Arrange
	method := &model.NotificationMethod{
		ProviderName: "telegram",
		Target:       "@testuser",
	}

	// Act
	result := ToProtoNotificationMethod(method)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "telegram", result.ProviderName)
	assert.Equal(t, "@testuser", result.Target)
}

func TestToProtoNotificationMethod_Nil(t *testing.T) {
	// Act
	result := ToProtoNotificationMethod(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToModelNotificationMethods_Success(t *testing.T) {
	// Arrange
	protoMethods := []*iamV1.NotificationMethod{
		{
			ProviderName: "telegram",
			Target:       "@user",
		},
		{
			ProviderName: "email",
			Target:       "user@test.com",
		},
	}

	// Act
	result := ToModelNotificationMethods(protoMethods)

	// Assert
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "telegram", result[0].ProviderName)
	assert.Equal(t, "@user", result[0].Target)
	assert.Equal(t, "email", result[1].ProviderName)
	assert.Equal(t, "user@test.com", result[1].Target)
}

func TestToModelNotificationMethods_Nil(t *testing.T) {
	// Act
	result := ToModelNotificationMethods(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToModelNotificationMethod_Success(t *testing.T) {
	// Arrange
	protoMethod := &iamV1.NotificationMethod{
		ProviderName: "telegram",
		Target:       "@testuser",
	}

	// Act
	result := ToModelNotificationMethod(protoMethod)

	// Assert
	assert.Equal(t, "telegram", result.ProviderName)
	assert.Equal(t, "@testuser", result.Target)
}

func TestToModelUserRegistrationInfo_Success(t *testing.T) {
	// Arrange
	req := &iamV1.RegisterRequest{
		Info: &iamV1.UserRegistrationInfo{
			UserInfo: &iamV1.UserInfo{
				Login: "newuser",
				Email: "newuser@example.com",
				NotificationMethods: []*iamV1.NotificationMethod{
					{
						ProviderName: "telegram",
						Target:       "@newuser",
					},
				},
			},
			Password: "SecurePassword123!",
		},
	}

	// Act
	result := ToModelUserRegistrationInfo(req)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "newuser", result.Login)
	assert.Equal(t, "newuser@example.com", result.Email)
	assert.Equal(t, "SecurePassword123!", result.Password)
	assert.Len(t, result.NotificationMethods, 1)
	assert.Equal(t, "telegram", result.NotificationMethods[0].ProviderName)
	assert.Equal(t, "@newuser", result.NotificationMethods[0].Target)
}

func TestToModelUserRegistrationInfo_Nil(t *testing.T) {
	// Act
	result := ToModelUserRegistrationInfo(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToModelUserRegistrationInfo_NilInfo(t *testing.T) {
	// Arrange
	req := &iamV1.RegisterRequest{
		Info: nil,
	}

	// Act
	result := ToModelUserRegistrationInfo(req)

	// Assert
	assert.Nil(t, result)
}

func TestToModelUserRegistrationInfo_NilUserInfo(t *testing.T) {
	// Arrange
	req := &iamV1.RegisterRequest{
		Info: &iamV1.UserRegistrationInfo{
			UserInfo: nil,
			Password: "password",
		},
	}

	// Act
	result := ToModelUserRegistrationInfo(req)

	// Assert
	assert.Nil(t, result)
}

func TestToModelUserRegistrationInfo_WithoutNotificationMethods(t *testing.T) {
	// Arrange
	req := &iamV1.RegisterRequest{
		Info: &iamV1.UserRegistrationInfo{
			UserInfo: &iamV1.UserInfo{
				Login:               "testuser",
				Email:               "test@example.com",
				NotificationMethods: []*iamV1.NotificationMethod{},
			},
			Password: "Pass123!",
		},
	}

	// Act
	result := ToModelUserRegistrationInfo(req)

	// Assert
	assert.NotNil(t, result)
	assert.Empty(t, result.NotificationMethods)
}

func TestConverters_RoundTrip_NotificationMethods(t *testing.T) {
	// Arrange - исходные методы уведомлений
	originalMethods := []model.NotificationMethod{
		{
			ProviderName: "telegram",
			Target:       "@testuser",
		},
		{
			ProviderName: "email",
			Target:       "test@example.com",
		},
	}

	// Act - конвертируем туда и обратно
	protoMethods := ToProtoNotificationMethods(originalMethods)
	convertedMethods := ToModelNotificationMethods(protoMethods)

	// Assert
	assert.Equal(t, originalMethods, convertedMethods)
}
