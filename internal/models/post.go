package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Body      string `bson:"body" json:"body" validate:"required"`
	OpenSlots int    `bson:"open_slots" json:"open_slots" validate:"required"`
	Game      string `bson:"game" json:"game" validate:"required, oneof=valorant csgo lol"`

	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	CreatedUser primitive.ObjectID `bson:"created_user" json:"created_user"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedUser primitive.ObjectID `bson:"updated_user" json:"updated_user"`
	ExpireAt    time.Time          `bson:"expire_at" json:"expire_at"`
}
