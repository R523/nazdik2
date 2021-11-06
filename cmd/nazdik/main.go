package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pterm/pterm"
	"github.com/r523/nazdik/internal/handler"
	"github.com/r523/nazdik/internal/store"
	"github.com/r523/nazdik/internal/ultrasonic"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	Interval = 1 * time.Second
)

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

	var st store.Distance

	app := fiber.New()

	d := handler.Distance{
		Store: &st,
	}
	d.Register(app.Group("/api"))

	go func() {
		if err := app.Listen(":1378"); err != nil {
			pterm.Error.Printf("listen on port 1378 failed %s\n", err)
		}
	}()

	stop := make(chan struct{})

	go func() {
		t := time.NewTicker(Interval)
		defer t.Stop()

		for {
			tp := rpi.P1_13
			ep := rpi.P1_11

			distance, err := ultrasonic.Read(tp, ep)
			if err != nil {
				pterm.Error.Printf("cannot read from ultrasonic %s\n", err)
			}

			pterm.Info.Printf("there is an object in %d cm\n", distance)

			st.Set(distance)

			select {
			case <-t.C:
			case <-stop:
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if err := app.Shutdown(); err != nil {
		pterm.Error.Printf("http server shutdown failed %s\n", err)
	}

	close(stop)
}
