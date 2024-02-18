package repositories

type Repository interface {
	Create()
	FindOne()
	Find()
	UpdateOne()
	UpdateMany()
	DeleteOne()
	DeleteMany()
}
