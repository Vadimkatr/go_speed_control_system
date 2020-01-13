package apiserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vadimkatr/go_speed_control_system/internal/app/config"
	"github.com/vadimkatr/go_speed_control_system/internal/app/store/csvfilestore"
	"github.com/vadimkatr/go_speed_control_system/internal/app/utils"
)

func Start(cfg *config.Config) error {
	startTime, endTime, err := utils.ParseAvailableServerTime(cfg)
	if err != nil {
		return err
	}
	store := &csvfilestore.CSVFileStore{
		Datapass: cfg.Store.Dirpass,
	}
	srv := newServer(store, startTime, endTime)
	log.Printf("Start server at: %s:%s\n", cfg.Server.Host, cfg.Server.Port)

	return http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), srv)
}
