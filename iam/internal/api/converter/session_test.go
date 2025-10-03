package converter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
)

func TestToProtoSession_Success(t *testing.T) {
	// Arrange
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	session := &model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}

	// Act
	result := ToProtoSession(session)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, sessionUUID.String(), result.Uuid)
	assert.NotNil(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
	assert.NotNil(t, result.ExpiresAt)
	assert.Equal(t, now.Unix(), result.CreatedAt.AsTime().Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.AsTime().Unix())
	assert.Equal(t, expiresAt.Unix(), result.ExpiresAt.AsTime().Unix())
}

func TestToProtoSession_Nil(t *testing.T) {
	// Act
	result := ToProtoSession(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoSession_ZeroTime(t *testing.T) {
	// Arrange
	session := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  uuid.New(),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		ExpiresAt: time.Time{},
	}

	// Act
	result := ToProtoSession(session)

	// Assert
	assert.NotNil(t, result)
	assert.NotNil(t, result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
	assert.NotNil(t, result.ExpiresAt)
}

func TestToProtoSession_DifferentTimeValues(t *testing.T) {
	tests := []struct {
		name      string
		createdAt time.Time
		updatedAt time.Time
		expiresAt time.Time
	}{
		{
			name:      "одинаковые времена",
			createdAt: time.Now(),
			updatedAt: time.Now(),
			expiresAt: time.Now().Add(1 * time.Hour),
		},
		{
			name:      "разные времена",
			createdAt: time.Now().Add(-2 * time.Hour),
			updatedAt: time.Now().Add(-1 * time.Hour),
			expiresAt: time.Now().Add(24 * time.Hour),
		},
		{
			name:      "сессия с длительным сроком",
			createdAt: time.Now(),
			updatedAt: time.Now(),
			expiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 дней
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			session := &model.Session{
				UUID:      uuid.New(),
				UserUUID:  uuid.New(),
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
				ExpiresAt: tt.expiresAt,
			}

			// Act
			result := ToProtoSession(session)

			// Assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.createdAt.Unix(), result.CreatedAt.AsTime().Unix())
			assert.Equal(t, tt.updatedAt.Unix(), result.UpdatedAt.AsTime().Unix())
			assert.Equal(t, tt.expiresAt.Unix(), result.ExpiresAt.AsTime().Unix())
		})
	}
}

func TestToProtoSession_ExpiredSession(t *testing.T) {
	// Arrange - истекшая сессия
	session := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  uuid.New(),
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour), // истекла час назад
	}

	// Act
	result := ToProtoSession(session)

	// Assert
	assert.NotNil(t, result)
	assert.True(t, result.ExpiresAt.AsTime().Before(time.Now()))
}

func TestToProtoSession_FutureSession(t *testing.T) {
	// Arrange - сессия истекает в будущем
	session := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}

	// Act
	result := ToProtoSession(session)

	// Assert
	assert.NotNil(t, result)
	assert.True(t, result.ExpiresAt.AsTime().After(time.Now()))
}

func TestToProtoSession_NilUUID(t *testing.T) {
	// Arrange
	session := &model.Session{
		UUID:      uuid.Nil,
		UserUUID:  uuid.Nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Act
	result := ToProtoSession(session)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, uuid.Nil.String(), result.Uuid)
}
