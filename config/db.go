package config

import "database/sql"

type DbManager struct {
	Db *sql.DB
}

func newDbManager(env *Env) *DbManager {
	// connect to the DB
	return &DbManager{}
}

func (m *DbManager) Save() interface{} {
	return "Not implemented"
}

func (m *DbManager) Update() interface{} {
	return "Not Implemented"
}

func (m *DbManager) FindOne() interface{} {
	return "Not implemented"
}

func (m *DbManager) Find() interface{} {
	return "Not implemented"
}

func (m *DbManager) SoftDelete() interface{} {
	return "Not implemented"
}

func (m *DbManager) HardDelete() interface{} {
	return "Not implemented"
}
