package handlers

import (
	"net/http"
	"time"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/models"
	"github.com/divinitymn/div-backend/internal/utils"

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

type UserUpdateMeBody struct {
	Avatar    string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	FirstName string `bson:"firstname,omitempty" json:"firstname,omitempty"`
	LastName  string `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Username  string `bson:"username" json:"username" validate:"required"`
	Bio       string `bson:"bio,omitempty" json:"bio,omitempty"`
	Password  string `bson:"password" json:"password" validate:"required"`
}

func UserUpdateMe(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Get("userId").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var r UserUpdateMeBody
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	var foundUser models.User

	if err := db.GetCollection("users").FindOne(
		c.Request().Context(),
		bson.M{"_id": id},
	).Decode(&foundUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  http.StatusNotFound,
				Message: "User not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to login",
		})
	}

	match, _, err := utils.CheckHash(r.Password, foundUser.Password)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.Response{
			Status:  http.StatusForbidden,
			Message: "Wrong password",
		})
	}

	if match {
		data := bson.M{}

		data["avatar"] = r.Avatar
		data["firstname"] = r.FirstName
		data["lastname"] = r.LastName
		data["username"] = r.Username
		data["bio"] = r.Bio

		if _, err := db.GetCollection("users").UpdateOne(
			c.Request().Context(),
			bson.M{"_id": id},
			bson.M{"$set": data},
		); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return c.JSON(http.StatusBadRequest, models.Response{
					Status:  http.StatusBadRequest,
					Message: "Username already exists",
				})
			}

			return err
		}

		return c.JSON(http.StatusOK, models.Response{
			Status:  http.StatusOK,
			Message: "User updated successfully",
		})
	}

	return c.JSON(http.StatusForbidden, models.Response{
		Status:  http.StatusForbidden,
		Message: "Wrong password",
	})
}
