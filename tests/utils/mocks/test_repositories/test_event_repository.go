package test_repositories

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/internal/repositories"
	"github.com/stretchr/testify/mock"
)

type TestEventRepository struct {
	repo *repositories.EventRepository
	mock.Mock
}

func NewEventRepository(db *pgxpool.Pool, timeout time.Duration) *TestEventRepository {
	return &TestEventRepository{repo: repositories.NewEventRepository(db, timeout)}
}

func (r *TestEventRepository) Create(e *models.Event, tx pgx.Tx) error {
	return r.repo.Create(e, tx)
}

func (r *TestEventRepository) Update(e *models.Event, tx pgx.Tx) error {
	return r.repo.Update(e, tx)
}

func (r *TestEventRepository) GetById(id string, tx pgx.Tx) (*models.Event, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestEventRepository) Delete(id string, tx pgx.Tx) error {
	return r.repo.Delete(id, tx)
}

func (r *TestEventRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
