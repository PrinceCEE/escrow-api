package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	repo *repositories.UserRepository
	mock.Mock
}

func NewUserRepository(db *pgxpool.Pool, timeout time.Duration) *UserRepository {
	return &UserRepository{repo: repositories.NewUserRepository(db, timeout)}
}

func (r *UserRepository) Create(u *models.User, tx pgx.Tx) error {
	return r.repo.Create(u, tx)
}

func (r *UserRepository) Update(u *models.User, tx pgx.Tx) error {
	return r.repo.Update(u, tx)
}

func (r *UserRepository) GetById(id string, tx pgx.Tx) (*models.User, error) {
	return r.repo.GetById(id, tx)
}

func (r *UserRepository) GetByBusinessId(id string, tx pgx.Tx) (*models.User, error) {
	return r.repo.GetByBusinessId(id, tx)
}

func (r *UserRepository) GetByEmail(email string, tx pgx.Tx) (*models.User, error) {
	return r.repo.GetByEmail(email, tx)
}

func (r *UserRepository) GetByPhoneNumber(phone string, tx pgx.Tx) (*models.User, error) {
	return r.repo.GetByPhoneNumber(phone, tx)
}

func (r *UserRepository) Delete(id string, tx pgx.Tx) error {
	return r.repo.Delete(id, tx)
}

func (r *UserRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
