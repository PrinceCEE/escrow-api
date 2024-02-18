package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type tokenRepository struct {
	DB *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *tokenRepository {
	return &tokenRepository{DB: db}
}

func (repo *tokenRepository) Create()     {}
func (repo *tokenRepository) FindOne()    {}
func (repo *tokenRepository) Find()       {}
func (repo *tokenRepository) UpdateOne()  {}
func (repo *tokenRepository) UpdateMany() {}
func (repo *tokenRepository) DeleteOne()  {}
func (repo *tokenRepository) DeleteMany() {}
