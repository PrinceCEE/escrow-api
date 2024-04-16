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

type ITransactionTimelineRepository interface {
	Create(tt *models.TransactionTimeline, tx pgx.Tx) error
	Update(tt *models.TransactionTimeline, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.TransactionTimeline, error)
	GetMany(args []any, where string, tx pgx.Tx) ([]*models.TransactionTimeline, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type TransactionTimelineRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewTransactionTimelineRepository(db *pgxpool.Pool, timeout time.Duration) *TransactionTimelineRepository {
	return &TransactionTimelineRepository{DB: db, Timeout: timeout}
}

func (repo *TransactionTimelineRepository) Create(tt *models.TransactionTimeline, tx pgx.Tx) error {
	now := time.Now().UTC()
	tt.CreatedAt = now
	tt.UpdatedAt = now

	args := []any{
		tt.Name,
		tt.TransactionID,
		tt.CreatedAt,
		tt.UpdatedAt,
	}

	query := `INSERT INTO transaction_timelines (name, transaction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, version`

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &tt.Version)
		if err != nil {
			return err
		}

		tt.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &tt.Version)
	if err != nil {
		return err
	}

	tt.ID = id.String()
	return nil
}

func (repo *TransactionTimelineRepository) Update(tt *models.TransactionTimeline, tx pgx.Tx) error {
	tt.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(tt, "transaction_timelines")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&tt.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&tt.Version)
}

func (repo *TransactionTimelineRepository) getByKey(key string, value any, tx pgx.Tx) (*models.TransactionTimeline, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	tt := new(models.TransactionTimeline)
	var id, transactionId uuid.UUID

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			transaction_id,
			created_at,
			updated_at,
			deleted_at

		FROM transaction_timelines
		WHERE %s = $1`,
		key,
	)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&id,
		&tt.Name,
		&transactionId,
		&tt.CreatedAt,
		&tt.UpdatedAt,
		&tt.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	tt.ID = id.String()
	tt.TransactionID = transactionId.String()

	return tt, nil
}

func (repo *TransactionTimelineRepository) GetById(id string, tx pgx.Tx) (*models.TransactionTimeline, error) {
	return repo.getByKey("id", id, tx)
}

func (repo *TransactionTimelineRepository) GetMany(args []any, where string, tx pgx.Tx) ([]*models.TransactionTimeline, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			transaction_id,
			created_at,
			updated_at,
			deleted_at
		%s
		ORDER BY created_at
	`, where)

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

	timelines := []*models.TransactionTimeline{}

	for rows.Next() {
		var tt *models.TransactionTimeline
		var id, transactionId uuid.UUID

		err := rows.Scan(
			&id,
			&tt.Name,
			&transactionId,
			&tt.CreatedAt,
			&tt.UpdatedAt,
			&tt.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		tt.ID = id.String()
		tt.TransactionID = transactionId.String()

		timelines = append(timelines, tt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timelines, nil
}

func (repo *TransactionTimelineRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM transaction_timelines WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query, id)
	} else {
		_, err = repo.DB.Exec(ctx, query, id)
	}

	return
}

func (repo *TransactionTimelineRepository) SoftDelete(id string, tx pgx.Tx) error {
	tt, err := repo.GetById(id, tx)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()
	tt.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}
	tt.UpdatedAt = now
	return repo.Update(tt, tx)
}
