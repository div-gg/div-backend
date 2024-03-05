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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TournamentGetAll(c echo.Context) error {
	page, limit := utils.GetPaginationValues(c)
	opts := options.
    Find().
    SetSort(bson.M{"updated_at": -1}).
    SetSkip(page).
    SetLimit(limit)

	cursor, err := db.GetCollection("tournaments").Find(
    c.Request().Context(),
		bson.D{},
		opts,
	)
	if err != nil {
		return err
	}

	var results []models.Tournament
	if err = cursor.All(c.Request().Context(), &results); err != nil {
		return err
	}

	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		log.Println(string(res))
	}

	var data interface{}

	if results != nil {
		data = results
	} else {
		data = []models.Tournament{}
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Tournaments retrieved successfully",
		Data:    data,
	})
}

func TournamentGetByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var result models.Tournament

	if err = db.GetCollection("tournaments").FindOne(
    c.Request().Context(),
		bson.M{"_id": id},
	).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  http.StatusNotFound,
				Message: "Tournament not found",
			})
		}

    return c.JSON(http.StatusInternalServerError, models.Response{
      Status:  http.StatusInternalServerError,
      Message: "Failed to retrieve tournament",
    })
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Tournament retrieved successfully",
		Data:    result,
	})
}

type TournamentCreateBody struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
	Game  string `json:"game" validate:"required,oneof=valorant csgo"`
}

func TournamentCreate(c echo.Context) error {
	tournament := new(TournamentCreateBody)
	if err := c.Bind(tournament); err != nil {
		return err
	}
  if err := c.Validate(&tournament); err != nil {
    return err
  }

	data := bson.M{}

  data["title"] = tournament.Title
  data["body"] = tournament.Body
  data["game"] = tournament.Game
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	data["created_user"] = c.Get("userId")
	data["updated_user"] = c.Get("userId")

	if _, err := db.GetCollection("tournaments").InsertOne(c.Request().Context(), data); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, models.Response{
		Status:  http.StatusCreated,
		Message: "Tournament created successfully",
	})
}

type TournamentUpdateBody struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
	Game  string `json:"game,omitempty"`
}

func TournamentUpdateByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	tournament := new(TournamentUpdateBody)
	if err := c.Bind(tournament); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	data := bson.M{}
	if tournament.Title != "" {
		data["title"] = tournament.Title
	}
	if tournament.Body != "" {
		data["body"] = tournament.Body
	}
	if tournament.Game != "" {
		data["game"] = tournament.Game
	}

	data["updated_at"] = time.Now()
	data["updated_user"] = c.Get("userId")

	if _, err = db.GetCollection("tournaments").UpdateByID(
    c.Request().Context(),
		id,
		bson.M{"$set": data},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Tournament updated successfully",
	})
}

func TournamentDeleteByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	if _, err = db.GetCollection("tournaments").DeleteOne(
    c.Request().Context(),
		bson.M{"_id": id},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Tournament deleted successfully",
	})
}
