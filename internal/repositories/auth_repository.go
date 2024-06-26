package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/pkg/utils"
)

type IAuthRepository interface {
	Create(a *models.Auth, tx pgx.Tx) error
	Update(a *models.Auth, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Auth, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
	GetByUserId(id string, tx pgx.Tx) (*models.Auth, error)
}

type AuthRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewAuthRepository(db *pgxpool.Pool, timeout time.Duration) *AuthRepository {
	return &AuthRepository{DB: db, Timeout: timeout}
}

func (repo *AuthRepository) Create(a *models.Auth, tx pgx.Tx) error {
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO auths (user_id, password, password_history, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version
	`

	args := []any{a.UserID, a.Password, a.PasswordHistory, a.CreatedAt, a.UpdatedAt}

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &a.Version)
		if err != nil {
			return err
		}

		a.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &a.Version)
	if err != nil {
		return err
	}

	a.ID = id.String()
	return nil
}

func (repo *AuthRepository) Update(a *models.Auth, tx pgx.Tx) error {
	a.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(a, "auths")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&a.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&a.Version)
}

func (repo *AuthRepository) getByKey(key string, value any, tx pgx.Tx) (*models.Auth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var id, userId uuid.UUID

	a := new(models.Auth)
	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			password,
			password_history,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			auths
		WHERE %s = $1
	`, key)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&id,
		&userId,
		&a.Password,
		&a.PasswordHistory,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.Version,
	)
	if err != nil {
		return nil, err
	}

	a.ID = id.String()
	*a.UserID = userId.String()

	return a, nil
}

func (repo *AuthRepository) GetById(id string, tx pgx.Tx) (*models.Auth, error) {
	return repo.getByKey("id", id, tx)
}

func (repo *AuthRepository) GetByUserId(id string, tx pgx.Tx) (*models.Auth, error) {
	return repo.getByKey("user_id", id, tx)
}

func (repo *AuthRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM auths WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *AuthRepository) SoftDelete(id string, tx pgx.Tx) error {
	a, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	a.UpdatedAt = now
	a.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(a, tx)
}
