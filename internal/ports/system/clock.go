package system

import (
	"github.com/lewis-od/wavelength/internal/clock"
	"time"
)

type systemClock struct {}

func NewClock() clock.Clock {
	return systemClock{}
}

func (systemClock) Now() time.Time {
	return time.Now()
}
