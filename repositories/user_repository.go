package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/Bupher-Co/bupher-api/models"
	"github.com/Bupher-Co/bupher-api/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(u *models.User) error
	Update(u *models.User) error
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPhoneNumber(phone string) (*models.User, error)
	Delete(id string) error
	SoftDelete(id string) error
}

type userRepository struct {
	DB      *pgxpool.Pool
	Timeout time.Duration
}

func NewUserRepository(db *pgxpool.Pool, timeout time.Duration) *userRepository {
	return &userRepository{DB: db, Timeout: timeout}
}

func (repo *userRepository) Create(u *models.User) error {
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

	return repo.DB.QueryRow(ctx, query, args...).Scan(u.ID, u.Version)
}

func (repo *userRepository) Update(u *models.User) error {
	u.UpdatedAt = time.Now().UTC()

	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query, err := pkg.GetUpdateQueryFromStruct(u, "users")
	if err != nil {
		return err
	}

	return repo.DB.QueryRow(ctx, query).Scan(u.Version)
}

func (repo *userRepository) getByKey(key string, value any) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	var u *models.User

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

	err := repo.DB.QueryRow(ctx, query, value).Scan(
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

func (repo *userRepository) GetById(id string) (*models.User, error) {
	return repo.getByKey("id", id)
}

func (repo *userRepository) GetByEmail(email string) (*models.User, error) {
	return repo.getByKey("email", email)
}

func (repo *userRepository) GetByPhoneNumber(phone string) (*models.User, error) {
	return repo.getByKey("phone_number", phone)
}

func (repo *userRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), repo.Timeout)
	defer cancel()

	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := repo.DB.Exec(ctx, query, id)
	return err
}

func (repo *userRepository) SoftDelete(id string) error {
	u, err := repo.GetById(id)
	if err != nil {
		return nil
	}

	u.DeletedAt = time.Now().UTC()
	return repo.Update(u)
}
