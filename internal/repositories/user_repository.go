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

type IUserRepository interface {
	Create(u *models.User, tx pgx.Tx) error
	Update(u *models.User, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.User, error)
	GetByEmail(email string, tx pgx.Tx) (*models.User, error)
	GetByPhoneNumber(phone string, tx pgx.Tx) (*models.User, error)
	GetByBusinessId(id string, tx pgx.Tx) (*models.User, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type UserRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewUserRepository(db *pgxpool.Pool, timeout time.Duration) *UserRepository {
	return &UserRepository{DB: db, Timeout: timeout}
}

func (repo *UserRepository) Create(u *models.User, tx pgx.Tx) error {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	query := `
		INSERT INTO users (email, account_type, reg_stage, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version
	`

	args := []any{u.Email, u.AccountType, u.RegStage, u.CreatedAt, u.UpdatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var id uuid.UUID
	if tx != nil {
		err := tx.QueryRow(ctx, query, args...).Scan(&id, &u.Version)
		if err != nil {
			return err
		}

		u.ID = id.String()
		return nil
	}

	err := repo.DB.QueryRow(ctx, query, args...).Scan(&id, &u.Version)
	if err != nil {
		return err
	}

	u.ID = id.String()
	return nil
}

func (repo *UserRepository) Update(u *models.User, tx pgx.Tx) error {
	u.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	qs, err := utils.GetUpdateQueryFromStruct(u, "users")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, qs.Query, qs.Args...).Scan(&u.Version)
	}

	return repo.DB.QueryRow(ctx, qs.Query, qs.Args...).Scan(&u.Version)
}

func (repo *UserRepository) getByKey(key string, value any, tx pgx.Tx) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	u := new(models.User)

	var id, businessId *uuid.UUID
	var imageUrl *string
	var business models.Business

	query := fmt.Sprintf(`
		SELECT
			u.id,
			u.email,
			u.phone_number,
			u.first_name,
			u.last_name,
			u.is_phone_number_verified,
			u.is_email_verified,
			u.reg_stage,
			u.account_type,
			u.business_id,
			u.image_url,
			u.created_at,
			u.updated_at,
			u.deleted_at,
			u.version,
			b.name,
			b.email,
			b.image_url,
			b.created_at,
			b.updated_at,
			b.deleted_at,
			b.version
		FROM
			users u
		LEFT JOIN businesses b ON b.id = u.business_id
		WHERE
			%s = $1
			AND u.deleted_at IS NULL`,
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
		&u.Email,
		&u.PhoneNumber,
		&u.FirstName,
		&u.LastName,
		&u.IsPhoneNumberVerified,
		&u.IsEmailVerified,
		&u.RegStage,
		&u.AccountType,
		&businessId,
		&imageUrl,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
		&u.Version,
		&business.Name,
		&business.Email,
		&business.ImageUrl,
		&business.CreatedAt,
		&business.UpdatedAt,
		&business.DeletedAt,
		&business.Version,
	)
	if err != nil {
		return nil, err
	}

	u.ID = id.String()
	u.ImageUrl = *imageUrl
	if u.AccountType == models.BusinessAccountType {
		u.BusinessID = businessId.String()
		u.Business = &business
	}

	return u, nil
}

func (repo *UserRepository) GetById(id string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("id", id, tx)
}

func (repo *UserRepository) GetByEmail(email string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("u.email", email, tx)
}

func (repo *UserRepository) GetByPhoneNumber(phone string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("u.phone_number", phone, tx)
}

func (repo *UserRepository) GetByBusinessId(id string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("u.business_id", id, tx)
}

func (repo *UserRepository) Delete(id string, tx pgx.Tx) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	if tx != nil {
		_, err = tx.Exec(ctx, query, id)
	} else {
		_, err = repo.DB.Exec(ctx, query, id)
	}

	return
}

func (repo *UserRepository) SoftDelete(id string, tx pgx.Tx) error {
	u, err := repo.GetById(id, tx)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()
	u.DeletedAt = models.NullTime{NullTime: sql.NullTime{Time: now}}
	u.UpdatedAt = now
	return repo.Update(u, tx)
}
