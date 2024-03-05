package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	FirstName string `bson:"firstname" json:"firstname" validate:"required"`
	LastName  string `bson:"lastname" json:"lastname" validate:"required"`
	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	Username  string `bson:"username" json:"username" validate:"required"`
	Password  string `bson:"password" json:"password" validate:"required"`

  DiscordID string `bson:"discord_id" json:"discord_id"`
  DiscordName string `bson:"discord_name" json:"discord_name"`
  RiotID string `bson:"riot_id" json:"riot_id"`
  RiotName string `bson:"riot_name" json:"riot_name"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
