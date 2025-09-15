package services

import (
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
