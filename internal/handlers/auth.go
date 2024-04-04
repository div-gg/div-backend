package handlers

import (
  "encoding/json"
	"net/http"
	"time"

  "github.com/divinitymn/div-backend/internal/config"
	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/models"
	"github.com/divinitymn/div-backend/internal/utils"


  "golang.org/x/oauth2"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	if err := c.Validate(&r); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
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
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
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

func DiscordCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")

	discordOAuth := &oauth2.Config{
		ClientID:     config.Env.DiscordClientID,
		ClientSecret: config.Env.DiscordClientSecret,
		RedirectURL:  config.Env.DiscordRedirectURI,
		Scopes:       []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}

	oAuthToken, err := discordOAuth.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to exchange token",
		})
	}

	client := discordOAuth.Client(c.Request().Context(), oAuthToken)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get user",
		})
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to decode user",
		})
	}

	var r models.User
	if err := db.GetCollection("users").FindOne(
		c.Request().Context(),
		bson.M{
			"$or": bson.A{
				bson.M{"discord_id": userInfo["id"]},
				bson.M{"email": userInfo["email"]},
			},
		},
	).Decode(&r); err != nil {
		if err == mongo.ErrNoDocuments {
			user := models.User{
				DiscordID:   userInfo["id"].(string),
				DiscordName: userInfo["global_name"].(string),
				Email:       userInfo["email"].(string),
				Avatar:      "https://cdn.discordapp.com/avatars/" + userInfo["id"].(string) + "/" + userInfo["avatar"].(string) + ".png",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			cUser, err := db.GetCollection("users").InsertOne(c.Request().Context(), user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, models.Response{
					Status:  http.StatusInternalServerError,
					Message: "Failed to create user",
				})
			}

			token, err := utils.CreateToken(cUser.InsertedID.(primitive.ObjectID).Hex())
			if err != nil {
				return c.JSON(http.StatusInternalServerError, models.Response{
					Status:  http.StatusInternalServerError,
					Message: "Failed to create token",
				})
			}

			return c.JSON(http.StatusOK, models.Response{
				Status:  http.StatusOK,
				Message: "User created successfully",
				Data: map[string]string{
					"token": token,
				},
			})
		}

		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create token",
		})
	}

	if r.DiscordID == "" {
		_, err := db.GetCollection("users").UpdateOne(
			c.Request().Context(),
			bson.M{"_id": r.ID},
			bson.M{"$set": bson.M{
				"discord_id":   userInfo["id"].(string),
				"discord_name": userInfo["global_name"].(string),
			}},
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.Response{
				Status:  http.StatusInternalServerError,
				Message: "Failed to update user",
			})
		}
	}

	token, err := utils.CreateToken(r.ID.Hex())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create token",
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data: map[string]string{
			"token": token,
		},
	})
}
