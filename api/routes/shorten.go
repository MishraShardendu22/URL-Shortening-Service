package routes

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ShardenduMishra22/url-shortener-service/api/database"
	"github.com/ShardenduMishra22/url-shortener-service/api/helpers"
	"github.com/ShardenduMishra22/url-shortener-service/api/models"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ShortenHandlerURL(c *fiber.Ctx) error {
	// Parses the incoming JSON request into a models.Request struct.
	var body models.Request

	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot Parse JSON"})
	}

	// Check the user's API quota by looking up their IP in Redis. If they have not been tracked,
	// set their quota with a 30-minute reset. If quota is exhausted, return an error with retry time.
	val, err := database.Client.Get(context.Background(), c.IP()).Result()
	if err == redis.Nil {
		// Initialize quota if not yet set
		_ = database.Client.Set(context.Background(), c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second)
	} else {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := database.Client.TTL(context.Background(), c.IP()).Result()
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate Limit Exceeded",
				"retry_after": limit / time.Minute / time.Nanosecond,
			})
		}
	}

	// Checks if body.URL is a valid URL format. If invalid, returns a 400 error response.
	if !govalidator.IsURL(body.URL) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// Prevents shortening URLs that point to the application’s own domain (self-referencing URLs).
	if !helpers.IsDifferentDomain(body.URL) {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "You Can't Hack this System (:",
		})
	}

	// Ensures the URL has a prefix of http or https.
	body.URL = helpers.EnsuredPrefixHTTP(body.URL)

	// The uuid library here generates a Universally Unique Identifier (UUID), which is a 128-bit number used to uniquely identify information in a distributed system.
	var id string
	// If no custom short URL is provided, generate a random one. Otherwise, use the provided custom short.
	if body.CustomShort == "" {
		id = uuid.New().String()[:8]
	} else {
		id = body.CustomShort
	}

	// Checks if the generated or provided short URL already exists. If it does, returns an error.
	val, _ = database.Client.Get(context.Background(), id).Result()
	if val != "" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "URL Custom Short is already in use",
		})
	}

	// Sets default expiry if not provided by the user (24 hours).
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	// Stores the URL in Redis with the provided or default expiry.
	err = database.Client.Set(context.Background(), id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect to server",
		})
	}

	// Constructs the response with remaining rate limit, reset time, and shortened URL.
	res := models.Response{
		Expiry:          body.Expiry,
		XRateLimitReset: 30,
		XRateRemaining:  10,
		URL:             body.URL,
		CustomShort:     "",
	}

	// Decreases the API quota for the user's IP by 1.
	database.Client.Decr(context.Background(), c.IP())
	val, _ = database.Client.Get(context.Background(), c.IP()).Result()
	res.XRateRemaining, _ = strconv.Atoi(val)

	// Retrieves the time remaining until the quota resets for the user's IP.
	ttl, _ := database.Client.TTL(context.Background(), c.IP()).Result()
	res.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	// Constructs the final shortened URL and assigns it to the response.
	res.CustomShort = os.Getenv("DOMAIN") + "/" + id

	// Returns the final response to the client.
	return c.Status(http.StatusOK).JSON(res)
}

// var (
// 	Ctx    = context.Background()
// 	Client *redis.Client
// )

// Request Struct
// Contains fields for the URL to shorten, an optional custom short code, and optional expiry in hours.
// type Request struct {
// 	URL         string        `json:"url"`
// 	CustomShort string        `json:"short"`
// 	Expiry      time.Duration `json:"expiry"`
// }

// Response Struct
// Holds the URL, custom short, expiry, and rate limit info.
// type Response struct {
// 	URL             string        `json:"url"`
// 	CustomShort     string        `json:"short"`
// 	Expiry          time.Duration `json:"expiry"`
// 	XRateRemaining  int           `json:"rate_limit"`
// 	XRateLimitReset time.Duration `json:"rate_limit_reset"`
// }

// TagRequest Struct
// Contains the ID of the shortened URL and an optional tag.
// type TagRequest struct {
// 	ShortID string `json:"shortID"`
// 	Tag     string `json:"tag"`
// }
