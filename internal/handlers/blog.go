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

func BlogGetAll(c echo.Context) error {
	page, limit := utils.GetPaginationValues(c)

	opts := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(page).
		SetLimit(limit)

	cursor, err := db.GetCollection("blogs").Find(
		c.Request().Context(),
		bson.M{},
		opts,
	)
	if err != nil {
		return err
	}

	var results []models.Blog
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
		data = []models.Blog{}
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Blogs retrieved successfully",
		Data:    data,
	})
}

func BlogGetByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	var result models.Blog

	err = db.GetCollection("blogs").FindOne(
		c.Request().Context(),
		bson.M{"_id": id},
	).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, models.Response{
				Status:  http.StatusNotFound,
				Message: "Blog not found",
			})
		}

		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Blog retrieved successfully",
		Data:    result,
	})
}

type BlogCreateBody struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

func BlogCreate(c echo.Context) error {
	blog := new(BlogCreateBody)
	if err := c.Bind(&blog); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&blog); err != nil {
		return err
	}

	data := bson.M{}

	data["title"] = blog.Title
	data["body"] = blog.Body
	data["created_at"] = time.Now()
	data["updated_at"] = time.Now()
	data["created_user"] = c.Get("userId")
	data["updated_user"] = c.Get("userId")

	res, err := db.GetCollection("blogs").InsertOne(c.Request().Context(), blog)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, models.Response{
		Status:  http.StatusCreated,
		Message: "Blog created successfully",
		Data:    res.InsertedID,
	})
}

type BlogUpdateBody struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

func BlogUpdateByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	blog := new(BlogUpdateBody)
	if err := c.Bind(&blog); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	data := bson.M{}
	if blog.Title != "" {
		data["title"] = blog.Title
	}
	if blog.Body != "" {
		data["body"] = blog.Body
	}
	data["updated_at"] = time.Now()
	data["updated_user"] = c.Get("userId")

	if _, err = db.GetCollection("blogs").UpdateOne(
		c.Request().Context(),
		id,
		bson.M{"$set": data},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Blog updated successfully",
	})
}

func BlogDeleteByID(c echo.Context) error {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
		})
	}

	if _, err = db.GetCollection("blogs").DeleteOne(
		c.Request().Context(),
		bson.M{"_id": id},
	); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, models.Response{
		Status:  http.StatusOK,
		Message: "Blog deleted successfully",
	})
}
