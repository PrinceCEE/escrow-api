package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type businessRepository struct {
	DB *pgxpool.Pool
}

func NewBusinessRepository(db *pgxpool.Pool) *businessRepository {
	return &businessRepository{DB: db}
}

func (repo *businessRepository) Create()     {}
func (repo *businessRepository) FindOne()    {}
func (repo *businessRepository) Find()       {}
func (repo *businessRepository) UpdateOne()  {}
func (repo *businessRepository) UpdateMany() {}
func (repo *businessRepository) DeleteOne()  {}
func (repo *businessRepository) DeleteMany() {}
