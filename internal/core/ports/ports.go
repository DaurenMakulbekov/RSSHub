package ports

import (
	"time"

	"RSSHub/internal/core/domain"
)

type PostgresRepository interface {
	AddFeed(feed domain.Feeds) error
	GetFeeds() ([]domain.Feeds, error)
	WriteArticles(articles []domain.RSSItem, feed domain.Feeds) error
	GetArticles(feed domain.Feeds) ([]domain.Articles, error)
	DeleteFeed(feed domain.Feeds) error
	GetArticlesByName(name string) ([]domain.Articles, error)
}

type Service interface {
	AddFeed(feed domain.Feeds) error
	Fetch()
	Stop()
	SetInterval(interval time.Duration) time.Duration
	SetWorkers(workers int) int
	DeleteFeed(feed domain.Feeds) error
	GetFeeds() ([]domain.Feeds, error)
	GetArticles(name string) ([]domain.Articles, error)
}
