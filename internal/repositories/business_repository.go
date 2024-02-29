package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BusinessRepository interface {
	Create(b *models.Business, tx pgx.Tx) error
	Update(b *models.Business, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Business, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type businessRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewBusinessRepository(db *pgxpool.Pool, timeout time.Duration) *businessRepository {
	return &businessRepository{DB: db, Timeout: timeout}
}

func (repo *businessRepository) Create(b *models.Business, tx pgx.Tx) error {
	now := time.Now().UTC()
	b.CreatedAt = now
	b.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO businesses (user_id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version
	`

	args := []any{b.UserID, b.Name, b.Email, b.CreatedAt, b.UpdatedAt}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(b.ID, b.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(b.ID, b.Version)
}

func (repo *businessRepository) Update(b *models.Business, tx pgx.Tx) error {
	b.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(b, "businesses")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, query).Scan(b.Version)
	}

	return repo.DB.QueryRow(ctx, query).Scan(b.Version)
}

func (repo *businessRepository) GetById(id string, tx pgx.Tx) (*models.Business, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var b *models.Business
	query := `
		SELECT
			id,
			user_id,
			name,
			email,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			businesses
		WHERE id = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		b.ID,
		b.UserID,
		b.Name,
		b.Email,
		b.CreatedAt,
		b.UpdatedAt,
		b.DeletedAt,
		b.Version,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (repo *businessRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM businesses WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *businessRepository) SoftDelete(id string, tx pgx.Tx) error {
	b, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	b.UpdatedAt = now
	b.DeletedAt = now

	return repo.Update(b, tx)
}
