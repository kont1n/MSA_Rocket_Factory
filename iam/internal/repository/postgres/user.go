package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/model"
	"github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/converter"
	repoModel "github.com/kont1n/MSA_Rocket_Factory/iam/internal/repository/model"
)

func (r *repository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	repoUser := converter.ToRepoUserPostgres(user)

	// Начинаем транзакцию для создания пользователя и способов уведомления
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, model.ErrFailedToCreateUser
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			// Ошибка rollback в defer, основная ошибка важнее
			return
		}
	}()

	// Создаем пользователя
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "login", "email", "password_hash", "created_at", "updated_at").
		Values(repoUser.UUID, repoUser.Login, repoUser.Email, repoUser.PasswordHash, repoUser.CreatedAt, repoUser.UpdatedAt)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, model.ErrFailedToCreateUser
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, model.ErrFailedToCreateUser
	}

	// Создаем способы уведомления
	if len(user.NotificationMethods) > 0 {
		for _, method := range user.NotificationMethods {
			builderNotification := sq.Insert("notification_methods").
				PlaceholderFormat(sq.Dollar).
				Columns("user_uuid", "provider_name", "target").
				Values(user.UUID, method.ProviderName, method.Target)

			query, args, err := builderNotification.ToSql()
			if err != nil {
				return nil, model.ErrFailedToCreateUser
			}

			_, err = tx.Exec(ctx, query, args...)
			if err != nil {
				return nil, model.ErrFailedToCreateUser
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, model.ErrFailedToCreateUser
	}

	return user, nil
}

func (r *repository) GetUserByUUID(ctx context.Context, userUUID uuid.UUID) (*model.User, error) {
	var repoUser repoModel.UserPostgres

	builderSelect := sq.Select("user_uuid", "login", "email", "password_hash", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_uuid": userUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, model.ErrFailedToGetUser
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&repoUser.UUID,
		&repoUser.Login,
		&repoUser.Email,
		&repoUser.PasswordHash,
		&repoUser.CreatedAt,
		&repoUser.UpdatedAt,
	)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Получаем способы уведомления
	notificationMethods, err := r.getNotificationMethods(ctx, userUUID)
	if err != nil {
		return nil, model.ErrFailedToGetUser
	}

	user := converter.ToModelUserFromPostgres(&repoUser, notificationMethods)
	return user, nil
}

func (r *repository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var repoUser repoModel.UserPostgres

	builderSelect := sq.Select("user_uuid", "login", "email", "password_hash", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"login": login})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, model.ErrFailedToGetUser
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&repoUser.UUID,
		&repoUser.Login,
		&repoUser.Email,
		&repoUser.PasswordHash,
		&repoUser.CreatedAt,
		&repoUser.UpdatedAt,
	)
	if err != nil {
		return nil, model.ErrUserNotFound
	}

	// Получаем способы уведомления
	notificationMethods, err := r.getNotificationMethods(ctx, repoUser.UUID)
	if err != nil {
		return nil, model.ErrFailedToGetUser
	}

	user := converter.ToModelUserFromPostgres(&repoUser, notificationMethods)
	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	repoUser := converter.ToRepoUserPostgres(user)

	builderUpdate := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set("login", repoUser.Login).
		Set("email", repoUser.Email).
		Set("password_hash", repoUser.PasswordHash).
		Set("updated_at", repoUser.UpdatedAt).
		Where(sq.Eq{"user_uuid": repoUser.UUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, model.ErrFailedToUpdateUser
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, model.ErrFailedToUpdateUser
	}

	return user, nil
}

// getNotificationMethods получает все способы уведомления для пользователя
func (r *repository) getNotificationMethods(ctx context.Context, userUUID uuid.UUID) ([]model.NotificationMethod, error) {
	builderSelect := sq.Select("provider_name", "target").
		From("notification_methods").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_uuid": userUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var methods []model.NotificationMethod
	for rows.Next() {
		var repoMethod repoModel.NotificationMethodPostgres
		err = rows.Scan(&repoMethod.ProviderName, &repoMethod.Target)
		if err != nil {
			return nil, err
		}
		methods = append(methods, converter.ToModelNotificationMethod(&repoMethod))
	}

	return methods, nil
}
