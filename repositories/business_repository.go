package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BusinessRepository interface {
	Create(ctx context.Context, b *models.Business) error
	Update(ctx context.Context, b *models.Business) error
	GetById(ctx context.Context, id string) (*models.Business, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) (time.Time, error)
}

type businessRepository struct {
	DB *pgxpool.Pool
}

func NewBusinessRepository(db *pgxpool.Pool) *businessRepository {
	return &businessRepository{DB: db}
}

func (repo *businessRepository) Create(ctx context.Context, b *models.Business) error {
	return nil
}

func (repo *businessRepository) Update(ctx context.Context, b *models.Business) error {
	return nil
}

func (repo *businessRepository) GetById(ctx context.Context, id string) (*models.Business, error) {
	return nil, nil
}

func (repo *businessRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *businessRepository) SoftDelete(ctx context.Context, id string) (time.Time, error) {
	return time.Now(), nil
}
