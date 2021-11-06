package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pterm/pterm"
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

	close(stop)
}
