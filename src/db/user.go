package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func (store *Store) CreateUser(ctx context.Context, data *model.UserCreate) (primitive.ObjectID, error) {

// 	data.CreatedAt = time.Now()
// 	data.UpdatedAt = time.Now()

// 	result, err := store.db.Collection(model.Collection).InsertOne(ctx, data)

// 	if err != nil {
// 		return primitive.NewObjectID(), common.ErrDB(err)
// 	}

// 	userId := result.InsertedID.(primitive.ObjectID)

// 	return userId, nil
// }

// func (noSql *Store) GetUser(ctx context.Context, condition map[string]string) (*model.User, error) {
// 	var user model.User
// 	filter := bson.M{}

// 	_, ok := condition["_id"]
// 	if ok {
// 		id, err := primitive.ObjectIDFromHex(condition["id"])
// 		if err != nil {
// 			return nil, common.ErrInvalidRequest(err)
// 		}
// 		filter["_id"] = id
// 	}

// 	email, ok := condition["email"]
// 	if ok {
// 		filter["email"] = email
// 	}

// 	err := noSql.db.Collection(model.Collection).FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return nil, common.RecordNotFound
// 		}
// 		log.Println(err)
// 		return nil, common.ErrDB(err)
// 	}

// 	fmt.Println(user)
// 	fmt.Println(&user)

// 	return &user, nil

// }

const (
	UserCollection = "users"
)

type CreateUserParams struct {
	Username    string                `json:"username" bson:"username"`
	Email       string                `json:"email" bson:"email"`
	Password    string                `json:"password" bson:"password"`
	Avatar      string                `json:"avatar" bson:"avatar"`
	Description string                `json:"description" bson:"description"`
	Posts       *[]primitive.ObjectID `json:"-" bson:"posts,omitempty"`
	CreatedAt   time.Time             `json:"-" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"-" bson:"updatedAt"`
}

func (q *Queries) CreateUser(ctx context.Context, data *CreateUserParams) (*primitive.ObjectID, error) {

	id, err := q.db.Collection(UserCollection).InsertOne(ctx, data)
	id2 := id.InsertedID.(primitive.ObjectID)
	return &id2, err
}

func (q *Queries) GetUserById(ctx context.Context, id *primitive.ObjectID) (*User, error) {
	var user User
	err := q.db.Collection(UserCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return &user, err
}
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := q.db.Collection(UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}
