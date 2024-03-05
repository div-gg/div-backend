package cmd

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/config"
)

func Execute() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

  config.InitEnv()
	db.InitDB()
	InitAPI()
}
