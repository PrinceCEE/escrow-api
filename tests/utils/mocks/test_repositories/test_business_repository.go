package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type TestBusinessRepository struct {
	repo *repositories.BusinessRepository
	mock.Mock
}

func NewBusinessRepository(db *pgxpool.Pool, timeout time.Duration) *TestBusinessRepository {
	return &TestBusinessRepository{repo: repositories.NewBusinessRepository(db, timeout)}
}

func (r *TestBusinessRepository) Create(b *models.Business, tx pgx.Tx) error {
	return r.repo.Create(b, tx)
}

func (r *TestBusinessRepository) Update(b *models.Business, tx pgx.Tx) error {
	return r.repo.Update(b, tx)
}

func (r *TestBusinessRepository) GetById(id string, tx pgx.Tx) (*models.Business, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestBusinessRepository) Delete(id string, tx pgx.Tx) error {
	return r.repo.Delete(id, tx)
}

func (r *TestBusinessRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
