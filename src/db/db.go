package db

import "go.mongodb.org/mongo-driver/mongo"

func New(db *mongo.Database) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db *mongo.Database
}
