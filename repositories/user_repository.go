package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	Update(ctx context.Context, u *models.User) error
	GetById(ctx context.Context, id string) (*models.User, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) (time.Time, error)
}

type userRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) Create(ctx context.Context, u *models.User) error {
	return nil
}

func (repo *userRepository) Update(ctx context.Context, u *models.User) error {
	return nil
}

func (repo *userRepository) GetById(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

func (repo *userRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *userRepository) SoftDelete(ctx context.Context, id string) (time.Time, error) {
	return time.Now(), nil
}
