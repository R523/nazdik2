package handler

import (
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

const MaxAge = 3600

func Static(app *fiber.App, content fs.FS) {
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(content),
		PathPrefix:   "",
		NotFoundFile: "404.html",
		Browse:       false,
		MaxAge:       MaxAge,
		Index:        "index.html",
		Next:         nil,
	}))
}
