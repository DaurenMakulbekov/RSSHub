package handlers

import (
	"RSSHub/internal/core/domain"
	"RSSHub/internal/core/ports"
	"fmt"
	"os"
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

func (handler *handler) AddFeedHandler(add domain.Add) {
	var feed = domain.Feeds{
		Name: add.Name,
		Url:  add.Url,
	}

	var err = handler.service.AddFeed(feed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func (handler *handler) FetchHandler() error {
	var err = handler.service.Fetch()

	return err
}
