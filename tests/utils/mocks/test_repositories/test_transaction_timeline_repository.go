package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type TestTransactionTimelineRepository struct {
	repo *repositories.TransactionTimelineRepository
	mock.Mock
}

func NewTransactionTimelineRepository(db *pgxpool.Pool, timeout time.Duration) *TestTransactionTimelineRepository {
	return &TestTransactionTimelineRepository{repo: repositories.NewTransactionTimelineRepository(db, timeout)}
}

func (r *TestTransactionTimelineRepository) Create(tt *models.TransactionTimeline, tx pgx.Tx) error {
	return r.repo.Create(tt, tx)
}

func (r *TestTransactionTimelineRepository) Update(tt *models.TransactionTimeline, tx pgx.Tx) error {
	return r.repo.Update(tt, tx)
}

func (r *TestTransactionTimelineRepository) GetById(id string, tx pgx.Tx) (*models.TransactionTimeline, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestTransactionTimelineRepository) GetMany(args []any, where string, tx pgx.Tx) ([]*models.TransactionTimeline, error) {
	return r.repo.GetMany(args, where, tx)
}

func (r *TestTransactionTimelineRepository) Delete(id string, tx pgx.Tx) (err error) {
	return r.repo.Delete(id, tx)
}

func (r *TestTransactionTimelineRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
