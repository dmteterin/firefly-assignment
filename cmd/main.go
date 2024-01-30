package main

import (
	"os"

	"github.com/dmteterin/firefly-assignment/internal/bank"
	"github.com/dmteterin/firefly-assignment/internal/config"
	"github.com/dmteterin/firefly-assignment/internal/crawler"
	"github.com/dmteterin/firefly-assignment/internal/validator"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	done := make(chan struct{})

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not load environment variables")
	}

	validator := validator.New(&cfg)

	bank, err := bank.New(&cfg, validator, done, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not create word bank")
	}

	crawler, err := crawler.New(&cfg, bank, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not create crawler")
	}

	err = crawler.RunScrapingQueue()
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not run scraping queue")
	}

	<-done
}
