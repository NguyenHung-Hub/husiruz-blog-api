package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID    `json:"_id" bson:"_id"`
	Username    string                `json:"username" bson:"username"`
	Email       string                `json:"email" bson:"email"`
	Password    string                `json:"password" bson:"password"`
	Avatar      string                `json:"avatar" bson:"avatar,omitempty"`
	Description string                `json:"description" bson:"description,omitempty"`
	Posts       *[]primitive.ObjectID `json:"posts,omitempty" bson:"posts,omitempty"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt" bson:"updatedAt"`
}

type Category struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Slug      string             `json:"slug" bson:"slug"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type PostStatus string

const (
	PostStatusVisibility PostStatus = "visibility"
	PostStatusHidden     PostStatus = "hidden"
	PostStatusDeleted    PostStatus = "deleted"
)

var AllPostStatus = map[string]string{
	"visibility": "visibility",
	"hidden":     "hidden",
	"deleted":    "deleted",
}

type Post struct {
	ID          primitive.ObjectID    `json:"_id" bson:"_id"`
	Title       string                `json:"title" bson:"title"`
	Description string                `json:"description" bson:"description"`
	Photo       string                `json:"photo" bson:"photo"`
	Author      primitive.ObjectID    `json:"author" bson:"author"`
	Categories  []*primitive.ObjectID `json:"categories" bson:"categories"`
	Slug        string                `json:"slug" bson:"slug"`
	Status      *PostStatus           `json:"status" bson:"status"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt" bson:"updatedAt"`
}

type Session struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	SessionID    string             `json:"session_id" bson:"session_id"`
	Username     string             `json:"username" bson:"username"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	UserAgent    string             `json:"user_agent" bson:"user_agent"`
	ClientIp     string             `json:"client_ip" bson:"client_ip"`
	IsBlocked    bool               `json:"is_blocked" bson:"is_blocked"`
	ExpireAt     time.Time          `json:"expire_at" bson:"expire_at"`
}
