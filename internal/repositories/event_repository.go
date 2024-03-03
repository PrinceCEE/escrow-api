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

type EventRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewEventRepository(db *pgxpool.Pool, timeout time.Duration) *EventRepository {
	return &EventRepository{DB: db, Timeout: timeout}
}

func (repo *EventRepository) Create(e *models.Event, tx pgx.Tx) error {
	now := time.Now().UTC()
	e.CreatedAt = now
	e.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO events (data, origin_environment, target_environment, event_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version
	`

	args := []any{
		e.Data,
		e.OriginEnvironment,
		e.TargetEnvironment,
		e.EventType,
		e.CreatedAt,
		e.UpdatedAt,
	}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(&e.ID, &e.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(&e.ID, &e.Version)
}

func (repo *EventRepository) Update(e *models.Event, tx pgx.Tx) error {
	e.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(e, "events")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, query).Scan(&e.Version)
	}

	return repo.DB.QueryRow(ctx, query).Scan(&e.Version)
}

func (repo *EventRepository) GetById(id string, tx pgx.Tx) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	e := new(models.Event)
	query := `
		SELECT
			id,
			data,
			origin_environment,
			target_environment,
			event_type
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			events
		WHERE id = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		&e.ID,
		&e.Data,
		&e.OriginEnvironment,
		&e.TargetEnvironment,
		&e.EventType,
		&e.CreatedAt,
		&e.UpdatedAt,
		&e.DeletedAt,
		&e.Version,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (repo *EventRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query, id)
	} else {
		_, err = repo.DB.Exec(ctx, query, id)
	}
	return err
}

func (repo *EventRepository) SoftDelete(id string, tx pgx.Tx) error {
	e, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	e.UpdatedAt = now
	e.DeletedAt = sql.NullTime{Time: now}

	return repo.Update(e, tx)
}
