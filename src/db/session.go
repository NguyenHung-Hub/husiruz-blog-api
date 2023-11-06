package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	SessionCollection = "sessions"
)

type CreateSessionParams struct {
	SessionID    string    `json:"session_id" bson:"session_id"`
	Username     string    `json:"username" bson:"username"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	UserAgent    string    `json:"user_agent" bson:"user_agent"`
	ClientIp     string    `json:"client_ip" bson:"client_ip"`
	IsBlocked    bool      `json:"is_blocked" bson:"is_blocked"`
	ExpireAt     time.Time `json:"expire_at" bson:"expire_at"`
}

func (q *Queries) CreateSession(ctx context.Context, data *CreateSessionParams) (Session, error) {
	var session Session
	id, err := q.db.Collection(SessionCollection).InsertOne(ctx, data)
	if err != nil {
		return session, err
	}
	id2 := id.InsertedID.(primitive.ObjectID)
	err = q.db.Collection(SessionCollection).FindOne(ctx, bson.M{"_id": id2}).Decode(&session)
	return session, err
}

func (q *Queries) GetSessionById(ctx context.Context, id string) (*Session, error) {
	var session Session

	err := q.db.Collection(SessionCollection).FindOne(ctx, bson.M{"session_id": id}).Decode(&session)
	return &session, err
}
