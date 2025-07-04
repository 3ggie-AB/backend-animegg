package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/3ggie-AB/backend-animegg/routes"
	"github.com/3ggie-AB/backend-animegg/config"
	"os"
	"log"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 1 * 1024 * 1024 * 1024,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Mengizinkan origin tertentu
		AllowMethods: "GET,POST,PUT,DELETE",          // Mengizinkan metode GET, POST, PUT
		AllowHeaders: "Content-Type,Authorization", // Mengizinkan header tertentu
	}))

	config.ConnectDatabase()

	routes.WebRoutes(app)
	routes.ApiRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
        port = "3000"
    }
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal("❌ Error menjalankan server:", err)
	}
}
