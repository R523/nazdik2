package store

import "sync"

type Distance struct {
	value int64
	lock  sync.RWMutex
}

func (d *Distance) Get() int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.value
}

func (d *Distance) Set(value int64) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.value = value
}
