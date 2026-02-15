package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//struct untuk config database
type Config struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	Port       string
}

//variable dari struct
var AppConfig *Config


func LoadEnv() {
	// cek env kalo ga ada kasih error
	if err := godotenv.Load(); err != nil {
		log.Println("Error Not Found file .env !⚠️")
	}

	//instalasi untuk config
	AppConfig = &Config{
		Port:       os.Getenv("PORT"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBUser:     os.Getenv("DB_USER"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}