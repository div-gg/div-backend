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
)

type UserObject struct {
	Avatar      string `bson:"avatar" json:"avatar"`
	DisplayName string `bson:"displayname" json:"displayname"`
	Username    string `bson:"username" json:"username"`
}

type PostResult struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Body  string `bson:"body" json:"body" validate:"required"`
	Slots int    `bson:"slots" json:"slots" validate:"required"`
	Game  string `bson:"game" json:"game" validate:"required,oneof=valorant cs2 lol"`

	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	CreatedUser []UserObject       `bson:"created_user" json:"created_user"`
	CreatedBy   primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by" json:"updated_by"`
	ExpireAt    time.Time          `bson:"expire_at" json:"expire_at"`
}

func PostGetAll(c echo.Context) error {
	page, limit := utils.GetPaginationValues(c)

	lookupStage := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "created_by"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "created_user"},
		}},
	}
	unwindStage := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$created_by"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}
	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "updated_at", Value: -1},
		}},
	}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	skipStage := bson.D{{Key: "$skip", Value: page}}

	cursor, err := db.GetCollection("posts").Aggregate(
		c.Request().Context(),
		mongo.Pipeline{
			lookupStage,
			unwindStage,
			sortStage,
			limitStage,
			skipStage,
		},
	)
	if err != nil {
		return err
	}

	var results []PostResult
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
		data = []PostResult{}
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
		c.Request().Context(),
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
	Body  string `bson:"body" json:"body" validate:"required"`
	Slots int    `bson:"slots" json:"slots" validate:"required"`
	Game  string `bson:"game" json:"game" validate:"required,oneof=valorant cs2 lol"`
}

func PostCreate(c echo.Context) error {
	post := PostCreateBody{}
	if err := c.Bind(&post); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if err := c.Validate(&post); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	data := bson.M{}

	if post.Body != "" {
		data["body"] = post.Body
	}
	if post.Slots != 0 {
		data["slots"] = post.Slots
	}
	if post.Game != "" {
		data["game"] = post.Game
	}

	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	data["created_by"] = c.Get("userId")
	data["updated_by"] = c.Get("userId")
	data["expire_at"] = time.Now().Add(time.Duration(6) * time.Hour)

	if _, err := db.GetCollection("posts").InsertOne(
		c.Request().Context(),
		data,
	); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, models.Response{
		Status:  http.StatusCreated,
		Message: "Post created successfully",
	})
}

type PostUpdateBody struct {
	Body  string `bson:"body,omitempty" json:"body,omitempty"`
	Slots int    `bson:"slots,omitempty" json:"slots,omitempty"`
	Game  string `bson:"game,omitempty" json:"game,omitempty" validate:"oneof=valorant cs2 lol"`
}

func PostUpdateByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	post := PostUpdateBody{}
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	if err := c.Validate(post); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	data := bson.M{}
	if post.Body != "" {
		data["body"] = post.Body
	}
	if post.Slots != 0 {
		data["slots"] = post.Slots
	}
	if post.Game != "" {
		data["game"] = post.Game
	}

	data["updated_at"] = time.Now()
	data["updated_by"] = c.Get("userId")

	_, err = db.GetCollection("posts").UpdateByID(
		c.Request().Context(),
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

	if _, err = db.GetCollection("posts").DeleteOne(
		c.Request().Context(),
		bson.M{"_id": id},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Post deleted successfully",
	})
}
