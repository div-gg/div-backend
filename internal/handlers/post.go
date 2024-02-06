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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PostGetAll(c echo.Context) error {
	page, limit := utils.GetPaginationValues(c)

	sort := bson.D{primitive.E{Key: "created_at", Value: -1}}
	opts := options.Find().SetSort(sort).SetSkip(page).SetLimit(limit)
	filter := bson.D{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.GetCollection("posts").Find(ctx, filter, opts)
  if err != nil {
    return err
  }

	var results []models.Post
	if err = cursor.All(ctx, &results); err != nil {
		return err
	}

	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		log.Println(string(res))
	}

	return c.JSON(http.StatusOK, results)
}
