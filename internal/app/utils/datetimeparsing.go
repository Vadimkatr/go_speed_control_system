package utils

import (
	"encoding/json"
	"time"

	"github.com/vadimkatr/go_speed_control_system/internal/app/config"
)

const (
	ServerTimeFormat    = "15:04"
	SavedFileTimeFormat = "02_01_2006"
	DateTimeReqFomat    = "02.01.2006 15:04:05"
	DateReqFormat       = "02.01.2006"
)

type DateTimeRF time.Time
type DateRF time.Time

// implement UnmarshalJSON interface
func (dt *DateTimeRF) UnmarshalJSON(b []byte) error {
	//s := strings.Trim(string(b), "\"")
	s := string(b)[1 : len(b)-1]
	t, err := time.Parse(DateTimeReqFomat, s)
	if err != nil {
		return err
	}
	*dt = DateTimeRF(t)
	return nil
}

// implement MarshalJSON interface
func (dt DateTimeRF) MarshalJSON() ([]byte, error) {
	return json.Marshal(dt)
}

// implement UnmarshalJSON interface
func (dt *DateRF) UnmarshalJSON(b []byte) error {
	//s := strings.Trim(string(b), "\"")
	s := string(b)[1 : len(b)-1]
	t, err := time.Parse(DateReqFormat, s)
	if err != nil {
		return err
	}
	*dt = DateRF(t)
	return nil
}

// ParseAvailableServerTime - parse time from config
func ParseAvailableServerTime(cfg *config.Config) (time.Time, time.Time, error) {
	startTime, err := time.Parse(ServerTimeFormat, cfg.Service.TimeStart)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, err := time.Parse(ServerTimeFormat, cfg.Service.TimeEnd)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return startTime, endTime, nil
}
