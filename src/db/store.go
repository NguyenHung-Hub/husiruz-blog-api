package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Querier
}

type NoSQLStore struct {
	*Queries
}

func NewStore(db *mongo.Database) Store {
	return &NoSQLStore{Queries: New(db)}
}
