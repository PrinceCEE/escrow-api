package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	Create(a *models.Auth, tx pgx.Tx) error
	Update(a *models.Auth, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Auth, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type authRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewAuthRepository(db *pgxpool.Pool, timeout time.Duration) *authRepository {
	return &authRepository{DB: db, Timeout: timeout}
}

func (repo *authRepository) Create(a *models.Auth, tx pgx.Tx) error {
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

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(a.ID, a.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(a.ID, a.Version)
}

func (repo *authRepository) Update(a *models.Auth, tx pgx.Tx) error {
	a.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(a, "auths")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, query).Scan(a.Version)
	}

	return repo.DB.QueryRow(ctx, query).Scan(a.Version)
}

func (repo *authRepository) GetById(id string, tx pgx.Tx) (*models.Auth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	a := new(models.Auth)
	query := `
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
		WHERE id = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		a.ID,
		a.UserID,
		a.Password,
		a.PasswordHistory,
		a.CreatedAt,
		a.UpdatedAt,
		a.DeletedAt,
		a.Version,
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (repo *authRepository) Delete(id string, tx pgx.Tx) (err error) {
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

func (repo *authRepository) SoftDelete(id string, tx pgx.Tx) error {
	a, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	a.UpdatedAt = now
	a.DeletedAt = now

	return repo.Update(a, tx)
}
