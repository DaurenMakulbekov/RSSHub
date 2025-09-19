package services

import (
	"RSSHub/internal/core/domain"
	"RSSHub/internal/core/ports"
	"RSSHub/internal/infrastructure/config"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type service struct {
	postgres ports.PostgresRepository
	ticker   *time.Ticker
	interval time.Duration
	workers  int
	done     chan bool
}

func NewService(config *config.Config, postgresRepository ports.PostgresRepository) *service {
	i, _ := strconv.Atoi(config.Workers)
	s, err := time.ParseDuration(config.Interval)

	if err != nil || strings.Contains(config.Interval, "-") {
		fmt.Fprintln(os.Stderr, "Error: incorrect input")
	}

	var service = &service{
		postgres: postgresRepository,
		interval: time.Duration(s.Seconds()),
		workers:  i,
		done:     make(chan bool),
	}

	return service
}

func (service *service) AddFeed(feed domain.Feeds) error {
	var err = service.postgres.AddFeed(feed)

	return err
}

func (service *service) Start() {
	go func() {
		for {
			select {
			case <-service.done:
				return
			case <-service.ticker.C:
			}
		}
	}()
}

func (service *service) Stop() {
	service.ticker.Stop()
	close(service.done)
}

func (service *service) Reset() {
	s, err := time.ParseDuration("1s")

	if err != nil || strings.Contains("1s", "-") {
		fmt.Fprintln(os.Stderr, "Error: incorrect input")
	}

	service.interval = time.Duration(s.Seconds())
	service.ticker.Reset(service.interval * time.Second)
}

func (service *service) Fetch() error {
	service.ticker = time.NewTicker(service.interval * time.Second)

	service.Start()

	return nil
}
