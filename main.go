package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("This is a URL Shortening Service!!")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	// Initialize the server
	app := fiber.New()
	
	// Setting Up CORS
	SetUpCORS(app)
	
	// Listening To The Port
	ListenToThePort(app)
}

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func SetUpCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE",
	}))
}

func ListenToThePort(app *fiber.App) {
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + os.Getenv(port)))
	fmt.Println("Listening to port: " + port)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}
}