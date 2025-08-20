package converter

import (
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
)

// ToRepoSessionPostgres конвертирует внутреннюю модель сессии в модель для PostgreSQL
func ToRepoSessionPostgres(session *model.Session) *repoModel.SessionPostgres {
	return &repoModel.SessionPostgres{
		UUID:      session.UUID,
		UserUUID:  session.UserUUID,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
		ExpiresAt: session.ExpiresAt,
	}
}

// ToModelSessionFromPostgres конвертирует модель сессии из PostgreSQL во внутреннюю модель
func ToModelSessionFromPostgres(repoSession *repoModel.SessionPostgres) *model.Session {
	return &model.Session{
		UUID:      repoSession.UUID,
		UserUUID:  repoSession.UserUUID,
		CreatedAt: repoSession.CreatedAt,
		UpdatedAt: repoSession.UpdatedAt,
		ExpiresAt: repoSession.ExpiresAt,
	}
}
