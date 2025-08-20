package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
)

func (r *repository) CreateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	repoSession := converter.ToRepoSessionPostgres(session)

	builderInsert := sq.Insert("sessions").
		PlaceholderFormat(sq.Dollar).
		Columns("session_uuid", "user_uuid", "created_at", "updated_at", "expires_at").
		Values(repoSession.UUID, repoSession.UserUUID, repoSession.CreatedAt, repoSession.UpdatedAt, repoSession.ExpiresAt)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, model.ErrFailedToCreateSession
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, model.ErrFailedToCreateSession
	}

	return session, nil
}

func (r *repository) GetSessionByUUID(ctx context.Context, sessionUUID uuid.UUID) (*model.Session, error) {
	var repoSession repoModel.SessionPostgres

	builderSelect := sq.Select("session_uuid", "user_uuid", "created_at", "updated_at", "expires_at").
		From("sessions").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"session_uuid": sessionUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, model.ErrFailedToGetSession
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&repoSession.UUID,
		&repoSession.UserUUID,
		&repoSession.CreatedAt,
		&repoSession.UpdatedAt,
		&repoSession.ExpiresAt,
	)
	if err != nil {
		return nil, model.ErrSessionNotFound
	}

	session := converter.ToModelSessionFromPostgres(&repoSession)
	return session, nil
}

func (r *repository) UpdateSession(ctx context.Context, session *model.Session) (*model.Session, error) {
	repoSession := converter.ToRepoSessionPostgres(session)

	builderUpdate := sq.Update("sessions").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", repoSession.UpdatedAt).
		Set("expires_at", repoSession.ExpiresAt).
		Where(sq.Eq{"session_uuid": repoSession.UUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, model.ErrFailedToUpdateSession
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, model.ErrFailedToUpdateSession
	}

	return session, nil
}

func (r *repository) DeleteSession(ctx context.Context, sessionUUID uuid.UUID) error {
	builderDelete := sq.Delete("sessions").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"session_uuid": sessionUUID})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return model.ErrFailedToDeleteSession
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return model.ErrFailedToDeleteSession
	}

	return nil
}

func (r *repository) CleanupExpiredSessions(ctx context.Context) error {
	builderDelete := sq.Delete("sessions").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Lt{"expires_at": "NOW()"})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return model.ErrFailedToDeleteSession
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return model.ErrFailedToDeleteSession
	}

	return nil
}
