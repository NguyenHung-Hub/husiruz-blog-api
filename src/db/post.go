package db

import (
	"context"
	"fmt"
	"husir_blog/src/common"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	PostCollection = "posts"
)

type CreatePostParams struct {
	Title       string                `json:"title" bson:"title"`
	Description string                `json:"description" bson:"description"`
	Photo       string                `json:"photo" bson:"photo"`
	Author      primitive.ObjectID    `json:"author" bson:"author"`
	Categories  []*primitive.ObjectID `json:"categories" bson:"categories"`
	Slug        string                `json:"slug" bson:"slug"`
	Status      string                `json:"status" bson:"status"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt" bson:"updatedAt"`
}
type PostResponseFull struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Photo       string             `json:"photo" bson:"photo"`
	Author      string             `json:"author" bson:"author"`
	Categories  []*Category        `json:"categories" bson:"categories"`
	Slug        string             `json:"slug" bson:"slug"`
	Status      string             `json:"status" bson:"status"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
type PostRecommendResponse struct {
	Title string `json:"title" bson:"title"`
	Slug  string `json:"slug" bson:"slug"`
}

func (q *Queries) CreatePost(ctx context.Context, data *CreatePostParams) (*primitive.ObjectID, error) {

	id, err := q.db.Collection(PostCollection).InsertOne(ctx, data)
	id2 := id.InsertedID.(primitive.ObjectID)
	return &id2, err
}

func (q *Queries) GetPostById(ctx context.Context, id *primitive.ObjectID) (*Post, error) {
	var res Post
	err := q.db.Collection(PostCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	return &res, err
}

func (q *Queries) GetPostBySlug(ctx context.Context, slug string) (*PostResponseFull, error) {

	p := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "slug", Value: slug}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: CategoryCollection},
			{Key: "localField", Value: "categories"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "categories"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: UserCollection},
			{Key: "localField", Value: "author"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "author"},
		}}},
		bson.D{{Key: "$unwind", Value: "$author"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "author", Value: "$author.username"}}}},
	}

	cursor, err := q.db.Collection(PostCollection).Aggregate(ctx, p)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	var post *PostResponseFull
	for cursor.Next(context.TODO()) {
		err = cursor.Decode(&post)
	}

	return post, err
}
func (q *Queries) GetPostByCategory(ctx context.Context, id string) ([]*PostResponseFull, error) {

	categoryId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	arr := []primitive.ObjectID{categoryId}
	p := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{{Key: CategoryCollection, Value: bson.D{{Key: "$in", Value: arr}}}}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: CategoryCollection},
			{Key: "localField", Value: "categories"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "categories"},
		}}},
	}

	cursor, err := q.db.Collection(PostCollection).Aggregate(ctx, p)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	var post []*PostResponseFull
	for cursor.Next(context.TODO()) {
		var p *PostResponseFull
		err = cursor.Decode(&p)
		post = append(post, p)
	}

	return post, err
}

func (q *Queries) ListPost(ctx context.Context, filter *PostFilter, paging *Paging) ([]*PostResponseFull, error) {

	skip := (paging.Page - 1) * paging.Limit

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: filter.Status}}}}
	if filter.CategoryId != "" {
		cateId, err := primitive.ObjectIDFromHex(filter.CategoryId)
		if err != nil {
			return nil, err
		}
		arr := []primitive.ObjectID{cateId}
		matchStage = bson.D{{Key: "$match", Value: bson.D{
			{Key: "status", Value: filter.Status},
			{Key: CategoryCollection, Value: bson.D{{Key: "$in", Value: arr}}},
		}}}
	}

	p := mongo.Pipeline{
		matchStage,
		bson.D{{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: -1}}}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: paging.Limit}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: CategoryCollection},
			{Key: "localField", Value: "categories"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "categories"},
		}}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: UserCollection},
			{Key: "localField", Value: "author"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "author"},
		}}},
		bson.D{{Key: "$unwind", Value: "$author"}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "author", Value: "$author.username"}}}},
	}

	fmt.Println(p)

	cursor, err := q.db.Collection(PostCollection).Aggregate(ctx, p)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	var list []*PostResponseFull
	for cursor.Next(context.TODO()) {
		var post PostResponseFull
		err = cursor.Decode(&post)
		list = append(list, &post)
	}

	return list, err
}

func (q *Queries) ListPostRandom(ctx context.Context, filter *PostFilter, nPost int) ([]*PostRecommendResponse, error) {

	p := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: filter.Status}}}},
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: nPost}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "title", Value: 1},
			{Key: "slug", Value: 1},
		}}},
	}

	cursor, err := q.db.Collection(PostCollection).Aggregate(ctx, p)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	var list []*PostRecommendResponse
	for cursor.Next(context.TODO()) {

		var post PostRecommendResponse
		err = cursor.Decode(&post)
		list = append(list, &post)
	}

	return list, err
}
