package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(e *models.Event) error
	Update(e *models.Event) error
	GetById(id string) (*models.Event, error)
	Delete(id string) error
	SoftDelete(id string) error
}

type eventRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewEventRepository(db *pgxpool.Pool, timeout time.Duration) *eventRepository {
	return &eventRepository{DB: db, Timeout: timeout}
}

func (repo *eventRepository) Create(e *models.Event) error {
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

	return repo.DB.QueryRow(ctx, query, args...).Scan(e.ID, e.Version)
}

func (repo *eventRepository) Update(e *models.Event) error {
	e.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(e, "events")
	if err != nil {
		return err
	}

	return repo.DB.QueryRow(ctx, query).Scan(e.Version)
}

func (repo *eventRepository) GetById(id string) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var e *models.Event
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

	err := repo.DB.QueryRow(ctx, query, id).Scan(
		e.ID,
		e.Data,
		e.OriginEnvironment,
		e.TargetEnvironment,
		e.EventType,
		e.CreatedAt,
		e.UpdatedAt,
		e.DeletedAt,
		e.Version,
	)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (repo *eventRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`

	_, err := repo.DB.Exec(ctx, query)
	return err
}

func (repo *eventRepository) SoftDelete(id string) error {
	e, err := repo.GetById(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	e.UpdatedAt = now
	e.DeletedAt = now

	return repo.Update(e)
}
