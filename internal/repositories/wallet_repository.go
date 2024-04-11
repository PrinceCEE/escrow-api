package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/gofrs/uuid"
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

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &w.Version)
		if err != nil {
			return err
		}

		w.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &w.Version)
	if err != nil {
		return err
	}

	w.ID = id.String()
	return nil
}

func (repo *WalletRepository) Update(w *models.Wallet, tx pgx.Tx) error {
	user := w.User
	business := w.Business

	w.User = nil
	w.Business = nil

	w.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(w, "wallets")
	if err != nil {
		return err
	}

	w.User = user
	w.Business = business

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
			w.version,
			COALESCE(u.email, ''),
			COALESCE(u.phone_number, ''),
			COALESCE(u.first_name, ''),
			COALESCE(u.last_name, ''),
			COALESCE(u.is_phone_number_verified, false),
			COALESCE(u.is_email_verified, false),
			COALESCE(u.reg_stage, 1),
			COALESCE(u.account_type, 'personal'),
			COALESCE(u.business_id, NULL),
			COALESCE(u.created_at, now()),
			COALESCE(u.updated_at, now()),
			COALESCE(u.deleted_at, NULL),
			COALESCE(u.version, 1),
			COALESCE(b.name, ''),
			COALESCE(b.email, ''),
			COALESCE(b.created_at, now()),
			COALESCE(b.updated_at, now()),
			COALESCE(b.deleted_at, NULL),
			COALESCE(b.version, 1)

		FROM wallets w
		LEFT JOIN users u ON u.id = w.identifier AND w.account_type = 'personal'
		LEFT JOIN businesses b ON b.id = w.identifier AND w.account_type = 'business'
		WHERE %s = $1
	`, key)

	var id, identifier uuid.UUID
	user := new(models.User)
	business := new(models.Business)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&id,
		&identifier,
		&w.Balance,
		&w.Receivable,
		&w.Payable,
		&w.AccountType,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.DeletedAt,
		&w.Version,
		&user.Email,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.IsPhoneNumberVerified,
		&user.IsEmailVerified,
		&user.RegStage,
		&user.AccountType,
		&user.BusinessID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Version,
		&business.Name,
		&business.Email,
		&business.CreatedAt,
		&business.UpdatedAt,
		&business.DeletedAt,
		&business.Version,
	)
	if err != nil {
		return nil, err
	}

	w.ID = id.String()
	if w.AccountType == models.PersonalAccountType {
		user.ID = identifier.String()
		w.Identifier = identifier.String()
		w.User = user
	} else {
		business.ID = identifier.String()
		w.Identifier = identifier.String()
		w.Business = business
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
