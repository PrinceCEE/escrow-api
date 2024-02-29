package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	Create(t *models.Token) error
	Update(t *models.Token) error
	GetById(id string) (*models.Token, error)
	Delete(id string) error
	SoftDelete(id string) error
}

type tokenRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewTokenRepository(db *pgxpool.Pool, timeout time.Duration) *tokenRepository {
	return &tokenRepository{DB: db, Timeout: timeout}
}

func (repo *tokenRepository) Create(t *models.Token) error {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO tokens (hash, user_id, token_type, in_use, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version
	`

	args := []any{
		t.Hash,
		t.UserID,
		t.TokenType,
		t.InUse,
		t.CreatedAt,
		t.UpdatedAt,
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(t.ID, t.Version)
}

func (repo *tokenRepository) Update(t *models.Token) error {
	t.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(t, "tokens")
	if err != nil {
		return err
	}

	return repo.DB.QueryRow(ctx, query).Scan(t.Version)
}

func (repo *tokenRepository) GetById(id string) (*models.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var t *models.Token
	query := `
		SELECT
			id,
			hash,
			user_id,
			token_type,
			in_use,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			tokens
		WHERE id = $1
	`

	err := repo.DB.QueryRow(ctx, query, id).Scan(
		t.ID,
		t.Hash,
		t.UserID,
		t.TokenType,
		t.InUse,
		t.CreatedAt,
		t.UpdatedAt,
		t.DeletedAt,
		t.Version,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (repo *tokenRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM tokens WHERE id = $1`

	_, err := repo.DB.Exec(ctx, query)
	return err
}

func (repo *tokenRepository) SoftDelete(id string) error {
	t, err := repo.GetById(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	t.UpdatedAt = now
	t.DeletedAt = now

	return repo.Update(t)
}
