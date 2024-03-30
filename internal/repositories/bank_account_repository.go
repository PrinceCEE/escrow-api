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

type IBankAccountRepository interface {
	Create(a *models.BankAccount, tx pgx.Tx) error
	Update(a *models.BankAccount, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.BankAccount, error)
	GetByWalletId(id string, tx pgx.Tx) ([]*models.BankAccount, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type BankAccountRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewBankAccountRepository(db *pgxpool.Pool, timeout time.Duration) *BankAccountRepository {
	return &BankAccountRepository{DB: db, Timeout: timeout}
}

func (repo *BankAccountRepository) Create(a *models.BankAccount, tx pgx.Tx) error {
	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO bank_accounts (wallet_id, bank_name, account_name, account_number, bvn, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, version
	`

	args := []any{a.WalletID, a.BankName, a.AccountName, a.AccountNumber, a.BVN, a.CreatedAt, a.UpdatedAt}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&a.ID, &a.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(&a.ID, &a.Version)
}

func (repo *BankAccountRepository) Update(a *models.BankAccount, tx pgx.Tx) error {
	a.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(a, "bank_accounts")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&a.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&a.Version)
}

func (repo *BankAccountRepository) GetById(id string, tx pgx.Tx) (*models.BankAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	a := new(models.BankAccount)
	query := `
		SELECT
			id,
			wallet_id,
			bank_name,
			account_name,
			account_number,
			bvn,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			bank_accounts
		WHERE id = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		&a.ID,
		&a.WalletID,
		&a.BankName,
		&a.AccountName,
		&a.AccountNumber,
		&a.BVN,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.Version,
	)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (repo *BankAccountRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM bank_accounts WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *BankAccountRepository) SoftDelete(id string, tx pgx.Tx) error {
	a, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	a.UpdatedAt = now
	a.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(a, tx)
}

func (repo *BankAccountRepository) GetByWalletId(id string, tx pgx.Tx) ([]*models.BankAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			id,
			wallet_id,
			bank_name,
			account_name,
			account_number,
			bvn,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			bank_accounts
		WHERE wallet_id = $1
	`

	var rows pgx.Rows
	if tx != nil {
		_rows, err := tx.Query(ctx, query, id)
		if err != nil {
			return nil, err
		}

		rows = _rows
	} else {
		_rows, err := repo.DB.Query(ctx, query, id)
		if err != nil {
			return nil, err
		}

		rows = _rows
	}

	bankAccounts, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.BankAccount])
	if err != nil {
		return nil, err
	}

	return bankAccounts, nil
}
