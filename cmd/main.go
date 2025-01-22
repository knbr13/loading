package main

import (
	"errors"
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

	err = validateCfg(cfg)
	if err != nil {
		log.Fatalf("error validating configurations: %v", err)
	}

	var i *int
	if cfg.RequestCount != nil {
		j := int(*cfg.RequestCount)
		i = &j
	}

	options := loader.RequestOptions{
		Method:       cfg.Method,
		URL:          cfg.TargetURL,
		Headers:      cfg.Headers,
		Concurrency:  cfg.Concurrency,
		Duration:     cfg.Duration,
		RequestCount: i,
	}

	loader.LoadTest(options)
}

func validateCfg(c config.Config) error {
	if (c.Duration == nil && c.RequestCount == nil) ||
		(c.Duration != nil && c.RequestCount != nil) {
		return errors.New("provide a value for one of 'request_count' or 'duration' (one of them only, not both)")
	}
	if c.Concurrency <= 0 {
		return errors.New("'concurrency' value should be greater than zero")
	}
	return nil
}
