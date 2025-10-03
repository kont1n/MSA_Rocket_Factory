package converter

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
)

func TestToRepoSessionPostgres_Success(t *testing.T) {
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
	result := ToRepoSessionPostgres(session)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, sessionUUID, result.UUID)
	assert.Equal(t, userUUID, result.UserUUID)
	assert.Equal(t, now.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix())
	assert.Equal(t, expiresAt.Unix(), result.ExpiresAt.Unix())
}

func TestToModelSessionFromPostgres_Success(t *testing.T) {
	// Arrange
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	repoSession := &repoModel.SessionPostgres{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
	}

	// Act
	result := ToModelSessionFromPostgres(repoSession)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, sessionUUID, result.UUID)
	assert.Equal(t, userUUID, result.UserUUID)
	assert.Equal(t, now.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix())
	assert.Equal(t, expiresAt.Unix(), result.ExpiresAt.Unix())
}

func TestToRepoSessionRedis_Success(t *testing.T) {
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
	result := ToRepoSessionRedis(session)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, sessionUUID.String(), result.UUID)
	assert.Equal(t, userUUID.String(), result.UserUUID)
	assert.Equal(t, now.Unix(), result.CreatedAt)
	assert.NotNil(t, result.UpdatedAt)
	assert.Equal(t, now.Unix(), *result.UpdatedAt)
	assert.Equal(t, expiresAt.Unix(), result.ExpiresAt)
}

func TestToRepoSessionRedis_ZeroUpdatedAt(t *testing.T) {
	// Arrange - сессия с нулевым UpdatedAt
	session := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{}, // zero time
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Act
	result := ToRepoSessionRedis(session)

	// Assert
	assert.NotNil(t, result)
	assert.Nil(t, result.UpdatedAt) // должно быть nil для zero time
}

func TestToModelSessionFromRedis_Success(t *testing.T) {
	// Arrange
	sessionUUID := uuid.New()
	userUUID := uuid.New()
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)
	updatedAt := now.Unix()

	repoSession := &repoModel.SessionRedis{
		UUID:      sessionUUID.String(),
		UserUUID:  userUUID.String(),
		CreatedAt: now.Unix(),
		UpdatedAt: &updatedAt,
		ExpiresAt: expiresAt.Unix(),
	}

	// Act
	result, err := ToModelSessionFromRedis(repoSession)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionUUID, result.UUID)
	assert.Equal(t, userUUID, result.UserUUID)
	assert.Equal(t, now.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix())
	assert.Equal(t, expiresAt.Unix(), result.ExpiresAt.Unix())
}

func TestToModelSessionFromRedis_NilUpdatedAt(t *testing.T) {
	// Arrange
	now := time.Now()
	repoSession := &repoModel.SessionRedis{
		UUID:      uuid.New().String(),
		UserUUID:  uuid.New().String(),
		CreatedAt: now.Unix(),
		UpdatedAt: nil, // nil UpdatedAt
		ExpiresAt: now.Add(1 * time.Hour).Unix(),
	}

	// Act
	result, err := ToModelSessionFromRedis(repoSession)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix()) // должно быть равно CreatedAt
}

func TestToModelSessionFromRedis_InvalidSessionUUID(t *testing.T) {
	// Arrange
	repoSession := &repoModel.SessionRedis{
		UUID:      "invalid-uuid",
		UserUUID:  uuid.New().String(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: nil,
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}

	// Act
	result, err := ToModelSessionFromRedis(repoSession)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestToModelSessionFromRedis_InvalidUserUUID(t *testing.T) {
	// Arrange
	repoSession := &repoModel.SessionRedis{
		UUID:      uuid.New().String(),
		UserUUID:  "bad-user-uuid",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: nil,
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	}

	// Act
	result, err := ToModelSessionFromRedis(repoSession)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSessionConverters_RoundTrip_Postgres(t *testing.T) {
	// Arrange
	originalSession := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Act - конвертируем туда и обратно
	repoSession := ToRepoSessionPostgres(originalSession)
	convertedSession := ToModelSessionFromPostgres(repoSession)

	// Assert
	assert.Equal(t, originalSession.UUID, convertedSession.UUID)
	assert.Equal(t, originalSession.UserUUID, convertedSession.UserUUID)
	assert.Equal(t, originalSession.CreatedAt.Unix(), convertedSession.CreatedAt.Unix())
	assert.Equal(t, originalSession.UpdatedAt.Unix(), convertedSession.UpdatedAt.Unix())
	assert.Equal(t, originalSession.ExpiresAt.Unix(), convertedSession.ExpiresAt.Unix())
}

func TestSessionConverters_RoundTrip_Redis(t *testing.T) {
	// Arrange
	originalSession := &model.Session{
		UUID:      uuid.New(),
		UserUUID:  uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Act - конвертируем туда и обратно
	repoSession := ToRepoSessionRedis(originalSession)
	convertedSession, err := ToModelSessionFromRedis(repoSession)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, originalSession.UUID, convertedSession.UUID)
	assert.Equal(t, originalSession.UserUUID, convertedSession.UserUUID)
	assert.Equal(t, originalSession.CreatedAt.Unix(), convertedSession.CreatedAt.Unix())
	assert.Equal(t, originalSession.UpdatedAt.Unix(), convertedSession.UpdatedAt.Unix())
	assert.Equal(t, originalSession.ExpiresAt.Unix(), convertedSession.ExpiresAt.Unix())
}

func TestToRepoSessionRedis_DifferentTimes(t *testing.T) {
	tests := []struct {
		name      string
		createdAt time.Time
		updatedAt time.Time
		expiresAt time.Time
	}{
		{
			name:      "короткая сессия",
			createdAt: time.Now(),
			updatedAt: time.Now(),
			expiresAt: time.Now().Add(1 * time.Hour),
		},
		{
			name:      "длинная сессия",
			createdAt: time.Now(),
			updatedAt: time.Now(),
			expiresAt: time.Now().Add(30 * 24 * time.Hour),
		},
		{
			name:      "обновленная сессия",
			createdAt: time.Now().Add(-1 * time.Hour),
			updatedAt: time.Now(),
			expiresAt: time.Now().Add(24 * time.Hour),
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
			result := ToRepoSessionRedis(session)

			// Assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.createdAt.Unix(), result.CreatedAt)
			assert.Equal(t, tt.expiresAt.Unix(), result.ExpiresAt)
		})
	}
}

func TestToModelSessionFromRedis_DifferentTimes(t *testing.T) {
	tests := []struct {
		name      string
		createdAt int64
		updatedAt *int64
		expiresAt int64
	}{
		{
			name:      "с UpdatedAt",
			createdAt: time.Now().Unix(),
			updatedAt: func() *int64 { t := time.Now().Unix(); return &t }(),
			expiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
		{
			name:      "без UpdatedAt",
			createdAt: time.Now().Unix(),
			updatedAt: nil,
			expiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repoSession := &repoModel.SessionRedis{
				UUID:      uuid.New().String(),
				UserUUID:  uuid.New().String(),
				CreatedAt: tt.createdAt,
				UpdatedAt: tt.updatedAt,
				ExpiresAt: tt.expiresAt,
			}

			// Act
			result, err := ToModelSessionFromRedis(repoSession)

			// Assert
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.createdAt, result.CreatedAt.Unix())
			assert.Equal(t, tt.expiresAt, result.ExpiresAt.Unix())
			if tt.updatedAt != nil {
				assert.Equal(t, *tt.updatedAt, result.UpdatedAt.Unix())
			} else {
				assert.Equal(t, tt.createdAt, result.UpdatedAt.Unix())
			}
		})
	}
}
