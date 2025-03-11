package main

import (
	"github.com/M-kos/anti_brutforce/internal/api"
	"github.com/M-kos/anti_brutforce/internal/config"
)

func main() {
	run() // TODO: graceful shutdown
}

func run() error {
	conf := config.LoadConfig()
	err := api.Run(conf)

	return err
}
