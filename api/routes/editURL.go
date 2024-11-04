package routes

import (
	"net/http"
	"time"

	"github.com/ShardenduMishra22/url-shortener-service/api/database"
	"github.com/ShardenduMishra22/url-shortener-service/api/models"
	"github.com/gofiber/fiber/v2"
)

func EditURL(c *fiber.Ctx) error {
	shortID := c.Params("shortID")
	var body models.Request

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Check if the shortID exists in the database
	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil || val == "" {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "ShortID does not exist"})
	}

	// Update the URL associated with the shortID
	err = database.Client.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to update shortened link"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Shortened link updated successfully"})
}
