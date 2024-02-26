package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

type ProfileUpdateBody struct {
	Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	FirstName string `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName  string `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Email     string `bson:"email,omitempty" json:"email,omitempty"`
}

func ProfileUpdate(c echo.Context) error {
	id := c.Get("userId").(string)

	var user ProfileUpdateBody
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	data := bson.M{}
	if user.Avatar != "" {
		data["avatar"] = user.Avatar
	}
	if user.Bio != "" {
		data["bio"] = user.Bio
	}
	if user.FirstName != "" {
		data["firstname"] = user.FirstName
	}
	if user.LastName != "" {
		data["lastname"] = user.LastName
	}
	if user.Email != "" {
		data["email"] = user.Email
	}

	data["updated_at"] = time.Now()

	_, err := db.GetCollection("users").UpdateByID(
		context.TODO(),
		id,
		bson.M{"$set": data},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Profile updated successfully",
	})
}
