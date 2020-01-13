package utils

import "time"

func MinutesOfDay(dt time.Time) int {
	return dt.Minute() + dt.Hour()*60
}
