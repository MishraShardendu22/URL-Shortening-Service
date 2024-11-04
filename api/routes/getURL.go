package routes

import (
	"net/http"

	"github.com/ShardenduMishra22/url-shortener-service/api/database"
	"github.com/gofiber/fiber/v2"
)

func GetByShortID(c *fiber.Ctx) error {
	shortID := c.Params("shortID")

	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Data not found for the given tagID"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": val})
}
