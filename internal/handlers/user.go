package handlers

import (
	"net/http"
	"time"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserResponse struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	FirstName string `bson:"firstname" json:"firstname" validate:"required"`
	LastName  string `bson:"lastname" json:"lastname" validate:"required"`
	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	Username  string `bson:"username" json:"username" validate:"required"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func UserGetByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var result UserResponse

	if err = db.GetCollection("users").FindOne(
		c.Request().Context(),
		bson.M{"_id": id},
	).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  http.StatusNotFound,
				Message: "Post not found",
			})
		}
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "User retrieved successfully",
		Data:    result,
	})
}

func UserGetMe(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Get("userId").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var r UserResponse

	if err := db.GetCollection("users").FindOne(
		c.Request().Context(),
    bson.M{"_id": id},
	).Decode(&r); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  http.StatusNotFound,
				Message: "User not found",
			})
		}

		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "User retrieved successfully",
		Data:    r,
	})
}

type UserUpdateBody struct {
	Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	FirstName string `bson:"firstname" json:"firstname" validate:"required"`
	LastName  string `bson:"lastname" json:"lastname" validate:"required"`
	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
}

func UserUpdateMe(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Get("userId").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	user.UpdatedAt = time.Now()

	if _, err = db.GetCollection("users").UpdateOne(
		c.Request().Context(),
    bson.M{"_id": id},
    bson.M{"$set": user},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "User updated successfully",
	})
}
