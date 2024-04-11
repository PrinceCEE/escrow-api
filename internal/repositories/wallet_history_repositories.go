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

type IWalletHistoryRepository interface {
	Create(h *models.WalletHistory, tx pgx.Tx) error
	Update(h *models.WalletHistory, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.WalletHistory, error)
	GetByWalletId(id string, pagination utils.Pagination, tx pgx.Tx) ([]*models.WalletHistory, error)
	GetMany(args []any, where string, tx pgx.Tx) ([]*models.WalletHistory, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type WalletHistoryRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewWalletHistoryRepository(db *pgxpool.Pool, timeout time.Duration) *WalletHistoryRepository {
	return &WalletHistoryRepository{DB: db, Timeout: timeout}
}

func (repo *WalletHistoryRepository) Create(h *models.WalletHistory, tx pgx.Tx) error {
	now := time.Now().UTC()
	h.CreatedAt = now
	h.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO wallet_histories (wallet_id, type, amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version
	`

	args := []any{h.WalletID, h.Type, h.Amount, h.Status, h.CreatedAt, h.UpdatedAt}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&h.ID, &h.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(&h.ID, &h.Version)
}

func (repo *WalletHistoryRepository) Update(h *models.WalletHistory, tx pgx.Tx) error {
	h.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(h, "wallet_histories")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&h.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&h.Version)
}

func (repo *WalletHistoryRepository) getByKey(key string, value any, tx pgx.Tx) (*models.WalletHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	h := new(models.WalletHistory)
	query := fmt.Sprintf(`
		SELECT
			h.id,
			h.wallet_id,
			h.type,
			h.amount,
			h.status,
			h.created_at,
			h.updated_at,
			h.deleted_at,
			h.version,
			w.identifier,
			w.balance,
			w.receivable_balance,
			w.payable_balance,
			w.account_type,
			w.created_at,
			w.updated_at,
			w.deleted_at,
			w.version
		FROM
			wallet_histories h
		INNER JOIN wallets w ON w.id = h.wallet_id
		WHERE %s = $1
	`, key)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	var id, walletId uuid.UUID
	var wallet models.Wallet
	err := row.Scan(
		&id,
		&walletId,
		&h.Type,
		&h.Amount,
		&h.Status,
		&h.CreatedAt,
		&h.UpdatedAt,
		&h.DeletedAt,
		&h.Version,
		&wallet.Identifier,
		&wallet.Balance,
		&wallet.Receivable,
		&wallet.Payable,
		&wallet.AccountType,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
		&wallet.DeletedAt,
		&wallet.Version,
	)
	if err != nil {
		return nil, err
	}

	h.ID = id.String()
	h.WalletID = walletId.String()
	wallet.ID = walletId.String()
	return h, nil
}

func (repo *WalletHistoryRepository) GetById(id string, tx pgx.Tx) (*models.WalletHistory, error) {
	return repo.getByKey("h.id", id, tx)
}

func (repo *WalletHistoryRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM wallet_histories WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *WalletHistoryRepository) SoftDelete(id string, tx pgx.Tx) error {
	h, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	h.UpdatedAt = now
	h.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(h, tx)
}

func (repo *WalletHistoryRepository) GetByWalletId(id string, pagination utils.Pagination, tx pgx.Tx) ([]*models.WalletHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			h.id,
			h.wallet_id,
			h.type,
			h.amount,
			h.status,
			h.created_at,
			h.updated_at,
			h.deleted_at,
			h.version
			w.identifier,
			w.balance,
			w.receivable_balance,
			w.payable_balance,
			w.account_type,
			w.created_at,
			w.updated_at,
			w.deleted_at,
			w.version
		FROM
			wallet_histories h
		INNER JOIN wallets w ON w.id = h.wallet_id
		WHERE h.wallet_id = $1
		OFFSET $2
		LIMIT $3
	`

	var rows pgx.Rows
	args := []any{id, pagination.Offset, pagination.Limit}
	if tx != nil {
		_rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rows = _rows
	} else {
		_rows, err := repo.DB.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rows = _rows
	}

	return returnFromRows(rows)
}

func (repo *WalletHistoryRepository) GetMany(args []any, where string, tx pgx.Tx) ([]*models.WalletHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	argLen := len(args)
	query := fmt.Sprintf(`
		SELECT
			h.id,
			h.wallet_id,
			h.type,
			h.amount,
			h.status,
			h.created_at,
			h.updated_at,
			h.deleted_at,
			h.version,
			w.identifier,
			w.balance,
			w.receivable_balance,
			w.payable_balance,
			w.account_type,
			w.created_at,
			w.updated_at,
			w.deleted_at,
			w.version
		FROM
			wallet_histories h
		INNER JOIN wallets w ON w.id = h.wallet_id
		%s
		OFFSET $%d
		LIMIT $%d
	`, where, argLen-1, argLen)

	var rows pgx.Rows
	if tx != nil {
		_rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rows = _rows
	} else {
		_rows, err := repo.DB.Query(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rows = _rows
	}

	return returnFromRows(rows)
}

func returnFromRows(rows pgx.Rows) ([]*models.WalletHistory, error) {
	walletHistories := []*models.WalletHistory{}

	for rows.Next() {
		var id, walletId uuid.UUID
		var h models.WalletHistory
		var wallet models.Wallet

		err := rows.Scan(
			&id,
			&walletId,
			&h.Type,
			&h.Amount,
			&h.Status,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.DeletedAt,
			&h.Version,
			&wallet.Identifier,
			&wallet.Balance,
			&wallet.Receivable,
			&wallet.Payable,
			&wallet.AccountType,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
			&wallet.DeletedAt,
			&wallet.Version,
		)

		if err != nil {
			return nil, err
		}

		h.ID = id.String()
		h.Wallet = wallet
		h.WalletID = walletId.String()
		wallet.ID = walletId.String()
		walletHistories = append(walletHistories, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return walletHistories, nil
}
