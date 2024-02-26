package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tournament struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`

  Name string `json:"name" bson:"name" validate:"required"`
  Body string `json:"body" bson:"body" validate:"required"` // As markdown text
  Game string `json:"game" bson:"game" validate:"required, oneof=valorant csgo"`

  CreatedAt time.Time `json:"created_at" bson:"created_at"`
	CreatedUser primitive.ObjectID `bson:"created_user" json:"created_user"`
  UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
  UpdatedUser primitive.ObjectID `bson:"updated_user" json:"updated_user"`

  // TODO: Registration and deadlines
}
