package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var configuration *Config

type DBConfig struct {
	Host string
	Port int
	Name string
	User string
	Password string
	EnableSSLMode bool

}

type Config struct {
	Version     string
	ServiceName string
	Port        int
	JwtScretKey string
	DB 			*DBConfig
}

func LoadConfig() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	version := os.Getenv("VERSION")
	if version == ""{
		fmt.Println("VERSION is required")
		os.Exit(1)
	}
	
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == ""{
		fmt.Println("SERVICE_NAME is required")
		os.Exit(1)
	}
	
	port := os.Getenv("PORT")
	if port == ""{
		fmt.Println("PORT is required")
		os.Exit(1)
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("PORT must be in Integer")
		os.Exit(1)
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == ""{
		fmt.Println("JWT_SECRET_KEY is required")
		os.Exit(1)
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		fmt.Println("DB_HOST is required")
		os.Exit(1)
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == ""{
		fmt.Println("DB_PORT is required")
		os.Exit(1)
	}

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		fmt.Println("PORT must be in Integer")
		os.Exit(1)
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		fmt.Println("DB_NAME is required")
		os.Exit(1)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		fmt.Println("DB_USER is required")
		os.Exit(1)
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		fmt.Println("DB_PASSWORD is required")
		os.Exit(1)
	}

	dbEnableMode := os.Getenv("DB_ENABLE_SSL_MODE")
	if dbEnableMode == "" {
		fmt.Println("DB_ENABLE_SSL_MODE is required")
		os.Exit(1)
	}

	parseddbEnableMode, err := strconv.ParseBool(dbEnableMode)
	if err != nil {
		fmt.Println("ENABLE_SSL_MODE must be a boolean (true/false)")
		os.Exit(1)
	}

	configuration = &Config{
		Version:     version,
		ServiceName: serviceName,
		Port:        portInt,
		JwtScretKey: jwtSecretKey,
		DB: &DBConfig{
			Host: dbHost,
			Port: dbPortInt,
			Name: dbName,
			User: dbUser,
			Password: dbPassword,
			EnableSSLMode: parseddbEnableMode,
		},
	}
}

func GetConfig() *Config {
	if configuration == nil {
		LoadConfig()
	}
	return configuration
}