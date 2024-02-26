package cmd

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/divinitymn/div-backend/internal/db"
)

func Execute() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	InitAPI()
}
