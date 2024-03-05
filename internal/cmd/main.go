package cmd

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/divinitymn/div-backend/internal/db"
	"github.com/divinitymn/div-backend/internal/config"
)

func Execute() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
  config.InitEnv()
	db.InitDB()
	InitAPI()
}
