package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type TestWalletRepository struct {
	repo *repositories.WalletRepository
	mock.Mock
}

func NewWalletRepository(db *pgxpool.Pool, timeout time.Duration) *TestWalletRepository {
	return &TestWalletRepository{repo: repositories.NewWalletRepository(db, timeout)}
}

func (r *TestWalletRepository) Create(w *models.Wallet, tx pgx.Tx) error {
	return r.repo.Create(w, tx)
}

func (r *TestWalletRepository) Update(w *models.Wallet, tx pgx.Tx) error {
	return r.repo.Update(w, tx)
}

func (r *TestWalletRepository) GetById(id string, tx pgx.Tx) (*models.Wallet, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestWalletRepository) Delete(id string, tx pgx.Tx) (err error) {
	return r.repo.Delete(id, tx)
}

func (r *TestWalletRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}

func (r *TestWalletRepository) GetByIdentifier(id string, tx pgx.Tx) (*models.Wallet, error) {
	return r.repo.GetByIdentifier(id, tx)
}
