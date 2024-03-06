package test_repositories

import (
	"time"

	"github.com/Bupher-Co/bupher-api/internal/models"
	"github.com/Bupher-Co/bupher-api/internal/repositories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OtpRepository struct {
	repo *repositories.OtpRepository
}

func NewOtpRepository(db *pgxpool.Pool, timeout time.Duration) *OtpRepository {
	return &OtpRepository{repo: repositories.NewOtpRepository(db, timeout)}
}

func (r *OtpRepository) Create(otp *models.Otp, tx pgx.Tx) error {
	return r.repo.Create(otp, tx)
}

func (r *OtpRepository) Update(otp *models.Otp, tx pgx.Tx) error {
	return r.repo.Update(otp, tx)
}

func (r *OtpRepository) GetById(id string, tx pgx.Tx) (*models.Otp, error) {
	return r.repo.GetById(id, tx)
}

func (r *OtpRepository) GetOneByWhere(where string, args []any, tx pgx.Tx) (*models.Otp, error) {
	return r.repo.GetOneByWhere(where, args, tx)
}

func (r *OtpRepository) Delete(id string, tx pgx.Tx) error {
	return r.repo.Delete(id, tx)
}

func (r *OtpRepository) SoftDelete(id string, tx pgx.Tx) error {
	return r.repo.SoftDelete(id, tx)
}
