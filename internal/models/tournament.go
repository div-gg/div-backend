package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tournament struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`

	Title string `bson:"title" json:"title" validate:"required"`
	Body  string `bson:"body" json:"body" validate:"required"` // As markdown text
	Game  string `bson:"game" json:"game" validate:"required,oneof=valorant cs2 lol"`

	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	CreatedUser primitive.ObjectID `bson:"created_user" json:"created_user"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedUser primitive.ObjectID `bson:"updated_user" json:"updated_user"`

	// TODO: Registration and deadlines
}
