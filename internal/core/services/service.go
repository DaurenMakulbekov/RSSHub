package services

import (
	"RSSHub/internal/core/domain"
	"RSSHub/internal/core/ports"
)

type service struct {
	postgres ports.PostgresRepository
}

func NewService(postgresRepository ports.PostgresRepository) *service {
	var service = &service{
		postgres: postgresRepository,
	}

	return service
}

func (service *service) AddFeed(feed domain.Feeds) error {
	var err = service.postgres.AddFeed(feed)

	return err
}
