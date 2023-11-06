package common

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseIDModel struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
}

type BaseTimestamp struct {
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt primitive.Timestamp `json:"updatedAt" bson:"updatedAt"`
}

type BaseModel struct {
	BaseIDModel
	BaseTimestamp
}
