package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
)

func TestToProtoLoginResponse_Success(t *testing.T) {
	// Arrange
	now := time.Now()
	tokenPair := &model.TokenPair{
		AccessToken:           "access-token-value",
		RefreshToken:          "refresh-token-value",
		AccessTokenExpiresAt:  now.Add(15 * time.Minute),
		RefreshTokenExpiresAt: now.Add(7 * 24 * time.Hour),
	}

	// Act
	result := ToProtoLoginResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "access-token-value", result.AccessToken)
	assert.Equal(t, "refresh-token-value", result.RefreshToken)
	assert.NotNil(t, result.AccessTokenExpiresAt)
	assert.NotNil(t, result.RefreshTokenExpiresAt)
	assert.Equal(t, now.Add(15*time.Minute).Unix(), result.AccessTokenExpiresAt.AsTime().Unix())
	assert.Equal(t, now.Add(7*24*time.Hour).Unix(), result.RefreshTokenExpiresAt.AsTime().Unix())
}

func TestToProtoLoginResponse_Nil(t *testing.T) {
	// Act
	result := ToProtoLoginResponse(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoLoginResponse_WithoutExpirationTimes(t *testing.T) {
	// Arrange - токены без времени истечения
	tokenPair := &model.TokenPair{
		AccessToken:           "access-token",
		RefreshToken:          "refresh-token",
		AccessTokenExpiresAt:  time.Time{}, // zero time
		RefreshTokenExpiresAt: time.Time{}, // zero time
	}

	// Act
	result := ToProtoLoginResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	// Время истечения не должно быть установлено
	assert.Nil(t, result.AccessTokenExpiresAt)
	assert.Nil(t, result.RefreshTokenExpiresAt)
}

func TestToProtoLoginResponse_EmptyTokens(t *testing.T) {
	// Arrange
	tokenPair := &model.TokenPair{
		AccessToken:           "",
		RefreshToken:          "",
		AccessTokenExpiresAt:  time.Now(),
		RefreshTokenExpiresAt: time.Now(),
	}

	// Act
	result := ToProtoLoginResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "", result.AccessToken)
	assert.Equal(t, "", result.RefreshToken)
}

func TestToProtoGetAccessTokenResponse_Success(t *testing.T) {
	// Arrange
	now := time.Now()
	tokenPair := &model.TokenPair{
		AccessToken:          "new-access-token",
		AccessTokenExpiresAt: now.Add(15 * time.Minute),
	}

	// Act
	result := ToProtoGetAccessTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.NotNil(t, result.AccessTokenExpiresAt)
	assert.Equal(t, now.Add(15*time.Minute).Unix(), result.AccessTokenExpiresAt.AsTime().Unix())
}

func TestToProtoGetAccessTokenResponse_Nil(t *testing.T) {
	// Act
	result := ToProtoGetAccessTokenResponse(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoGetAccessTokenResponse_WithoutExpirationTime(t *testing.T) {
	// Arrange
	tokenPair := &model.TokenPair{
		AccessToken:          "access-token",
		AccessTokenExpiresAt: time.Time{}, // zero time
	}

	// Act
	result := ToProtoGetAccessTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Nil(t, result.AccessTokenExpiresAt)
}

func TestToProtoGetAccessTokenResponse_EmptyToken(t *testing.T) {
	// Arrange
	tokenPair := &model.TokenPair{
		AccessToken:          "",
		AccessTokenExpiresAt: time.Now(),
	}

	// Act
	result := ToProtoGetAccessTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "", result.AccessToken)
}

func TestToProtoGetRefreshTokenResponse_Success(t *testing.T) {
	// Arrange
	now := time.Now()
	tokenPair := &model.TokenPair{
		RefreshToken:          "new-refresh-token",
		RefreshTokenExpiresAt: now.Add(7 * 24 * time.Hour),
	}

	// Act
	result := ToProtoGetRefreshTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "new-refresh-token", result.RefreshToken)
	assert.NotNil(t, result.RefreshTokenExpiresAt)
	assert.Equal(t, now.Add(7*24*time.Hour).Unix(), result.RefreshTokenExpiresAt.AsTime().Unix())
}

func TestToProtoGetRefreshTokenResponse_Nil(t *testing.T) {
	// Act
	result := ToProtoGetRefreshTokenResponse(nil)

	// Assert
	assert.Nil(t, result)
}

func TestToProtoGetRefreshTokenResponse_WithoutExpirationTime(t *testing.T) {
	// Arrange
	tokenPair := &model.TokenPair{
		RefreshToken:          "refresh-token",
		RefreshTokenExpiresAt: time.Time{}, // zero time
	}

	// Act
	result := ToProtoGetRefreshTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	assert.Nil(t, result.RefreshTokenExpiresAt)
}

func TestToProtoGetRefreshTokenResponse_EmptyToken(t *testing.T) {
	// Arrange
	tokenPair := &model.TokenPair{
		RefreshToken:          "",
		RefreshTokenExpiresAt: time.Now(),
	}

	// Act
	result := ToProtoGetRefreshTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "", result.RefreshToken)
}

func TestToProtoLoginResponse_DifferentExpirationTimes(t *testing.T) {
	tests := []struct {
		name                 string
		accessTokenDuration  time.Duration
		refreshTokenDuration time.Duration
	}{
		{
			name:                 "короткий access, длинный refresh",
			accessTokenDuration:  5 * time.Minute,
			refreshTokenDuration: 30 * 24 * time.Hour,
		},
		{
			name:                 "стандартные сроки",
			accessTokenDuration:  15 * time.Minute,
			refreshTokenDuration: 7 * 24 * time.Hour,
		},
		{
			name:                 "длинные сроки",
			accessTokenDuration:  1 * time.Hour,
			refreshTokenDuration: 90 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			now := time.Now()
			tokenPair := &model.TokenPair{
				AccessToken:           "access",
				RefreshToken:          "refresh",
				AccessTokenExpiresAt:  now.Add(tt.accessTokenDuration),
				RefreshTokenExpiresAt: now.Add(tt.refreshTokenDuration),
			}

			// Act
			result := ToProtoLoginResponse(tokenPair)

			// Assert
			assert.NotNil(t, result)
			assert.True(t, result.RefreshTokenExpiresAt.AsTime().After(result.AccessTokenExpiresAt.AsTime()))
		})
	}
}

func TestToProtoGetAccessTokenResponse_LongToken(t *testing.T) {
	// Arrange - очень длинный токен
	longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	tokenPair := &model.TokenPair{
		AccessToken:          longToken,
		AccessTokenExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// Act
	result := ToProtoGetAccessTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, longToken, result.AccessToken)
}

func TestToProtoGetRefreshTokenResponse_LongToken(t *testing.T) {
	// Arrange - очень длинный токен
	longToken := "very.long.refresh.token.with.many.segments.and.data.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.multiple.segments"
	tokenPair := &model.TokenPair{
		RefreshToken:          longToken,
		RefreshTokenExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	// Act
	result := ToProtoGetRefreshTokenResponse(tokenPair)

	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, longToken, result.RefreshToken)
}
