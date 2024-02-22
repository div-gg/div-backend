package handlers

import (
  "context"
  "net/http"

  "github.com/divinitymn/aion-backend/internal/db"
  "github.com/divinitymn/aion-backend/internal/models"

  "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserGetByID(c echo.Context) error {
  id, err := primitive.ObjectIDFromHex(c.Param("id"))
  if err != nil {
    return c.JSON(http.StatusBadRequest, models.Response{
      Status: http.StatusBadRequest,
      Message: "Invalid ID",
    })
  }

  filter := bson.D{primitive.E{Key: "_id", Value: id}}
  var result models.User

  err = db.GetCollection("users").FindOne(context.TODO(), filter).Decode(&result)
  if err != nil {
    if err == mongo.ErrNoDocuments {
      return c.JSON(http.StatusNotFound, models.Response{
        Status: http.StatusNotFound,
        Message: "Post not found",
      })
    }
    return err
  }

  return c.JSON(http.StatusOK, models.Response{
    Status: http.StatusOK,
    Message: "User retrieved successfully",
    Data: result,
  })
}
