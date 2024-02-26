package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/middlewares"
	r "github.com/divinitymn/div-backend/internal/routes"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func InitAPI() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(middlewares.ErrorHandler)

	e.Validator = &CustomValidator{Validator: validator.New()}
	r.InitRoutes(e)

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
      e.Logger.Fatal(err)
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Shutdown server
	if err := e.Shutdown(context.TODO()); err != nil {
		e.Logger.Fatal(err)
	}

	// Disconnect from MongoDB
	if err := db.Client.Disconnect(context.Background()); err != nil {
		log.Fatal("⇨ Error disconnecting to MongoDB: ", err)
		log.Fatal(err)
	}

	return e
}
