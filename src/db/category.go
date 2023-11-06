package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CategoryCollection = "categories"
)

type CreateCategoryParams struct {
	Name      string    `json:"name" bson:"name"`
	Slug      string    `json:"slug" bson:"slug"`
	CreatedAt time.Time `json:"-" bson:"createdAt"`
	UpdatedAt time.Time `json:"-" bson:"updatedAt"`
}

type CategoryResponse struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Slug string             `json:"slug" bson:"slug"`
}

func (q *Queries) CreateCategory(ctx context.Context, data *CreateCategoryParams) (*primitive.ObjectID, error) {

	id, err := q.db.Collection(CategoryCollection).InsertOne(ctx, data)
	id2 := id.InsertedID.(primitive.ObjectID)
	return &id2, err
}

func (q *Queries) GetCategoryById(ctx context.Context, id *primitive.ObjectID) (*Category, error) {
	var res Category
	err := q.db.Collection(CategoryCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return &res, err
}
func (q *Queries) GetCategoryBySlug(ctx context.Context, slug string) (*Category, error) {
	var res Category
	err := q.db.Collection(CategoryCollection).FindOne(ctx, bson.M{"slug": slug}).Decode(&res)
	return &res, err
}

func (q *Queries) GetCategoriesName(ctx context.Context) ([]*CategoryResponse, error) {
	var res []*CategoryResponse
	cursor, err := q.db.Collection(CategoryCollection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var c CategoryResponse
		err = cursor.Decode(&c)
		if err != nil {
			log.Println(err)
		}
		res = append(res, &c)
	}

	return res, nil
}

func (q *Queries) DeleteCategoryBySlug(ctx context.Context, slug string) error {
	return q.db.Collection(CategoryCollection).FindOneAndDelete(ctx, bson.M{"slug": slug}).Err()
}
