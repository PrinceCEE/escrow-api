package test_repositories

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/princecee/escrow-api/internal/models"
	"github.com/princecee/escrow-api/internal/repositories"
	"github.com/stretchr/testify/mock"
)

type TokenRepository struct {
	repo *repositories.TokenRepository
	mock.Mock
}

func NewTokenRepository(db *pgxpool.Pool, timeout time.Duration) *TokenRepository {
	return &TokenRepository{repo: repositories.NewTokenRepository(db, timeout)}
}

func (r *TokenRepository) Create(t *models.Token, tx pgx.Tx) error {
	return r.repo.Create(t, tx)
}

func (r *TokenRepository) Update(t *models.Token, tx pgx.Tx) error {
	return r.repo.Update(t, tx)
}

func (r *TokenRepository) GetById(id string, tx pgx.Tx) (*models.Token, error) {
	return r.repo.GetById(id, tx)
}

func (r *TokenRepository) Delete(id string, tx pgx.Tx) error {
	return r.repo.Delete(id, tx)
}

func (r *TokenRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
