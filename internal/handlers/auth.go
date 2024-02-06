package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/divinitymn/aion-backend/internal/db"
	"github.com/divinitymn/aion-backend/internal/models"
	"github.com/divinitymn/aion-backend/internal/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterRequest struct {
	FirstName string    `bson:"firstname" json:"firstname" validate:"required"`
	LastName  string    `bson:"lastname" json:"lastname" validate:"required"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Username  string    `bson:"username" json:"username" validate:"required"`
	Password  string    `bson:"password" json:"password" validate:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
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

	r.Password, _ = utils.CreateHash(r.Password, utils.DefaultParams)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.GetCollection("users").InsertOne(ctx, r)
	if err != nil {
		log.Print(err)
		return err
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
	r := LoginRequest{}

	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&r); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result models.User
	filter := bson.D{primitive.E{Key: "username", Value: r.Username}}

	err := db.GetCollection("users").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Print(err)
		return err
	}

	match, _, err := utils.CheckHash(r.Password, result.Password)
	if err != nil {
		log.Print(err)
		return err
	}

	if match {
		token, err := utils.CreateToken(result.ID.String())
		if err != nil {
			log.Print(err)
			return err
		}

		return c.JSON(
			http.StatusOK,
			models.Response{
				Status:  http.StatusOK,
				Message: "Login success",
				Data: map[string]string{
					"token": token,
				},
			},
		)
	}

	return c.JSON(http.StatusUnauthorized, models.Response{
		Status:  http.StatusUnauthorized,
		Message: "Login failed",
	})
}
