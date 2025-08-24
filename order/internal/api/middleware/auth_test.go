package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	httpAuth "github.com/kont1n/MSA_Rocket_Factory/platform/pkg/middleware/http"
	iamV1 "github.com/kont1n/MSA_Rocket_Factory/shared/pkg/proto/iam/v1"
)

// MockIAMClient мок для IAM клиента
type MockIAMClient struct {
	mock.Mock
}

func (m *MockIAMClient) Login(ctx context.Context, req *iamV1.LoginRequest, opts ...grpc.CallOption) (*iamV1.LoginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*iamV1.LoginResponse), args.Error(1)
}

func (m *MockIAMClient) Whoami(ctx context.Context, req *iamV1.WhoamiRequest, opts ...grpc.CallOption) (*iamV1.WhoamiResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*iamV1.WhoamiResponse), args.Error(1)
}

func TestAuthMiddleware_Handle_MissingSession(t *testing.T) {
	// Arrange
	mockIAM := &MockIAMClient{}
	middleware := NewAuthMiddleware(mockIAM)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	// Act
	middleware.Handle(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.False(t, handlerCalled, "Следующий handler не должен быть вызван")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "MISSING_SESSION")
}

func TestAuthMiddleware_Handle_ValidSession(t *testing.T) {
	// Arrange
	mockIAM := &MockIAMClient{}
	middleware := NewAuthMiddleware(mockIAM)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(httpAuth.SessionUUIDHeader, "test-session-uuid")
	w := httptest.NewRecorder()

	// Мокаем успешный ответ от IAM
	user := &iamV1.User{
		Uuid: "test-user-uuid",
		Info: &iamV1.UserInfo{
			Email: "test@example.com",
		},
	}

	mockIAM.On("Whoami", mock.Anything, &iamV1.WhoamiRequest{
		SessionUuid: "test-session-uuid",
	}).Return(&iamV1.WhoamiResponse{User: user}, nil)

	handlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		// Проверяем, что пользователь добавлен в контекст
		ctx := r.Context()
		userFromCtx, ok := GetUserFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, user.Uuid, userFromCtx.Uuid)

		sessionUUID, ok := GetSessionUUIDFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "test-session-uuid", sessionUUID)
	})

	// Act
	middleware.Handle(nextHandler).ServeHTTP(w, req)

	// Assert
	assert.True(t, handlerCalled, "Следующий handler должен быть вызван")
	mockIAM.AssertExpectations(t)
}
