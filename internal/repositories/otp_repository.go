package repositories

import (
	"context"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OtpRepository interface {
	Create(otp *models.Otp, tx pgx.Tx) error
	Update(otp *models.Otp, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Otp, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type otpRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewOtpRepository(db *pgxpool.Pool, timeout time.Duration) *otpRepository {
	return &otpRepository{DB: db, Timeout: timeout}
}

func (repo *otpRepository) Create(otp *models.Otp, tx pgx.Tx) error {
	now := time.Now().UTC()
	otp.CreatedAt = now
	otp.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO otps (user_id, code, is_used, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version
	`

	args := []any{otp.UserID, otp.Code, otp.IsUsed, otp.CreatedAt, otp.UpdatedAt}

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(otp.ID, otp.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(otp.ID, otp.Version)
}

func (repo *otpRepository) Update(otp *models.Otp, tx pgx.Tx) error {
	otp.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(otp, "otps")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, query).Scan(otp.Version)
	}

	return repo.DB.QueryRow(ctx, query).Scan(otp.Version)
}

func (repo *otpRepository) GetById(id string, tx pgx.Tx) (*models.Otp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	otp := new(models.Otp)
	query := `
		SELECT
			id,
			user_id,
			code,
			is_used,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			otps
		WHERE id = $1
	`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, id)
	} else {
		row = repo.DB.QueryRow(ctx, query, id)
	}

	err := row.Scan(
		otp.ID,
		otp.UserID,
		otp.Code,
		otp.IsUsed,
		otp.CreatedAt,
		otp.UpdatedAt,
		otp.DeletedAt,
		otp.Version,
	)
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (repo *otpRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM otps WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query)
	} else {
		_, err = repo.DB.Exec(ctx, query)
	}

	return
}

func (repo *otpRepository) SoftDelete(id string, tx pgx.Tx) error {
	a, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	a.UpdatedAt = now
	a.DeletedAt = now

	return repo.Update(a, tx)
}
