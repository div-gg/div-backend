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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PostGetAll(c echo.Context) error {
	page, limit := utils.GetPaginationValues(c)

	sort := bson.D{primitive.E{Key: "updated_at", Value: -1}}
	opts := options.Find().SetSort(sort).SetSkip(page).SetLimit(limit)
	filter := bson.D{}

	cursor, err := db.GetCollection("posts").Find(context.TODO(), filter, opts)
	if err != nil {
    if err == mongo.ErrNoDocuments {
      return c.JSON(http.StatusOK, []models.Post{})
    }
		return err
	}

	var results []models.Post
	if err = cursor.All(context.TODO(), &results); err != nil {
    if err == mongo.ErrNoDocuments {
      return c.JSON(http.StatusOK, results)
    }

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
    data = []models.Post{}
  }

	return c.JSON(http.StatusOK, models.Response{
    Status: http.StatusOK,
    Message: "Posts retrieved successfully",
    Data: data,
  })
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

	_, err := db.GetCollection("posts").InsertOne(context.TODO(), post)
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

type PostUpdateRequest struct {
	Body      string `json:"body" validate:"required"`
	OpenSlots int    `json:"open_slots" validate:"required"`
	Game      string `json:"game" validate:"required,oneof=valorant csgo lol"`

  UpdatedAt time.Time `json:"updated_at"`
}

func PostUpdate(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return err
	}

	post := new(PostUpdateRequest)
	if err := c.Bind(post); err != nil {
		return err
	}
	post.UpdatedAt = time.Now()

	update := bson.M{ "$set": post }

	_, err = db.GetCollection("posts").UpdateByID(context.TODO(), id, update)
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

	_, err = db.GetCollection("posts").DeleteOne(context.TODO(), filter)
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
