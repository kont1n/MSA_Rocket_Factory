package converter

import (
	"time"

	"github.com/google/uuid"

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

// ToRepoSessionRedis конвертирует внутреннюю модель сессии в модель для Redis
func ToRepoSessionRedis(session *model.Session) *repoModel.SessionRedis {
	var updatedAt *int64
	if !session.UpdatedAt.IsZero() {
		ts := session.UpdatedAt.Unix()
		updatedAt = &ts
	}

	return &repoModel.SessionRedis{
		UUID:      session.UUID.String(),
		UserUUID:  session.UserUUID.String(),
		CreatedAt: session.CreatedAt.Unix(),
		UpdatedAt: updatedAt,
		ExpiresAt: session.ExpiresAt.Unix(),
	}
}

// ToModelSessionFromRedis конвертирует модель сессии из Redis во внутреннюю модель
func ToModelSessionFromRedis(repoSession *repoModel.SessionRedis) (*model.Session, error) {
	sessionUUID, err := uuid.Parse(repoSession.UUID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(repoSession.UserUUID)
	if err != nil {
		return nil, err
	}

	createdAt := time.Unix(repoSession.CreatedAt, 0)
	expiresAt := time.Unix(repoSession.ExpiresAt, 0)

	var updatedAt time.Time
	if repoSession.UpdatedAt != nil {
		updatedAt = time.Unix(*repoSession.UpdatedAt, 0)
	} else {
		updatedAt = createdAt
	}

	return &model.Session{
		UUID:      sessionUUID,
		UserUUID:  userUUID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiresAt: expiresAt,
	}, nil
}
