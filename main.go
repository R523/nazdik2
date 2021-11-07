package main

import (
	"embed"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
	"github.com/r523/nazdik/internal/http/handler"
	"github.com/r523/nazdik/internal/store"
	"github.com/r523/nazdik/internal/ultrasonic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

//go:embed web/nazdik/out
var content embed.FS

func main() {
	if err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Naz", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("dik", pterm.NewStyle(pterm.FgLightRed)),
	).Render(); err != nil {
		_ = err
	}

	// load all the drivers:
	if _, err := host.Init(); err != nil {
		pterm.Error.Printf("host initiation failed %s\n", err)

		return
	}

	st := new(store.Distance)

	app := fiber.New()

	d := handler.Distance{
		Store: st,
	}
	d.Register(app.Group("/api"))

	handler.Static(app, content)

	go func() {
		if err := app.Listen(":1378"); err != nil {
			pterm.Error.Printf("listen on port 1378 failed %s\n", err)
		}
	}()

	stop := make(chan struct{})

	tp := rpi.P1_13
	ep := rpi.P1_11
	ultrasonic.New(tp, ep).Run(st, stop)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	pterm.Info.Printf("Bye!\n")

	if err := app.Shutdown(); err != nil {
		pterm.Error.Printf("http server shutdown failed %s\n", err)
	}

	close(stop)
}