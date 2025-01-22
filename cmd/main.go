package main

import (
	"flag"
	"log"

	"github.com/knbr13/loading/internal/config"
	"github.com/knbr13/loading/internal/loader"
)

func main() {
	cfgFile := flag.String("c", "", "the .yaml configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgFile)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
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
