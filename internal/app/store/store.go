package store

import (
	"time"

	"github.com/vadimkatr/go_speed_control_system/internal/app/record"
)

type Store interface {
	SaveRecord(*record.Record) (string, error)
	CarsExceedingSpeed(time.Time, float32) ([]*record.Record, error)
	MinAndMaxSpeed(time.Time) (*record.Record, *record.Record, error)
}
