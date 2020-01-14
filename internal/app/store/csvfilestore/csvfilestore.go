package csvfilestore

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/vadimkatr/go_speed_control_system/internal/app/record"
	"github.com/vadimkatr/go_speed_control_system/internal/app/utils"
)

var (
	ErrNoRecords = errors.New("there are no records for this request")
    mu sync.Mutex
)

type CSVFileStore struct {
	Datapass string
}

func (s *CSVFileStore) SaveRecord(r *record.Record) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(
		fmt.Sprintf("%s/%s.csv", s.Datapass, r.DateTime.Format(utils.SavedFileTimeFormat)),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		os.ModePerm,
	)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{
		r.RecordId,
		r.DateTime.Format(utils.DateTimeReqFomat),
		r.VehicleNumber,
		fmt.Sprintf("%f", r.Speed),
	})
	if err != nil {
		return "", err
	}

	return r.RecordId, nil
}

func (s *CSVFileStore) CarsExceedingSpeed(date time.Time, speed float32) ([]*record.Record, error) {
	filepath := fmt.Sprintf("%s/%s.csv", s.Datapass, date.Format(utils.SavedFileTimeFormat))
	if !fileExists(filepath) {
		return nil, ErrNoRecords
	}

	file, err := os.OpenFile(
		filepath,
		os.O_RDONLY,
		os.ModePerm,
	)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records := make([]*record.Record, 0)
	parser := csv.NewReader(file)
	for {
		r, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		rec, err := fromCsvRecToStruct(r)
		if err != nil {
			return nil, err
		}
		if date.Truncate(24*time.Hour).Equal(rec.DateTime.Truncate(24*time.Hour)) && rec.Speed >= speed {
			records = append(records, rec)
		}
	}

	return records, nil
}

func (s *CSVFileStore) MinAndMaxSpeed(date time.Time) (*record.Record, *record.Record, error) {
	filepath := fmt.Sprintf("%s/%s.csv", s.Datapass, date.Format(utils.SavedFileTimeFormat))
	if !fileExists(filepath) {
		return nil, nil, ErrNoRecords
	}

	file, err := os.OpenFile(
		filepath,
		os.O_RDONLY,
		os.ModePerm,
	)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	minRec := record.Record{
		Speed: -1,
	}
	maxRec := record.Record{
		Speed: -1,
	}
	parser := csv.NewReader(file)
	for {
		r, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		rec, err := fromCsvRecToStruct(r)
		if err != nil {
			return nil, nil, err
		}
		if date.Truncate(24 * time.Hour).Equal(rec.DateTime.Truncate(24 * time.Hour)) {
			if minRec.Speed == -1 && maxRec.Speed == -1 {
				minRec = *rec
				maxRec = *rec
				continue
			}
			if rec.Speed < minRec.Speed {
				minRec = *rec
			}
			if rec.Speed > maxRec.Speed {
				maxRec = *rec
			}
		}
	}

	if minRec.Speed == -1 && maxRec.Speed == -1 {
		return nil, nil, errors.New("no records for this day")
	}

	return &minRec, &maxRec, nil
}

func fromCsvRecToStruct(r []string) (*record.Record, error) {
	v := reflect.ValueOf(record.Record{})
	if len(r) != v.NumField() {
		return nil, errors.New("parsing csv file error: invalid number of fields")
	}

	dt, err := time.Parse(utils.DateTimeReqFomat, r[1])
	if err != nil {
		return nil, errors.New("parsing csv file error: " + err.Error())
	}

	speed, err := strconv.ParseFloat(r[3], 32)
	if err != nil {
		return nil, errors.New("parsing csv file error: " + err.Error())
	}

	return &record.Record{
		RecordId:      r[0],
		DateTime:      dt,
		VehicleNumber: r[2],
		Speed:         float32(speed),
	}, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
