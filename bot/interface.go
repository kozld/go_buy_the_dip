package bot

import "time"

type Bot interface {
	HandleCandle(time.Time, float64) error
	Start() float64
}
