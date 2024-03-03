package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(u *models.User, tx pgx.Tx) error
	Update(u *models.User, tx pgx.Tx) error
	GetById(id string, tx pgx.Tx) (*models.User, error)
	GetByEmail(email string, tx pgx.Tx) (*models.User, error)
	GetByPhoneNumber(phone string, tx pgx.Tx) (*models.User, error)
	Delete(id string, tx pgx.Tx) error
	SoftDelete(id string, tx pgx.Tx) error
}

type userRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewUserRepository(db *pgxpool.Pool, timeout time.Duration) *userRepository {
	return &userRepository{DB: db, Timeout: timeout}
}

func (repo *userRepository) Create(u *models.User, tx pgx.Tx) error {
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

	if tx != nil {
		return tx.QueryRow(ctx, query, args...).Scan(u.ID, u.Version)
	}

	return repo.DB.QueryRow(ctx, query, args...).Scan(u.ID, u.Version)
}

func (repo *userRepository) Update(u *models.User, tx pgx.Tx) error {
	u.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := utils.GetUpdateQueryFromStruct(u, "users")
	if err != nil {
		return err
	}

	if tx != nil {
		return tx.QueryRow(ctx, query).Scan(u.Version)
	}

	return repo.DB.QueryRow(ctx, query).Scan(u.Version)
}

func (repo *userRepository) getByKey(key string, value any, tx pgx.Tx) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	u := new(models.User)
	query := fmt.Sprintf(`
		SELECT
			id,
			email,
			phone_number,
			first_name,
			last_name,
			is_phone_number_verified,
			is_email_verified,
			reg_stage,
			account_type,
			created_at,
			updated_at,
			deleted_at,
			version
		FROM
			users
		WHERE
			%s = $1
			AND deleted_at IS NULL`,
		key,
	)

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, value)
	} else {
		row = repo.DB.QueryRow(ctx, query, value)
	}

	err := row.Scan(
		u.ID,
		u.Email,
		u.PhoneNumber,
		u.FirstName,
		u.LastName,
		u.IsPhoneNumberVerified,
		u.IsEmailVerified,
		u.RegStage,
		u.AccountType,
		u.CreatedAt,
		u.UpdatedAt,
		u.DeletedAt,
		u.Version,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (repo *userRepository) GetById(id string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("id", id, tx)
}

func (repo *userRepository) GetByEmail(email string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("email", email, tx)
}

func (repo *userRepository) GetByPhoneNumber(phone string, tx pgx.Tx) (*models.User, error) {
	return repo.getByKey("phone_number", phone, tx)
}

func (repo *userRepository) Delete(id string, tx pgx.Tx) (err error) {
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

func (repo *userRepository) SoftDelete(id string, tx pgx.Tx) error {
	u, err := repo.GetById(id, tx)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()
	u.DeletedAt = now
	u.UpdatedAt = now
	return repo.Update(u, tx)
}