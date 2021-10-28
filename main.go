package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	Interval        = 1 * time.Second
	PulseInTimesout = 10 * time.Second

	// Thresholds in cm
	Threshold0 = 10
	Threshold1 = 20
	Threshold2 = 30
	Threshold3 = 40
)

func SetLEDs(lvls [4]gpio.Level) error {
	leds := [4]gpio.PinOut{rpi.P1_29, rpi.P1_31, rpi.P1_33, rpi.P1_32}

	for i := 0; i < 4; i++ {
		if err := leds[i].Out(lvls[i]); err != nil {
			return fmt.Errorf("failed to set output on pin %s %w", leds[i], err)
		}
	}

	return nil
}

func PulseIn(pin gpio.PinIn, lvl gpio.Level, t time.Duration) (time.Duration, error) {
	var e1, e2 gpio.Edge

	if lvl == gpio.High {
		e1 = gpio.RisingEdge
		e2 = gpio.FallingEdge
	} else {
		e1 = gpio.FallingEdge
		e2 = gpio.RisingEdge
	}

	if err := pin.In(gpio.PullNoChange, e1); err != nil {
		return 0, fmt.Errorf("pin %s in setup failed %w", pin, err)
	}

	pin.WaitForEdge(t)

	now := time.Now()

	if err := pin.In(gpio.PullNoChange, e2); err != nil {
		return 0, fmt.Errorf("pin %s in setup failed %w", pin, err)
	}

	pin.WaitForEdge(t)

	return time.Since(now), nil
}

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

	if err := SetLEDs([4]gpio.Level{gpio.Low, gpio.Low, gpio.Low, gpio.Low}); err != nil {
		pterm.Error.Printf("set leds failed %s\n", err)

		return
	}

	t := time.NewTicker(Interval)
	defer t.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		tp := rpi.P1_13
		ep := rpi.P1_11

		if err := tp.Out(gpio.Low); err != nil {
			pterm.Error.Printf("pinout %s %s\n", tp, err)

			return
		}

		// nolint: gomnd
		time.Sleep(2 * time.Microsecond)

		if err := tp.Out(gpio.High); err != nil {
			pterm.Error.Printf("pinout %s %s\n", tp, err)

			return
		}

		// nolint: gomnd
		time.Sleep(10 * time.Microsecond)

		duration, err := PulseIn(ep, gpio.High, PulseInTimesout)
		if err != nil {
			pterm.Error.Printf("failed to pulse in %s", err)

			return
		}

		// nolint: gomnd
		distance := duration.Microseconds() / 29 / 2

		pterm.Info.Printf("there is an object in %d cm\n", distance)

		switch {
		case distance <= Threshold0:
			if err := SetLEDs([4]gpio.Level{gpio.High, gpio.Low, gpio.Low, gpio.Low}); err != nil {
				pterm.Error.Printf("set leds failed %s\n", err)

				return
			}
		case distance > Threshold0 && distance <= Threshold1:
			if err := SetLEDs([4]gpio.Level{gpio.Low, gpio.High, gpio.Low, gpio.Low}); err != nil {
				pterm.Error.Printf("set leds failed %s\n", err)

				return
			}
		case distance > Threshold1 && distance <= Threshold2:
			if err := SetLEDs([4]gpio.Level{gpio.Low, gpio.Low, gpio.High, gpio.Low}); err != nil {
				pterm.Error.Printf("set leds failed %s\n", err)

				return
			}
		case distance > Threshold2 && distance <= Threshold3:
			if err := SetLEDs([4]gpio.Level{gpio.Low, gpio.Low, gpio.Low, gpio.High}); err != nil {
				pterm.Error.Printf("set leds failed %s\n", err)

				return
			}
		case distance > Threshold3:
			if err := SetLEDs([4]gpio.Level{gpio.Low, gpio.Low, gpio.Low, gpio.Low}); err != nil {
				pterm.Error.Printf("set leds failed %s\n", err)

				return
			}
		}

		select {
		case <-t.C:
		case <-quit:
			return
		}
	}
}
