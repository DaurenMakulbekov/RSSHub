package ports

import (
	"RSSHub/internal/core/domain"
)

type PostgresRepository interface {
	AddFeed(feed domain.Feeds) error
}

type Service interface {
	AddFeed(feed domain.Feeds) error
	Fetch() error
}
