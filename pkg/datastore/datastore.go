package datastore

import (
	"github.com/jmoiron/sqlx"
)

type Storage interface {
}

type Datastore struct {
	DB *sqlx.DB
}

func NewDatastore(db *sqlx.DB) *Datastore {
	return &Datastore{
		DB: db,
	}
}
