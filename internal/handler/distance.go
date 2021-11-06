package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/r523/nazdik/internal/store"
)

type Distance struct {
	Store *store.Distance
}

func (d *Distance) Get(c *fiber.Ctx) error {
	// nolint: wrapcheck
	return c.Status(http.StatusOK).JSON(d.Store.Get())
}

func (d *Distance) Register(g fiber.Router) {
	g.Get("/distance", d.Get)
}
