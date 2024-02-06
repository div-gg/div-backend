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

func PostGetOne(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result models.Post

	err = db.GetCollection("posts").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func PostCreate(c echo.Context) error {
	post := new(models.Post)
	if err := c.Bind(post); err != nil {
		return err
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
  post.ExpireAt = time.Now().Add(time.Second * 10)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.GetCollection("posts").InsertOne(ctx, post)
	if err != nil {
		return err
	}

	return c.JSON(
    http.StatusCreated,
		models.Response{
			Status:  http.StatusCreated,
			Message: "Post created successfully",
		},
	)
}

func PostUpdate(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	post := new(models.Post)
	if err := c.Bind(post); err != nil {
		return err
	}

	post.UpdatedAt = time.Now()

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{primitive.E{Key: "$set", Value: post}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

  _, err = db.GetCollection("posts").UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return c.JSON(
    http.StatusOK,
    models.Response{
      Status:  http.StatusOK,
      Message: "Post updated successfully",
    },
  )
}

func PostDelete(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = db.GetCollection("posts").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return c.JSON(
    http.StatusOK,
    models.Response{
      Status:  http.StatusOK,
      Message: "Post deleted successfully",
    },
  )
}
