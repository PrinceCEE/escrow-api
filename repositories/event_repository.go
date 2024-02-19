package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, e *models.Event) error
	Update(ctx context.Context, e *models.Event) error
	GetById(ctx context.Context, id string) (*models.Event, error)
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) (time.Time, error)
}

type eventRepository struct {
	DB *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *eventRepository {
	return &eventRepository{DB: db}
}

func (repo *eventRepository) Create(ctx context.Context, e *models.Event) error {
	return nil
}

func (repo *eventRepository) Update(ctx context.Context, e *models.Event) error {
	return nil
}

func (repo *eventRepository) GetById(ctx context.Context, id string) (*models.Event, error) {
	return nil, nil
}

func (repo *eventRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (repo *eventRepository) SoftDelete(ctx context.Context, id string) (time.Time, error) {
	return time.Now(), nil
}
