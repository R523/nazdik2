package alarm

import (
	"time"

	"github.com/r523/nazdik/internal/store"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3/rpi"
)

const Interval = 2 * time.Second

type Alarm struct {
	Store     *store.Distance
	Threshold int64
}

func (a *Alarm) Run() {
	for {
		if a.Store.Get() < a.Threshold {
			if err := rpi.P1_37.Out(gpio.High); err != nil {
				_ = err
			}
		} else {
			if err := rpi.P1_37.Out(gpio.High); err != nil {
				_ = err
			}
		}

		<-time.After(Interval)
	}
}
