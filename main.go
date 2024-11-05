package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)


type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Post("/shorten", func(c *fiber.Ctx) error {
		var data struct {
			URL string `json:"url"`
		}

		if err := json.Unmarshal(c.Body(), &data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		shortURL := createURL(data.URL)
		response := struct {
			ShortURL string `json:"short_url"`
		}{ShortURL: shortURL}

		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Get("/redirect/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		url, err := getURL(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("URL not found")
		}
		return c.Redirect(url.OriginalURL, fiber.StatusFound)
	})

	fmt.Println("Starting server on port 3000...")
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error on starting server:", err)
	}
}
