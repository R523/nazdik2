package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const MaxAge = 3600

func Static(app *fiber.App) {
	app.Static("/", "../../web/nazdik/out", fiber.Static{
		Compress:      true,
		ByteRange:     false,
		Browse:        false,
		CacheDuration: time.Hour,
		MaxAge:        MaxAge,
		Index:         "index.html",
		Next:          nil,
	})
}
