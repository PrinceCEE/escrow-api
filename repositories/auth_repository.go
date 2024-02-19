package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	Create(ctx context.Context, a *models.Auth) error
	Update(ctx context.Context, a *models.Auth) error
	GetById(ctx context.Context, id string) (*models.Auth, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) (time.Time, error)
}

type authRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *authRepository {
	return &authRepository{DB: db}
}

func (repo *authRepository) Create(ctx context.Context, a *models.Auth) error {
	return nil
}

func (repo *authRepository) Update(ctx context.Context, a *models.Auth) error {
	return nil
}

func (repo *authRepository) GetById(ctx context.Context, id string) (*models.Auth, error) {
	return nil, nil
}

func (repo *authRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *authRepository) SoftDelete(ctx context.Context, id string) (time.Time, error) {
	return time.Now(), nil
}
