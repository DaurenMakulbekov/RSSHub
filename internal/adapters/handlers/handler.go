package handlers

import (
	"RSSHub/internal/core/ports"
)

type handler struct {
	service ports.Service
}

func NewHandler(service ports.Service) *handler {
	var handler = &handler{
		service: service,
	}

	return handler
}
