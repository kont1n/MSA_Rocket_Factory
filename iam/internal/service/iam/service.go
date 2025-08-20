package iam

import (
	"context"

	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository"
	def "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service"
	authService "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service/auth"
	userService "github.com/kont1n/MSA_Rocket_Factory/iam/internal/service/user"
)

var _ def.IAMService = (*service)(nil)

type service struct {
	authService def.AuthService
	userService def.UserService
}

func NewService(iamRepository repository.IAMRepository) *service {
	return &service{
		authService: authService.NewService(iamRepository),
		userService: userService.NewService(iamRepository),
	}
}

// AuthService методы
func (s *service) Login(ctx context.Context, login, password string) (*model.Session, error) {
	return s.authService.Login(ctx, login, password)
}

func (s *service) Whoami(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, *model.User, error) {
	return s.authService.Whoami(ctx, sessionUUID)
}

func (s *service) Logout(ctx context.Context, sessionUUID uuid.UUID) error {
	return s.authService.Logout(ctx, sessionUUID)
}

// UserService методы
func (s *service) Register(ctx context.Context, registrationInfo *model.UserRegistrationInfo) (*model.User, error) {
	return s.userService.Register(ctx, registrationInfo)
}

func (s *service) GetUser(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	return s.userService.GetUser(ctx, userUUID)
}
