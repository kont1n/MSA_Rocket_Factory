package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	jwtV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/jwt/v1"
)

// ToProtoLoginResponse конвертирует TokenPair в protobuf LoginResponse
func ToProtoLoginResponse(tokenPair *model.TokenPair) *jwtV1.LoginResponse {
	if tokenPair == nil {
		return nil
	}

	response := &jwtV1.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	// Добавляем время истечения токенов, если они заданы
	if !tokenPair.AccessTokenExpiresAt.IsZero() {
		response.AccessTokenExpiresAt = timestamppb.New(tokenPair.AccessTokenExpiresAt)
	}
	if !tokenPair.RefreshTokenExpiresAt.IsZero() {
		response.RefreshTokenExpiresAt = timestamppb.New(tokenPair.RefreshTokenExpiresAt)
	}

	return response
}

// ToProtoGetAccessTokenResponse конвертирует TokenPair в protobuf GetAccessTokenResponse
func ToProtoGetAccessTokenResponse(tokenPair *model.TokenPair) *jwtV1.GetAccessTokenResponse {
	if tokenPair == nil {
		return nil
	}

	response := &jwtV1.GetAccessTokenResponse{
		AccessToken: tokenPair.AccessToken,
	}

	// Добавляем время истечения токена, если оно задано
	if !tokenPair.AccessTokenExpiresAt.IsZero() {
		response.AccessTokenExpiresAt = timestamppb.New(tokenPair.AccessTokenExpiresAt)
	}

	return response
}

// ToProtoGetRefreshTokenResponse конвертирует TokenPair в protobuf GetRefreshTokenResponse
func ToProtoGetRefreshTokenResponse(tokenPair *model.TokenPair) *jwtV1.GetRefreshTokenResponse {
	if tokenPair == nil {
		return nil
	}

	response := &jwtV1.GetRefreshTokenResponse{
		RefreshToken: tokenPair.RefreshToken,
	}

	// Добавляем время истечения токена, если оно задано
	if !tokenPair.RefreshTokenExpiresAt.IsZero() {
		response.RefreshTokenExpiresAt = timestamppb.New(tokenPair.RefreshTokenExpiresAt)
	}

	return response
}
