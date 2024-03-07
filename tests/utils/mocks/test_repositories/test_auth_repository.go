package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type TestAuthRepository struct {
	repo *repositories.AuthRepository
	mock.Mock
}

func NewAuthRepository(db *pgxpool.Pool, timeout time.Duration) *TestAuthRepository {
	return &TestAuthRepository{repo: repositories.NewAuthRepository(db, timeout)}
}

func (r *TestAuthRepository) Create(a *models.Auth, tx pgx.Tx) error {
	return r.repo.Create(a, tx)
}

func (r *TestAuthRepository) Update(a *models.Auth, tx pgx.Tx) error {
	return r.repo.Update(a, tx)
}

func (r *TestAuthRepository) GetById(id string, tx pgx.Tx) (*models.Auth, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestAuthRepository) GetByUserId(id string, tx pgx.Tx) (*models.Auth, error) {
	return r.repo.GetByUserId(id, tx)
}

func (r *TestAuthRepository) Delete(id string, tx pgx.Tx) (err error) {
	return r.repo.Delete(id, tx)
}

func (r *TestAuthRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
