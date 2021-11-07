package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const MaxAge = 3600

func Static(app *fiber.App) {
	app.Static("/", "web/nazdik/out", fiber.Static{
		Compress:      true,
		ByteRange:     false,
		Browse:        false,
		Index:         "index.html",
		CacheDuration: time.Hour,
		MaxAge:        MaxAge,
		Next:          nil,
	})
}
