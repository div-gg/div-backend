package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	Title string `bson:"title" json:"title"`
	Body  string `bson:"body" json:"body"`

	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	CreatedUser primitive.ObjectID `bson:"created_user" json:"created_user"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedUser primitive.ObjectID `bson:"updated_user" json:"updated_user"`
}
