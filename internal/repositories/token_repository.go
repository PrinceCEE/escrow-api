package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/pkg/utils"
)

type ITokenRepository interface {
	Create(t *models.Token, tx pgx.Tx) error
	Update(t *models.Token, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Token, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type TokenRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewTokenRepository(db *pgxpool.Pool, timeout time.Duration) *TokenRepository {
	return &TokenRepository{DB: db, Timeout: timeout}
}

func (repo *TokenRepository) Create(t *models.Token, tx pgx.Tx) error {
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

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &t.Version)
		if err != nil {
			return err
		}

		t.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &t.Version)
	if err != nil {
		return err
	}

	t.ID = id.String()
	return nil
}

func (repo *TokenRepository) Update(t *models.Token, tx pgx.Tx) error {
	t.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(t, "tokens")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&t.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&t.Version)
}

func (repo *TokenRepository) getByKey(key string, value any, tx pgx.Tx) (*models.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var id, userId uuid.UUID
	t := new(models.Token)
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

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&id,
		&t.Hash,
		&userId,
		&t.TokenType,
		&t.InUse,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
		&t.Version,
	)
	if err != nil {
		return nil, err
	}

	t.ID = id.String()
	t.UserID = userId.String()
	return t, nil
}

func (repo *TokenRepository) GetById(id string, tx pgx.Tx) (*models.Token, error) {
	return repo.getByKey("id", id, tx)
}

func (repo *TokenRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM tokens WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *TokenRepository) SoftDelete(id string, tx pgx.Tx) error {
	t, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	t.UpdatedAt = now
	t.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(t, tx)
}
