package handlers

import (
	"encoding/json"
	"log"
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
		log.Println(err)
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
				DiscordID: userInfo["id"].(string),
        DiscordName: userInfo["global_name"].(string),
				Email:     userInfo["email"].(string),
				Avatar:    "https://cdn.discordapp.com/avatars/" + userInfo["id"].(string) + "/" + userInfo["avatar"].(string) + ".png",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			cUser, err := db.GetCollection("users").InsertOne(c.Request().Context(), user)
			if err != nil {
				log.Println("Error creating user", err)
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
        "discord_id": userInfo["id"].(string),
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
