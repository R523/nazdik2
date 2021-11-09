package alarm

import (
	"time"

	"github.com/r523/nazdik/internal/pcf8574"
	"github.com/r523/nazdik/internal/store"
)

const (
	Interval    = 2 * time.Second
	BuzzTimeout = 1 * time.Second

	onCommand  = 0b0000_0000
	offCommand = 0b1111_0000
)

type Alarm struct {
	Store     *store.Distance
	Threshold int64
}

func New(st *store.Distance, t int64) *Alarm {
	return &Alarm{
		Store:     st,
		Threshold: t,
	}
}

func (a *Alarm) Run() {
	for {
		if a.Store.Get() < a.Threshold {
			p, err := pcf8574.New("/dev/i2c-1")
			if err != nil {
				_ = err

				continue
			}

			if err := p.Write(onCommand); err != nil {
				_ = err

				continue
			}

			time.Sleep(BuzzTimeout)

			if err := p.Write(offCommand); err != nil {
				_ = err

				continue
			}
		}

		<-time.After(Interval)
	}
}
