package models

import (
  "time"

  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
  ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

  Description string `bson:"description" json:"description" validate:"required"`
  OpenSlots int `bson:"open_slots" json:"open_slots" validate:"required"`
  Game string `bson:"game" json:"game" validate:"required, oneof=valorant csgo lol"`

  CreatedAt time.Time `bson:"created_at" json:"created_at"`
  UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
  ExpireAt time.Time `bson:"expire_at" json:"expire_at"`
}
