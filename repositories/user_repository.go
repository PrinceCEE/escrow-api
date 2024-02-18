package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type userRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) Create()     {}
func (repo *userRepository) FindOne()    {}
func (repo *userRepository) Find()       {}
func (repo *userRepository) UpdateOne()  {}
func (repo *userRepository) UpdateMany() {}
func (repo *userRepository) DeleteOne()  {}
func (repo *userRepository) DeleteMany() {}
