package config

import (
	"os"
)

type EnvConfig struct {
	DbURL               string
	DbName              string
	JWTSecret           string
	DiscordClientID     string
	DiscordClientSecret string
	DiscordRedirectURI  string
}

var Env = &EnvConfig{}

func InitEnv() {
	Env = &EnvConfig{
		DbURL:               os.Getenv("DB_URL"),
		DbName:              os.Getenv("DB_NAME"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		DiscordClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		DiscordClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		DiscordRedirectURI:  os.Getenv("DISCORD_REDIRECT_URI"),
	}
}
