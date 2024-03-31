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

type IWalletRepository interface {
	Create(w *models.Wallet, tx pgx.Tx) error
	Update(w *models.Wallet, tx pgx.Tx) error
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

func (repo *WalletRepository) getByKey(key string, value any, tx pgx.Tx) (*models.Wallet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	w := new(models.Wallet)
	query := fmt.Sprintf(`
		SELECT
			w.id,
			w.identifier,
			w.balance,
			w.receivable_balance,
			w.payable_balance,
			w.account_type,
			w.created_at,
			w.updated_at,
			w.deleted_at,
			w.version
			u.id,
			u.email,
			u.phone_number,
			u.first_name,
			u.last_name,
			u.is_phone_number_verified,
			u.is_email_verified,
			u.reg_stage,
			u.account_type,
			u.business_id,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			u.version
			b.id,
			b.name,
			b.email,
			b.created_at,
			b.updated_at,
			b.deleted_at,
			b.version
		FROM
			wallets w
		INNER JOIN users u ON u.id = w.identifier
		INNER JOIN businesses b ON b.id = w.identifier
		WHERE %s = $1
	`, key)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
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
		&w.User.ID,
		&w.User.Email,
		&w.User.PhoneNumber,
		&w.User.FirstName,
		&w.User.LastName,
		&w.User.IsPhoneNumberVerified,
		&w.User.IsEmailVerified,
		&w.User.RegStage,
		&w.User.AccountType,
		&w.User.BusinessID,
		&w.User.CreatedAt,
		&w.User.UpdatedAt,
		&w.User.DeletedAt,
		&w.User.Version,
		&w.Business.ID,
		&w.Business.Name,
		&w.Business.Email,
		&w.Business.CreatedAt,
		&w.Business.UpdatedAt,
		&w.Business.DeletedAt,
		&w.Business.Version,
	)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (repo *WalletRepository) GetById(id string, tx pgx.Tx) (*models.Wallet, error) {
	return repo.getByKey("w.id", id, tx)
}

func (repo *WalletRepository) GetByIdentifier(id string, tx pgx.Tx) (*models.Wallet, error) {
	return repo.getByKey("w.identifier", id, tx)
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
