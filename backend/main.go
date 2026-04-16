package main

import (
	"log"
	"os"

	"url-shortener/backend/db"
	"url-shortener/backend/handlers"
	"url-shortener/backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	database := db.Connect()
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	port := getEnv("PORT", "8080")

	handler := handlers.NewURLHandler(database, baseURL)
	router := routes.SetupRouter(handler)

	log.Printf("backend is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
