package handlers

import (
	"log"
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

	Avatar      string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	DisplayName string `bson:"displayname" json:"displayname" validate:"required"`
	Bio         string `bson:"bio,omitempty" json:"bio,omitempty"`
	Email       string `bson:"email" json:"email" validate:"required,email"`
	Username    string `bson:"username" json:"username" validate:"required"`

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

	var r UserResponse
	if err = db.GetCollection("users").FindOne(
		c.Request().Context(),
		bson.M{"_id": id},
	).Decode(&r); err != nil {
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
		Data:    r,
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
	Avatar      string `bson:"avatar,omitempty" json:"avatar,omitempty"`
	DisplayName string `bson:"displayname,omitempty" json:"displayname,omitempty"`
	Username    string `bson:"username" json:"username" validate:"required"`
	Bio         string `bson:"bio,omitempty" json:"bio,omitempty"`
	Password    string `bson:"password" json:"password" validate:"required"`
}

func UserUpdateMe(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Get("userId").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var body UserUpdateMeBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	var r models.User

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

		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to login",
		})
	}

	match, _, err := utils.CheckHash(body.Password, r.Password)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.Response{
			Status:  http.StatusForbidden,
			Message: "Wrong password",
		})
	}

	if match {
		data := bson.M{}

    if body.Avatar != "" {
      data["avatar"] = body.Avatar
    }
    if body.Bio != "" {
      data["bio"] = body.Bio
    }
    if body.DisplayName != "" {
      data["displayname"] = body.DisplayName
    }
    data["username"] = body.Username

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

type UserChangePasswordMeBody struct {
	Password    string `bson:"password,omitempty" json:"password,omitempty"`
	NewPassword string `bson:"newPassword" json:"newPassword"`
}

func UserChangePasswordMe(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Get("userId").(string))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var body UserChangePasswordMeBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	var r models.User
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

    return c.JSON(http.StatusInternalServerError, models.Response{
      Status:  http.StatusInternalServerError,
      Message: "Failed to login",
    })
  }

  if hashedPassword, err := utils.CreateHash(body.NewPassword, utils.DefaultParams); err != nil {
    return c.JSON(http.StatusInternalServerError, models.Response{
      Status:  http.StatusInternalServerError,
      Message: "Failed to hash password",
    })
  } else {
    body.NewPassword = hashedPassword
  }

  if r.Password == "" {
    if _, err := db.GetCollection("users").UpdateOne(
      c.Request().Context(),
      bson.M{"_id": id},
      bson.M{"$set": bson.M{
        "password": body.NewPassword,
      }},
    ); err != nil {
      return err
    }

    return c.JSON(http.StatusOK, models.Response{
      Status:  http.StatusOK,
      Message: "Password updated successfully",
    })
  }

  match, _, err := utils.CheckHash(body.Password, r.Password)
  if err != nil {
    return err
  }

  if match {
    if _, err := db.GetCollection("users").UpdateOne(
      c.Request().Context(),
      bson.M{"_id": id},
      bson.M{"$set": bson.M{
        "password": body.NewPassword,
      }},
    ); err != nil {
      return err
    }

		return c.JSON(http.StatusOK, models.Response{
			Status:  http.StatusOK,
			Message: "Password updated successfully",
		})
  }

  return c.JSON(http.StatusForbidden, models.Response{
    Status:  http.StatusForbidden,
    Message: "Current password does not match",
  })

}
