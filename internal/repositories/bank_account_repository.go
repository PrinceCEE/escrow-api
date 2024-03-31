package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IBankAccountRepository interface {
	Create(b *models.BankAccount, tx pgx.Tx) error
	Update(b *models.BankAccount, tx pgx.Tx) error
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

func (repo *BankAccountRepository) Create(b *models.BankAccount, tx pgx.Tx) error {
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO bank_accounts (wallet_id, bank_name, account_name, account_number, bvn, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, version
	`

	args := []any{b.WalletID, b.BankName, b.AccountName, b.AccountNumber, b.BVN, b.CreatedAt, b.UpdatedAt}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&b.ID, &b.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(&b.ID, &b.Version)
}

func (repo *BankAccountRepository) Update(b *models.BankAccount, tx pgx.Tx) error {
	b.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(b, "bank_accounts")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&b.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&b.Version)
}

func (repo *BankAccountRepository) getByKey(key string, value any, tx pgx.Tx) (*models.BankAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	b := new(models.BankAccount)
	query := fmt.Sprintf(`
		SELECT
			b.id,
			b.wallet_id,
			b.bank_name,
			b.account_name,
			b.account_number,
			b.bvn,
			b.created_at,
			b.updated_at,
			b.deleted_at,
			b.version
			w.id,
			w.balance,
			w.receivable_balance,
			w.payable_balance,
			w.account_type,
			w.identifier,
			w.created_at,
			w.updated_at,
			w.deleted_at,
			w.version,
		FROM
			bank_accounts b
		INNER JOIN wallets w ON w.id = b.wallet_id
		WHERE %s = $1
	`, key)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&b.ID,
		&b.WalletID,
		&b.BankName,
		&b.AccountName,
		&b.AccountNumber,
		&b.BVN,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
		&b.Version,
		&b.Wallet.ID,
		&b.Wallet.Balance,
		&b.Wallet.Receivable,
		&b.Wallet.Payable,
		&b.Wallet.AccountType,
		&b.Wallet.Identifier,
		&b.Wallet.CreatedAt,
		&b.Wallet.UpdatedAt,
		&b.Wallet.DeletedAt,
		&b.Wallet.Version,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (repo *BankAccountRepository) GetById(id string, tx pgx.Tx) (*models.BankAccount, error) {
	return repo.getByKey("b.id", id, tx)
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
	b, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	b.UpdatedAt = now
	b.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(b, tx)
}

func (repo *BankAccountRepository) GetByWalletId(id string, tx pgx.Tx) ([]*models.BankAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			b.id,
			b.wallet_id,
			b.bank_name,
			b.account_name,
			b.account_number,
			b.bvn,
			b.created_at,
			b.updated_at,
			b.deleted_at,
			b.version
			w.id,
			w.balance,
			w.receivable_balance,
			w.payable_balance,
			w.account_type,
			w.identifier,
			w.created_at,
			w.updated_at,
			w.deleted_at,
			w.version,
		FROM
			bank_accounts b
		INNER JOIN wallets w ON w.id = b.wallet_id
		WHERE b.wallet_id = $1
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

	bankAccounts := []*models.BankAccount{}
	for rows.Next() {
		var b models.BankAccount
		err := rows.Scan(
			&b.ID,
			&b.WalletID,
			&b.BankName,
			&b.AccountName,
			&b.AccountNumber,
			&b.BVN,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.DeletedAt,
			&b.Version,
			&b.Wallet.ID,
			&b.Wallet.Balance,
			&b.Wallet.Receivable,
			&b.Wallet.Payable,
			&b.Wallet.AccountType,
			&b.Wallet.Identifier,
			&b.Wallet.CreatedAt,
			&b.Wallet.UpdatedAt,
			&b.Wallet.DeletedAt,
			&b.Wallet.Version,
		)

		if err != nil {
			return nil, err
		}

		bankAccounts = append(bankAccounts, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bankAccounts, nil
}
