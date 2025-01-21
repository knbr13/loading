package main

import (
	"log"

	"github.com/knbr13/loading/internal/config"
	"github.com/knbr13/loading/internal/loader"
)

func main() {
	cfg, err := config.LoadConfig("testdata/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	options := loader.RequestOptions{
		Method:       cfg.Method,
		URL:          cfg.TargetURL,
		Headers:      cfg.Headers,
		Concurrency:  cfg.Concurrency,
		RequestCount: cfg.RequestCount,
	}

	loader.LoadTest(options)
}
