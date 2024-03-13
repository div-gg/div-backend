package handlers

import (
	"net/http"
	"time"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type ProfileUpdateBody struct {
	Avatar      string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	Bio         string `bson:"bio,omitempty" json:"bio,omitempty"`
	DisplayName string `bson:"displayname,omitempty" json:"displayname,omitempty"`
	Email       string `bson:"email,omitempty" json:"email,omitempty"`
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
	if user.DisplayName != "" {
		data["displayname"] = user.DisplayName
	}
	if user.Email != "" {
		data["email"] = user.Email
	}

	data["updated_at"] = time.Now()

	if _, err := db.GetCollection("users").UpdateByID(
		c.Request().Context(),
		id,
		bson.M{"$set": data},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Profile updated successfully",
	})
}
