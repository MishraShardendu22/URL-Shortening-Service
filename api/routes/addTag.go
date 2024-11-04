package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ShardenduMishra22/url-shortener-service/api/database"
	"github.com/ShardenduMishra22/url-shortener-service/api/models"
	"github.com/gofiber/fiber/v2"
)

func AddTag(c *fiber.Ctx) error {
	var tagRequest models.TagRequest
	if err := c.BodyParser(&tagRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	shortID := tagRequest.ShortID
	tag := tagRequest.Tag

	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Data not found for the given shortID"})
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		data = make(map[string]interface{})
		data["data"] = val
	}

	var tags []string
	if existingTags, ok := data["tags"].([]interface{}); ok {
		for _, t := range existingTags {
			if strTag, ok := t.(string); ok {
				tags = append(tags, strTag)
			}
		}
	}

	for _, existingTag := range tags {
		if existingTag == tag {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Tag already exists"})
		}
	}

	tags = append(tags, tag)
	data["tags"] = tags

	updatedData, err := json.Marshal(data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to marshal updated data"})
	}

	err = database.Client.Set(database.Ctx, shortID, updatedData, 0).Err()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update database"})
	}

	return c.Status(http.StatusOK).JSON(data)
}
