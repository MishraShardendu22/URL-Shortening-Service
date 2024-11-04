package main

import (
	"fmt"
	"os"

	"github.com/ShardenduMishra22/url-shortener-service/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// Test Route Set-up
	TestRouteSetUp(app)

	// Setting Up Routes
	SetUpRoutes(app)

	// Listening To The Port
	ListenToThePort(app)
}

func HandleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// Set up CORS to prevent CSRF attacks
func SetUpCORS(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE",
	}))
}

// Listening To The Port
func ListenToThePort(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Listening to port: " + port)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// Route to check if the server is running or not
func TestRouteSetUp(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"Message": "This is Working!!"})
	})
}

// Setting Up Routes for the application
func SetUpRoutes(app *fiber.App) {

	app.Get("/api/v1/:shortID", routes.GetByShortID)

	app.Post("/api/v1/addTag", routes.AddTag)

	app.Post("/api/v1", routes.ShortenHandlerURL)

	app.Put("/api/v1/:shortID", routes.EditURL)

	app.Delete("/api/v1/:shortID", routes.DeleteURL)
}
