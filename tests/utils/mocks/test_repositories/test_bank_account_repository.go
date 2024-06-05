package test_repositories

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/internal/repositories"
	"github.com/princecee/escrow-api/pkg/utils"
	"github.com/stretchr/testify/mock"
)

type TestBankAccountRepository struct {
	repo *repositories.BankAccountRepository
	mock.Mock
}

func NewBankAccountRepository(db *pgxpool.Pool, timeout time.Duration) *TestBankAccountRepository {
	return &TestBankAccountRepository{repo: repositories.NewBankAccountRepository(db, timeout)}
}

func (r *TestBankAccountRepository) Create(a *models.BankAccount, tx pgx.Tx) error {
	return r.repo.Create(a, tx)
}

func (r *TestBankAccountRepository) Update(a *models.BankAccount, tx pgx.Tx) error {
	return r.repo.Update(a, tx)
}

func (r *TestBankAccountRepository) GetById(id string, tx pgx.Tx) (*models.BankAccount, error) {
	return r.repo.GetById(id, tx)
}

func (r *TestBankAccountRepository) Delete(id string, tx pgx.Tx) (err error) {
	return r.repo.Delete(id, tx)
}

func (r *TestBankAccountRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}

func (r *TestBankAccountRepository) GetByWalletId(id string, pagination utils.Pagination, tx pgx.Tx) ([]*models.BankAccount, error) {
	return r.repo.GetByWalletId(id, pagination, tx)
}
