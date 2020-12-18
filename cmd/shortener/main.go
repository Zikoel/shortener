package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/zikoel/shortener/pkg/persistence"
	"github.com/zikoel/shortener/pkg/shortener"
)

type shortenerParams struct {
	URL string `json:"URL" form:"URL" query:"URL"`
	Key string `json:"key" form:"key" query:"key"`
}

type config struct {
	ServerPort          int    `env:"SERVER_PORT" envDefault:"5000"`
	RedisHost           string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort           int    `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword       string `env:"REDIS_PASSWORD" envDefault:""`
	InMemoryPersistence bool   `env:"IN_MEMORY_PERSISTENCE" envDefault:"false"`
}

type stats struct {
	TargetURL  string `json:"targetURL"`
	VisitCount uint64 `json:"visitCount"`
}

// APIDocs A simple string to describe the endpoints
const APIDocs = `
endpoint: GET /:key
description:
	Redirect to the pointed key pointed URL or return 404 if URL associated with that key doesn't exist

endpoint: POST /api/urls
JSON Params:
	URL: The url to be register
	key: The suggested key
description:
	Register the long URL with the suggested key, if the key parameter is not provided an 8 char length key is
	automatically provided

endpoint: DELETE /api/urls/:key
description:
	Deleted the provided key

endpoint: GET /api/urls/:key/stats
description:
	Return a JSON that provide some key related stats

endpoint: GET /api/usage
description:
	Return this documentation
`

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic("Invalid arguments")
	}

	db := persistence.CreateRedisAdapter(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)
	short, err := shortener.CreateShortener(db, db, 1234)

	if err != nil {
		panic("Impossibile initialize shortener")
	}

	lookupHandler := func(c *fiber.Ctx) error {
		key := c.Params("key")
		url, err := short.URLFromKey(key)
		if err != nil {
			c.SendStatus(404)
			return err
		}
		c.Redirect(url)
		val, err := short.CountVisit(key)

		if err != nil {
			fmt.Printf("Error increment counter for key %s: %v\n", key, err)
			return err
		}

		fmt.Printf("key %s reach the value of %d\n", key, val)
		return nil
	}

	registerHandler := func(c *fiber.Ctx) error {
		params := new(shortenerParams)

		err = c.BodyParser(params)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(400)
		}

		key, err := short.KeyFromURL(params.URL, params.Key)

		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
			return err
		}

		c.SendString(key)
		fmt.Printf("Registered new key %s with target URL %s\n", key, params.URL)

		return nil
	}

	deleteHandler := func(c *fiber.Ctx) error {
		key := c.Params("key")
		err := short.DeleteURLByKey(key)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
			return err
		}
		fmt.Printf("Removed key %s\n", key)
		return nil
	}

	statsHandler := func(c *fiber.Ctx) error {
		key := c.Params("key")

		url, err := short.URLFromKey(key)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(404)
			return err
		}

		count, err := short.CollectStats(key)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
			return err
		}

		stats := stats{
			TargetURL:  url,
			VisitCount: count,
		}

		err = c.JSON(stats)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
			return err
		}

		return nil
	}

	usageHandler := func(c *fiber.Ctx) error {
		c.SendString(APIDocs)
		return nil
	}

	app := fiber.New()

	app.Use(cors.New())

	api := app.Group("/api")

	api.Post("/urls", registerHandler)
	api.Delete("/urls/:key", deleteHandler)
	api.Get("/urls/:key/stats", statsHandler)
	api.Get("/usage", usageHandler)

	app.Get("/:key", lookupHandler)

	app.Listen(fmt.Sprintf(":%d", cfg.ServerPort))
}
