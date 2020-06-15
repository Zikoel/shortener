package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber"
	"github.com/teris-io/shortid"
)

func addNewURL(ctx *context.Context, client *redis.Client, generator *shortid.Shortid, url string) (string, error) {
	code, err := generator.Generate()

	if err != nil {
		return "", err
	}

	err = client.Set(code, url, 0).Err()

	if err != nil {
		return "", err
	}

	err = client.Set(fmt.Sprintf("%s/count", code), 1, 0).Err()

	if err != nil {
		return "", err
	}

	return code, nil
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	app := fiber.New()

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		panic("Random generator init error")
	}

	app.Post("/api/:url", func(c *fiber.Ctx) {
		ctx := context.Background()
		url, err := addNewURL(&ctx, rdb, sid, c.Params("url"))

		if err != nil {
			fmt.Printf("err: %v\n", err)
			c.SendStatus(500)
			return
		}

		c.Send(fmt.Sprintf("Registered url %s", url))
	})

	app.Delete("/api/:shortURL", func(c *fiber.Ctx) {
		r := rdb.Del(c.Params("shortURL"))

		if r.Err() != nil {
			c.SendStatus(500)
			return
		}

		c.Send("Removed")
	})

	app.Get("/api/stats/:shortURL", func(c *fiber.Ctx) {
		r := rdb.Get(fmt.Sprintf("%s/count", c.Params("shortURL")))

		if r.Err() != nil {
			c.Send("No data available")
			return
		}

		count, err := r.Int64()

		if err != nil {
			c.Send("Error on counter")
			return
		}

		c.Send(fmt.Sprintf("The URL was hitted %d times", count))
	})

	app.Get("/:shortURL", func(c *fiber.Ctx) {
		r := rdb.Get(c.Params("shortURL"))
		if r.Err() != nil {
			c.SendStatus(404)
			return
		}

		rdb.Incr(fmt.Sprintf("%s/count", c.Params("shortURL")))

		c.Redirect(r.Val())
	})

	app.Listen(3000)
}
