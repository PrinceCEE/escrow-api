package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/pkg/utils"
)

type IOtpRepository interface {
	Create(otp *models.Otp, tx pgx.Tx) error
	Update(otp *models.Otp, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.Otp, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
	GetOneByWhere(where string, args []any, tx pgx.Tx) (*models.Otp, error)
}

type OtpRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewOtpRepository(db *pgxpool.Pool, timeout time.Duration) *OtpRepository {
	return &OtpRepository{DB: db, Timeout: timeout}
}

func (repo *OtpRepository) Create(otp *models.Otp, tx pgx.Tx) error {
	now := time.Now().UTC()
	otp.CreatedAt = now
	otp.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		INSERT INTO otps (user_id, code, is_used, otp_type, expires_in, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, version
	`

	args := []any{
		otp.UserID,
		otp.Code,
		otp.IsUsed,
		otp.OtpType,
		otp.ExpiresIn,
		otp.CreatedAt,
		otp.UpdatedAt,
	}

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &otp.Version)
		if err != nil {
			return err
		}

		otp.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &otp.Version)
	if err != nil {
		return err
	}

	otp.ID = id.String()
	return nil
}

func (repo *OtpRepository) Update(otp *models.Otp, tx pgx.Tx) error {
	otp.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(otp, "otps")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&otp.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&otp.Version)
}

func (repo *OtpRepository) getByKey(key string, value any, tx pgx.Tx) (*models.Otp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var id, userId uuid.UUID
	otp := new(models.Otp)
	query := `
		SELECT
			id,
			user_id,
			code,
			is_used,
			otp_type,
			expires_in,
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
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		&id,
		&userId,
		&otp.Code,
		&otp.IsUsed,
		&otp.OtpType,
		&otp.ExpiresIn,
		&otp.CreatedAt,
		&otp.UpdatedAt,
		&otp.DeletedAt,
		&otp.Version,
	)
	if err != nil {
		return nil, err
	}

	otp.ID = id.String()
	otp.UserID = userId.String()

	return otp, nil
}

func (repo *OtpRepository) GetById(id string, tx pgx.Tx) (*models.Otp, error) {
	return repo.getByKey("id", id, tx)
}

func (repo *OtpRepository) GetOneByWhere(where string, args []any, tx pgx.Tx) (*models.Otp, error) {
	otp := new(models.Otp)

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			code,
			is_used,
			otp_type,
			expires_in,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			otps
		%s
	`, where)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, args...)
	} else {
		row = repo.DB.QueryRow(ctx, query, args...)
	}

	var id, userId uuid.UUID
	err := row.Scan(
		&id,
		&userId,
		&otp.Code,
		&otp.IsUsed,
		&otp.OtpType,
		&otp.ExpiresIn,
		&otp.CreatedAt,
		&otp.UpdatedAt,
		&otp.DeletedAt,
		&otp.Version,
	)

	if err != nil {
		return nil, err
	}

	otp.ID = id.String()
	otp.UserID = userId.String()

	return otp, nil
}

func (repo *OtpRepository) Delete(id string, tx pgx.Tx) (err error) {
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

func (repo *OtpRepository) SoftDelete(id string, tx pgx.Tx) error {
	otp, err := repo.GetById(id, tx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	otp.UpdatedAt = now
	otp.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}

	return repo.Update(otp, tx)
}
