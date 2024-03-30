package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type IWalletHistoryRepository interface {
	Create(a *models.WalletHistory, tx pgx.Tx) error
	Update(a *models.WalletHistory, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.WalletHistory, error)
	GetByWalletId(id string, tx pgx.Tx) ([]*models.WalletHistory, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type TestWalletHistoryRepository struct {
	repo *repositories.WalletHistoryRepository
	mock.Mock
}

func NewWalletHistoryRepository(db *pgxpool.Pool, timeout time.Duration) *TestWalletHistoryRepository {
	return &TestWalletHistoryRepository{repo: repositories.NewWalletHistoryRepository(db, timeout)}
}

func (r *TestWalletHistoryRepository) Create(h *models.WalletHistory, tx pgx.Tx) error {
	return r.repo.Create(h, tx)
}

func (r *TestWalletHistoryRepository) Update(h *models.WalletHistory, tx pgx.Tx) error {
	return r.repo.Update(h, tx)
}

func (r *TestWalletHistoryRepository) GetById(id string, tx pgx.Tx) (*models.WalletHistory, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestWalletHistoryRepository) Delete(id string, tx pgx.Tx) (err error) {
	return r.repo.Delete(id, tx)
}

func (r *TestWalletHistoryRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}

func (r *TestWalletHistoryRepository) GetByWalletId(id string, tx pgx.Tx) ([]*models.WalletHistory, error) {
	return r.repo.GetByWalletId(id, tx)
}
