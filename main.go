package main

import (
	"log"
	"os"
	"pojok_baca_api/database"
	"pojok_baca_api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    // Connect to database
    database.ConnectDB()

    // Initialize Fiber app
    app := fiber.New()

    // Register routes
    routes.SetupRoutes(app)

    // Start server
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000" // Default port if not specified
    }
    log.Fatal(app.Listen(":" + port))
}