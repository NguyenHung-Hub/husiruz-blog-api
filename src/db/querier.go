package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Querier interface {
	CreateUser(ctx context.Context, data *CreateUserParams) (*primitive.ObjectID, error)
	GetUserById(ctx context.Context, id *primitive.ObjectID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	CreateCategory(ctx context.Context, data *CreateCategoryParams) (*primitive.ObjectID, error)
	GetCategoryById(ctx context.Context, id *primitive.ObjectID) (*Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	GetCategoriesName(ctx context.Context) ([]*CategoryResponse, error)
	DeleteCategoryBySlug(ctx context.Context, slug string) error

	CreatePost(ctx context.Context, data *CreatePostParams) (*primitive.ObjectID, error)
	GetPostById(ctx context.Context, id *primitive.ObjectID) (*Post, error)
	GetPostBySlug(ctx context.Context, slug string) (*PostResponseFull, error)
	GetPostByCategory(ctx context.Context, id string) ([]*PostResponseFull, error)
	ListPost(ctx context.Context, filter *PostFilter, paging *Paging) ([]*PostResponseFull, error)
	ListPostRandom(ctx context.Context, filter *PostFilter, nPost int) ([]*PostRecommendResponse, error)

	CreateSession(ctx context.Context, data *CreateSessionParams) (Session, error)
	GetSessionById(ctx context.Context, id string) (*Session, error)
}

var _ Querier = (*Queries)(nil)
