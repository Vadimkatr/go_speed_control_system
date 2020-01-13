package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/vadimkatr/go_speed_control_system/internal/app/apiserver"
	"github.com/vadimkatr/go_speed_control_system/internal/app/config"
)

func main() {
	f, err := os.Open("config/conf.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg *config.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(cfg); err != nil {
		log.Fatal(err)
	}
}
