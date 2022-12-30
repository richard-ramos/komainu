package main

import (
	"github.com/richard-ramos/komainu/pkg/config"
	"github.com/richard-ramos/komainu/pkg/services"
)

func main() {

	cfg := &config.Config{
		DataDir: "./",
	}

	services.App(cfg).Run()
}
