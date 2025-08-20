package converter

import (
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
)

// ToRepoUserPostgres конвертирует внутреннюю модель пользователя в модель для PostgreSQL
func ToRepoUserPostgres(user *model.User) *repoModel.UserPostgres {
	return &repoModel.UserPostgres{
		UUID:         user.UUID,
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

// ToModelUserFromPostgres конвертирует модель пользователя из PostgreSQL во внутреннюю модель
func ToModelUserFromPostgres(repoUser *repoModel.UserPostgres, notificationMethods []model.NotificationMethod) *model.User {
	return &model.User{
		UUID:                repoUser.UUID,
		Login:               repoUser.Login,
		Email:               repoUser.Email,
		PasswordHash:        repoUser.PasswordHash,
		NotificationMethods: notificationMethods,
		CreatedAt:           repoUser.CreatedAt,
		UpdatedAt:           repoUser.UpdatedAt,
	}
}

// ToRepoNotificationMethodPostgres конвертирует способ уведомления для PostgreSQL
func ToRepoNotificationMethodPostgres(userUUID string, method model.NotificationMethod) *repoModel.NotificationMethodPostgres {
	return &repoModel.NotificationMethodPostgres{
		ProviderName: method.ProviderName,
		Target:       method.Target,
	}
}

// ToModelNotificationMethod конвертирует способ уведомления из PostgreSQL во внутреннюю модель
func ToModelNotificationMethod(repoMethod *repoModel.NotificationMethodPostgres) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: repoMethod.ProviderName,
		Target:       repoMethod.Target,
	}
}
