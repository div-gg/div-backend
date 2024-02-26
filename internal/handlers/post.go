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

func PostGetAll(c echo.Context) error {
	page, limit := utils.GetPaginationValues(c)

	sort := bson.D{primitive.E{Key: "updated_at", Value: -1}}
	opts := options.Find().SetSort(sort).SetSkip(page).SetLimit(limit)
	filter := bson.D{}

	cursor, err := db.GetCollection("posts").Find(
		context.TODO(),
		filter,
		opts,
	)
	if err != nil {
		return err
	}

	var results []models.Post
	if err = cursor.All(context.TODO(), &results); err != nil {
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
		Status:  http.StatusOK,
		Message: "Posts retrieved successfully",
		Data:    data,
	})
}

func PostGetByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var result models.Post

	err = db.GetCollection("posts").FindOne(
		context.TODO(),
		bson.M{"_id": id},
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  http.StatusNotFound,
				Message: "Post not found",
			})
		}

		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Post retrieved successfully",
		Data:    result,
	})
}

type PostCreateBody struct {
	Body      string `bson:"body" json:"body" validate:"required"`
	OpenSlots int    `bson:"open_slots" json:"open_slots" validate:"required"`
	Game      string `bson:"game" json:"game" validate:"required, oneof=valorant csgo lol"`

	ExpireAfter int `bson:"expire_after" json:"expire_after" validate:"required"`
}

func PostCreate(c echo.Context) error {
	post := new(PostCreateBody)
	if err := c.Bind(post); err != nil {
		return err
	}

	data := bson.M{}

	if post.Body == "" {
		data["body"] = post.Body
	}
	if post.OpenSlots == 0 {
		data["open_slots"] = post.OpenSlots
	}
	if post.Game == "" {
		data["game"] = post.Game
	}

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	data["created_user"] = c.Get("userId")
	data["updated_user"] = c.Get("userId")
	if post.ExpireAfter == 0 {
		data["expire_at"] = time.Now().Add(time.Duration(post.ExpireAfter) * time.Hour)
	} else {
		data["expire_at"] = time.Now().Add(time.Duration(6) * time.Hour)
	}

	_, err := db.GetCollection("posts").InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, models.Response{
		Status:  http.StatusCreated,
		Message: "Post created successfully",
	})
}

type PostUpdateBody struct {
	Body      string `bson:"body,omitempty" json:"body,omitempty"`
	OpenSlots int    `bson:"open_slots,omitempty" json:"open_slots,omitempty"`
	Game      string `bson:"game,omitempty" json:"game,omitempty" validate:"oneof=valorant csgo lol"`
}

func PostUpdateByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	post := new(PostUpdateBody)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	data := bson.M{}
	if post.Body != "" {
		data["body"] = post.Body
	}
	if post.OpenSlots != 0 {
		data["open_slots"] = post.OpenSlots
	}
	if post.Game != "" {
		data["game"] = post.Game
	}

	data["updated_at"] = time.Now()
	data["updated_user"] = c.Get("userId")

	_, err = db.GetCollection("posts").UpdateByID(
		context.TODO(),
		id,
		bson.M{"$set": data},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Post updated successfully",
	})
}

func PostDeleteByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	_, err = db.GetCollection("posts").DeleteOne(
		context.TODO(),
		bson.M{"_id": id},
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Post deleted successfully",
	})
}
