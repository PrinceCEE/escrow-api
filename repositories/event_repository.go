package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type eventRepository struct {
	DB *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *eventRepository {
	return &eventRepository{DB: db}
}

func (repo *eventRepository) Create()     {}
func (repo *eventRepository) FindOne()    {}
func (repo *eventRepository) Find()       {}
func (repo *eventRepository) UpdateOne()  {}
func (repo *eventRepository) UpdateMany() {}
func (repo *eventRepository) DeleteOne()  {}
func (repo *eventRepository) DeleteMany() {}
