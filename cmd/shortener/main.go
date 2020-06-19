package main

import (
	"fmt"

	"github.com/gofiber/fiber"
	"github.com/zikoel/shortener/pkg/persistence"
	"github.com/zikoel/shortener/pkg/shortener"
)

type shortenerParams struct {
	URL string `json:"URL" form:"URL" query:"URL"`
	Key string `json:"key" form:"key" query:"key"`
}

func main() {
	db := persistence.CreateRedisAdapter()
	short, err := shortener.CreateShortener(db, db, 1234)

	if err != nil {
		panic("Impossibile initialize short")
	}

	lookupHandler := func(c *fiber.Ctx) {
		key := c.Params("key")
		url, err := short.URLFromKey(key)
		if err != nil {
			c.SendStatus(404)
			return
		}
		c.Redirect(url)
		val, err := short.CountVisit(key)

		if err != nil {
			fmt.Printf("Error increment counter for key %s: %v\n", key, err)
			return
		}

		fmt.Printf("key %s reach the value of %d\n", key, val)
	}
	registerHandler := func(c *fiber.Ctx) {
		params := new(shortenerParams)

		err = c.BodyParser(params)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
		}

		key, err := short.KeyFromURL(params.URL, params.Key)

		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
			return
		}

		c.Send(fmt.Sprintf("Registered url %s", key))
	}
	deleteHandler := func(c *fiber.Ctx) {
		err := short.DeleteURLByKey(c.Params("key"))
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
		}
	}
	statsHandler := func(c *fiber.Ctx) {
		key := c.Params("key")

		count, err := short.CollectStats(key)
		if err != nil {
			fmt.Printf("%v\n", err)
			c.SendStatus(500)
			return
		}

		c.Send(fmt.Sprintf("Count: %d", count))
	}

	app := fiber.New()
	api := app.Group("/api")

	api.Post("/register", registerHandler)
	api.Delete("/:key", deleteHandler)
	api.Get("/stats/:key", statsHandler)

	app.Get("/:key", lookupHandler)

	app.Listen(5000)
}
