package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/vadimkatr/go_speed_control_system/internal/app/record"
	"github.com/vadimkatr/go_speed_control_system/internal/app/store"
	"github.com/vadimkatr/go_speed_control_system/internal/app/utils"
)

var (
	ErrServerIsNotAvailable = errors.New("server is not available now. Try another time")
)

type server struct {
	router             *mux.Router
	logger             *CustomLogger
	store              store.Store
	availableStartTime time.Time
	availableEndTime   time.Time
}

func newServer(store store.Store, startTime, endTime time.Time) *server {
	s := &server{
		router:             mux.NewRouter(),
		logger:             newLogger(os.Stdout, os.Stdout, os.Stdout),
		store:              store,
		availableStartTime: startTime,
		availableEndTime:   endTime,
	}

	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.router)
	loggedRouter.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc(
		"/create_record",
		s.availableServerMiddleware(s.handleRecordCreate()),
	).Methods("POST")
	s.router.HandleFunc(
		"/get_exceeding_speed",
		s.availableServerMiddleware(s.handleGetCarsExceedingSpeedForDate()),
	).Methods("GET")
	s.router.HandleFunc(
		"/get_minmax_record",
		s.availableServerMiddleware(s.handleMinAndMaxSpeed()),
	).Methods("GET")
}

func (s *server) availableServerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()
		if !(utils.MinutesOfDay(timeNow) <= utils.MinutesOfDay(s.availableEndTime) &&
			utils.MinutesOfDay(timeNow) >= utils.MinutesOfDay(s.availableStartTime)) {
			s.logger.Error.Printf("error check server available time: %v\n", ErrServerIsNotAvailable)
			s.error(w, r, http.StatusUnauthorized, ErrServerIsNotAvailable)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (s *server) handleRecordCreate() http.HandlerFunc {
	type request struct {
		DateTime      utils.DateTimeRF `json:"datetime"`
		VehicleNumber string           `json:"vehicle_number"`
		Speed         float32          `json:"speed,string"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error.Printf("error while creating record: %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// validate and create record
		dt := time.Time(req.DateTime)
		newRecord, err := record.CreateRecord(dt, req.VehicleNumber, req.Speed)
		if err != nil {
			s.logger.Error.Printf("error while validating or creating record: %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// save record to store
		recordId, err := s.store.SaveRecord(newRecord)
		if err != nil {
			s.logger.Error.Printf("error while saving record: %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Info.Printf("save record with id: %s\n", recordId)
		s.respond(w, r, http.StatusCreated, map[string]string{"recordId": recordId})
	}
}

func (s *server) handleGetCarsExceedingSpeedForDate() http.HandlerFunc {
	type request struct {
		Date  utils.DateRF `json:"date"`
		Speed float32      `json:"speed,string"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error.Printf("error while getting exceed speed cars: %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// get cars from store
		dt := time.Time(req.Date)
		records, err := s.store.CarsExceedingSpeed(dt, req.Speed)
		if err != nil {
			s.logger.Error.Printf("error while getting exceed speed cars: %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		res := make([]map[string]string, 0)
		for _, rec := range records {
			res = append(res, rec.ToMap())
		}

		s.logger.Info.Printf("get exceed speed cars\n")
		s.respond(w, r, http.StatusCreated, res)
	}
}

func (s *server) handleMinAndMaxSpeed() http.HandlerFunc {
	type request struct {
		Date utils.DateRF `json:"date"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.logger.Error.Printf("error while getting minmax records: %v\n", err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		dt := time.Time(req.Date)
		recMin, recMax, err := s.store.MinAndMaxSpeed(dt)
		if err != nil {
			s.logger.Error.Printf("error while getting minmax records: %v\n", err)
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.logger.Info.Printf("get minmax records\n")
		s.respond(w, r, http.StatusCreated, []map[string]string{recMin.ToMap(), recMax.ToMap()})
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
