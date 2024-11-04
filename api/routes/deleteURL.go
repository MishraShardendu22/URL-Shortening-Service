package routes

import (
	"net/http"

	"github.com/ShardenduMishra22/url-shortener-service/api/database"
	"github.com/gofiber/fiber/v2"
)

func DeleteURL(c *fiber.Ctx) error {
	shortID := c.Params("shortID")

	err := database.Client.Del(database.Ctx, shortID).Err()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to delete shortened link"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Shortened link deleted successfully"})
}
