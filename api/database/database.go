package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       dbNo,
	})
	return rdb
}

func InitializeClient() {
	Client = CreateClient(0)
}

// Redis allows you to create multiple databases within a single Redis server instance. Each of these databases is indexed with integers starting from 0. By default, Redis has 16 databases (0 to 15), but this can be configured.

// Key Characteristics:

// Each database is isolated; keys in one database do not affect keys in another.
// You specify which database to use by providing its index when creating a client connection.
// Example
// Let's say you have the following two databases in your Redis server:

// Database 0:

// Keys: user:1 -> {"name": "Alice"}
// Keys: user:2 -> {"name": "Bob"}
// Database 1:

// Keys: order:1 -> {"item": "Laptop", "quantity": 1}
// Keys: order:2 -> {"item": "Phone", "quantity": 2}
// When you connect to Database 0 with CreateClient(0), you can access the user:1 and user:2 keys, but you cannot see the order:1 and order:2 keys in Database 1 unless you connect to that database (using CreateClient(1)).