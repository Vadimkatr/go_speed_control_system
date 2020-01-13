package record

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/vadimkatr/go_speed_control_system/internal/app/utils"
)

var (
	minSpeed = float32(0.0)
	maxSpeed = float32(400.0)
	minDate  = time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC)
	maxDate  = time.Date(2030, time.January, 1, 1, 0, 0, 0, time.UTC)
)

type Record struct {
	RecordId      string
	DateTime      time.Time `json:"datetime"`
	VehicleNumber string    `json:"vehicle_number"`
	Speed         float32   `json:"speed"`
}

func CreateRecordId(dt time.Time, vn string, s float32) string {
	speed := fmt.Sprintf("%f", s)
	id := sha256.Sum256([]byte(fmt.Sprintf("%s %s %s", dt.String(), vn, speed)))
	return hex.EncodeToString(id[:])
}

func (r *Record) ToMap() map[string]string {
	return map[string]string{
		"datetime":       r.DateTime.Format(utils.DateTimeReqFomat),
		"vehicle_number": r.VehicleNumber,
		"speed":          fmt.Sprintf("%v", r.Speed),
	}
}

func CreateRecord(dt time.Time, vn string, s float32) (*Record, error) {
	if s < minSpeed || s > maxSpeed {
		return nil, ErrValidateRecSpeed
	}
	if vn == "" {
		return nil, ErrValidateRecVehNum
	}
	if dt.Before(minDate) || dt.After(maxDate) {
		return nil, ErrValidateRecDatetime
	}

	return &Record{
		RecordId:      CreateRecordId(dt, vn, s),
		DateTime:      dt,
		VehicleNumber: vn,
		Speed:         s,
	}, nil
}
