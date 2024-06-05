package test_repositories

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/internal/repositories"
	"github.com/stretchr/testify/mock"
)

type TestTransactionRepository struct {
	repo *repositories.TransactionRepository
	mock.Mock
}

func NewTransactionRepository(db *pgxpool.Pool, timeout time.Duration) *TestTransactionRepository {
	return &TestTransactionRepository{repo: repositories.NewTransactionRepository(db, timeout)}
}

func (r *TestTransactionRepository) Create(t *models.Transaction, tx pgx.Tx) error {
	return r.repo.Create(t, tx)
}

func (r *TestTransactionRepository) Update(t *models.Transaction, tx pgx.Tx) error {
	return r.repo.Update(t, tx)
}

func (r *TestTransactionRepository) GetById(id string, tx pgx.Tx) (*models.Transaction, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestTransactionRepository) GetMany(args []any, where string, tx pgx.Tx) ([]*models.Transaction, error) {
	return r.repo.GetMany(args, where, tx)
}

func (r *TestTransactionRepository) Delete(id string, tx pgx.Tx) (err error) {
	return r.repo.Delete(id, tx)
}

func (r *TestTransactionRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
