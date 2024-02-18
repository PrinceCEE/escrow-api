package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type authRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *authRepository {
	return &authRepository{DB: db}
}

func (repo *authRepository) Create()     {}
func (repo *authRepository) FindOne()    {}
func (repo *authRepository) Find()       {}
func (repo *authRepository) UpdateOne()  {}
func (repo *authRepository) UpdateMany() {}
func (repo *authRepository) DeleteOne()  {}
func (repo *authRepository) DeleteMany() {}
