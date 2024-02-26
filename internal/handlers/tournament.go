package handlers

import (
	"context"
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

  sort := bson.D{primitive.E{Key: "updated_at", Value: -1}}
  opts := options.Find().SetSort(sort).SetSkip(page).SetLimit(limit)
  filter := bson.D{}

  cursor, err := db.GetCollection("tournaments").Find(
    context.TODO(),
    filter,
    opts,
  )
  if err != nil {
    return err
  }

  var results []models.Tournament
  if err = cursor.All(context.Background(), &results); err != nil {
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

  return c.JSON(200, models.Response{
    Status: http.StatusOK,
    Message: "Tournaments retrieved successfully",
    Data: data,
  })
}

func TournamentGetByID(c echo.Context) error {
  id, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil {
    return c.JSON(http.StatusBadRequest, models.Response{
      Status: http.StatusBadRequest,
      Message: "Invalid ID",
    })
  }

  var result models.Tournament

  err = db.GetCollection("tournaments").FindOne(
    context.TODO(),
    bson.M{"_id": id},
  ).Decode(&result)
  if err != nil {
    if err == mongo.ErrNoDocuments {
      return c.JSON(http.StatusNotFound, models.Response{
        Status: http.StatusNotFound,
        Message: "Tournament not found",
      })
    }

    return err
  }

  return c.JSON(http.StatusOK, models.Response{
    Status: http.StatusOK,
    Message: "Tournament retrieved successfully",
    Data: result,
  })
}

type TournamentCreateBody struct {
  Name string `json:"name" validate:"required"`
  Body string `json:"body" validate:"required"`
  Game string `json:"game" validate:"required,oneof=valorant csgo"`
}

func TournamentCreate(c echo.Context) error {
  tournament := new(TournamentCreateBody)
  if err := c.Bind(tournament); err != nil {
    return err
  }

  data := bson.M{}

  if tournament.Name != "" {
    data["name"] = tournament.Name
  }
  if tournament.Body != "" {
    data["body"] = tournament.Body
  }
  if tournament.Game != "" {
    data["game"] = tournament.Game
  }

  data["created_at"] = time.Now()
  data["updated_at"] = time.Now()
  data["created_user"] = c.Get("userId")
  data["updated_user"] = c.Get("userId")

  _, err := db.GetCollection("tournaments").InsertOne(context.TODO(), data)
  if err != nil {
    return err
  }

  return c.JSON(http.StatusCreated, models.Response{
    Status: http.StatusCreated,
    Message: "Tournament created successfully",
  })
}

type TournamentUpdateBody struct {
  Name string `json:"name,omitempty"`
  Body string `json:"body,omitempty"`
  Game string `json:"game,omitempty"`
}

func TournamentUpdateByID(c echo.Context) error {
  id, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil {
    return c.JSON(http.StatusBadRequest, models.Response{
      Status: http.StatusBadRequest,
      Message: "Invalid ID",
    })
  }

  tournament := new(TournamentUpdateBody)
  if err := c.Bind(tournament); err != nil {
    return c.JSON(http.StatusBadRequest, models.Response{
      Status: http.StatusBadRequest,
      Message: "Invalid request body",
    })
  }

  data := bson.M{}
  if tournament.Name != "" {
    data["name"] = tournament.Name
  }
  if tournament.Body != "" {
    data["body"] = tournament.Body
  }
  if tournament.Game != "" {
    data["game"] = tournament.Game
  }

  data["updated_at"] = time.Now()
  data["updated_user"] = c.Get("userId")

  _, err = db.GetCollection("tournaments").UpdateByID(
    context.TODO(),
    id,
    bson.M{"$set": data},
  )
  if err != nil {
    return err
  }

  return c.JSON(http.StatusOK, models.Response{
    Status: http.StatusOK,
    Message: "Tournament updated successfully",
  })
}

func TournamentDeleteByID(c echo.Context) error {
  id, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil {
    return c.JSON(http.StatusBadRequest, models.Response{
      Status: http.StatusBadRequest,
      Message: "Invalid ID",
    })
  }

  _, err = db.GetCollection("tournaments").DeleteOne(
    context.TODO(),
    bson.M{"_id": id},
  )
  if err != nil {
    return err
  }

  return c.JSON(http.StatusOK, models.Response{
    Status: http.StatusOK,
    Message: "Tournament deleted successfully",
  })
}
