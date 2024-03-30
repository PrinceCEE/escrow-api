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

type IWalletHistoryRepository interface {
	Create(a *models.WalletHistory, tx pgx.Tx) error
	Update(a *models.WalletHistory, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.WalletHistory, error)
	GetByWalletId(id string, tx pgx.Tx) ([]*models.WalletHistory, error)
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

func (repo *WalletHistoryRepository) GetById(id string, tx pgx.Tx) (*models.WalletHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	h := new(models.WalletHistory)
	query := `
		SELECT
			id,
			wallet_id,
			type,
			amount,
			status,
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
		&h.ID,
		&h.WalletID,
		&h.Type,
		&h.Amount,
		&h.Status,
		&h.CreatedAt,
		&h.UpdatedAt,
		&h.DeletedAt,
		&h.Version,
	)
	if err != nil {
		return nil, err
	}

	return h, nil
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

func (repo *WalletHistoryRepository) GetByWalletId(id string, tx pgx.Tx) ([]*models.WalletHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		SELECT
			id,
			wallet_id,
			type,
			amount,
			status,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			wallet_histories
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

	walletHistories, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.WalletHistory])
	if err != nil {
		return nil, err
	}

	return walletHistories, nil
}
