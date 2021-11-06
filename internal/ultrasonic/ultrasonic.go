package ultrasonic

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/gpio"
)

const PulseInTimeout = 10 * time.Second

// PulseIn measures the duration of specified level on the given pin.
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

func Read(tp, ep gpio.PinIO) (int64, error) {
	if err := tp.Out(gpio.Low); err != nil {
		return 0, fmt.Errorf("pinout %s %w", tp, err)
	}

	// nolint: gomnd
	time.Sleep(2 * time.Microsecond)

	if err := tp.Out(gpio.High); err != nil {
		return 0, fmt.Errorf("pinout %s %w", tp, err)
	}

	// nolint: gomnd
	time.Sleep(10 * time.Microsecond)

	duration, err := PulseIn(ep, gpio.High, PulseInTimeout)
	if err != nil {
		return 0, fmt.Errorf("failed to pulse in %w", err)
	}

	// nolint: gomnd
	distance := duration.Microseconds() / 29 / 2

	return distance, nil
}
