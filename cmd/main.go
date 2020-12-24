package main

import (
	"flag"
	"log"

	"iziBookTest/internal/config"
	"iziBookTest/internal/web"
)

func main() {
	configPath := flag.String("c", "cmd/config.yml", "specify path to a config.yaml")
	flag.Parse()

	conf, err := config.GetConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	controller, err := web.NewController(conf)
	if err != nil {
		log.Fatal(err)
	}

	controller.StartService()
}
