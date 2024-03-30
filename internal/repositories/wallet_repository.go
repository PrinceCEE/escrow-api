package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IWalletRepository interface {
	Create(a *models.Wallet, tx pgx.Tx) error
	Update(a *models.Wallet, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Wallet, error)
	GetByIdentifier(id string, tx pgx.Tx) (*models.Wallet, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type WalletRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewWalletRepository(db *pgxpool.Pool, timeout time.Duration) *WalletRepository {
	return &WalletRepository{DB: db, Timeout: timeout}
}

func (repo *WalletRepository) Create(w *models.Wallet, tx pgx.Tx) error {
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO wallets (identifier, balance, receivable_balance, payable_balance, account_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, version
	`

	args := []any{w.Identifier, w.Balance, w.Receivable, w.Payable, w.AccountType, w.CreatedAt, w.UpdatedAt}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&w.ID, &w.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(&w.ID, &w.Version)
}

func (repo *WalletRepository) Update(w *models.Wallet, tx pgx.Tx) error {
	w.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(w, "wallets")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&w.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&w.Version)
}

func (repo *WalletRepository) GetById(id string, tx pgx.Tx) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	w := new(models.Wallet)
	query := `
		SELECT
			id,
			identifier,
			balance,
			receivable_balance,
			payable_balance,
			account_type,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			wallets
		WHERE id = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		&w.ID,
		&w.Identifier,
		&w.Balance,
		&w.Receivable,
		&w.Payable,
		&w.AccountType,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.DeletedAt,
		&w.Version,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (repo *WalletRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM wallets WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *WalletRepository) SoftDelete(id string, tx pgx.Tx) error {
	w, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	w.UpdatedAt = now
	w.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(w, tx)
}

func (repo *WalletRepository) GetByIdentifier(id string, tx pgx.Tx) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	w := new(models.Wallet)
	query := `
		SELECT
			id,
			identifier,
			balance,
			receivable_balance,
			payable_balance,
			account_type,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			wallets
		WHERE identifier = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		&w.ID,
		&w.Identifier,
		&w.Balance,
		&w.Receivable,
		&w.Payable,
		&w.AccountType,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.DeletedAt,
		&w.Version,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}
