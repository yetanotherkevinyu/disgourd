package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Environment string

var Environments = struct {
	LOCAL string
}{
	LOCAL: "local",
}

type Config struct {
	DiscordClientSecret string
}

func LoadConfig() Config {
	env := getEnv("ENV", Environments.LOCAL)
	err := godotenv.Load(".env." + strings.ToLower(env))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
