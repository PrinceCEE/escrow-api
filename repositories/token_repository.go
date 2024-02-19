package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	Create(ctx context.Context, t *models.Token) error
	Update(ctx context.Context, t *models.Token) error
	GetById(ctx context.Context, id string) (*models.Token, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) (time.Time, error)
}

type tokenRepository struct {
	DB *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *tokenRepository {
	return &tokenRepository{DB: db}
}

func (repo *tokenRepository) Create(ctx context.Context, t *models.Token) error {
	return nil
}

func (repo *tokenRepository) Update(ctx context.Context, t *models.Token) error {
	return nil
}

func (repo *tokenRepository) GetById(ctx context.Context, id string) (*models.Token, error) {
	return nil, nil
}

func (repo *tokenRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *tokenRepository) SoftDelete(ctx context.Context, id string) (time.Time, error) {
	return time.Now(), nil
}
