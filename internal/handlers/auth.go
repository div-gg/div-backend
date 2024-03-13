package handlers

import (
	"net/http"
	"time"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/models"
	"github.com/divinitymn/div-backend/internal/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RegisterRequest struct {
	DisplayName string    `bson:"displayname" json:"displayname" validate:"required"`
	Email       string    `bson:"email" json:"email" validate:"required,email"`
	Username    string    `bson:"username" json:"username" validate:"required"`
	Password    string    `bson:"password" json:"password" validate:"required"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func RegisterHandler(c echo.Context) error {
	r := RegisterRequest{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&r); err != nil {
		return err
	}

	if hashedPassword, err := utils.CreateHash(r.Password, utils.DefaultParams); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to hash password",
		})
	} else {
		r.Password = hashedPassword
	}

	if _, err := db.GetCollection("users").InsertOne(
		c.Request().Context(),
		r,
	); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.JSON(http.StatusConflict, models.Response{
				Status:  http.StatusConflict,
				Message: "Username or email already exists",
			})
		}

		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to register",
		})
	}

	return c.JSON(http.StatusCreated, models.Response{
		Status:  http.StatusCreated,
		Message: "Register success",
	})
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func LoginHandler(c echo.Context) error {
	data := LoginRequest{}

	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&data); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var result models.User

	if err := db.GetCollection("users").FindOne(
		c.Request().Context(),
		bson.M{"username": data.Username},
	).Decode(&result); err != nil {
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

	match, _, err := utils.CheckHash(data.Password, result.Password)
	if err != nil {
		return c.JSON(http.StatusForbidden, models.Response{
			Status:  http.StatusForbidden,
			Message: "Wrong password",
		})
	}

	if match {
		token, err := utils.CreateToken(result.ID.Hex())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, models.Response{
			Status:  http.StatusOK,
			Message: "Login success",
			Data: map[string]string{
				"token": token,
			},
		})
	}

	return c.JSON(http.StatusForbidden, models.Response{
		Status:  http.StatusForbidden,
		Message: "Wrong password",
	})
}
