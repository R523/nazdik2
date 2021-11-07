package ultrasonic

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
	"github.com/r523/nazdik/internal/store"
	"periph.io/x/conn/v3/gpio"
)

const (
	PulseInTimeout = 10 * time.Second
	Interval       = 1 * time.Second
)

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

type Ultrasonic struct {
	TriggerPin gpio.PinOut
	EchoPin    gpio.PinIn
}

func New(tp gpio.PinOut, ep gpio.PinIn) Ultrasonic {
	return Ultrasonic{
		TriggerPin: tp,
		EchoPin:    ep,
	}
}

func (u Ultrasonic) Read() (int64, error) {
	if err := u.TriggerPin.Out(gpio.Low); err != nil {
		return 0, fmt.Errorf("pinout %s %w", u.TriggerPin, err)
	}

	// nolint: gomnd
	time.Sleep(2 * time.Microsecond)

	if err := u.TriggerPin.Out(gpio.High); err != nil {
		return 0, fmt.Errorf("pinout %s %w", u.TriggerPin, err)
	}

	// nolint: gomnd
	time.Sleep(10 * time.Microsecond)

	duration, err := PulseIn(u.EchoPin, gpio.High, PulseInTimeout)
	if err != nil {
		return 0, fmt.Errorf("failed to pulse in %w", err)
	}

	// nolint: gomnd
	distance := duration.Microseconds() / 29 / 2

	return distance, nil
}

func (u Ultrasonic) Run(st *store.Distance, stop <-chan struct{}) {
	go func() {
		t := time.NewTicker(Interval)
		defer t.Stop()

		for {
			distance, err := u.Read()
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
}
