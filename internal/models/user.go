package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	DisplayName string `bson:"displayname,omitempty" json:"displayname,omitempty" validate:"required"`
	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	Email     string `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Username  string `bson:"username,omitempty" json:"username,omitempty" validate:"required"`
	Password  string `bson:"password,omitempty" json:"password,omitempty" validate:"required"`

	DiscordID   string `bson:"discord_id,omitempty" json:"discord_id,omitempty"`
	DiscordName string `bson:"discord_name,omitempty" json:"discord_name,omitempty"`
	RiotID      string `bson:"riot_id,omitempty" json:"riot_id,omitempty"`
	RiotName    string `bson:"riot_name,omitempty" json:"riot_name,omitempty"`

	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
