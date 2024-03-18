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

	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy primitive.ObjectID `bson:"updated_by" json:"updated_by"`

	// TODO: Registration and deadlines
}
