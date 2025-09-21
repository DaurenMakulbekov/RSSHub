package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"RSSHub/internal/core/domain"
	"RSSHub/internal/core/ports"
	"RSSHub/internal/infrastructure/config"
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

	service := &service{
		postgres: postgresRepository,
		interval: time.Duration(s.Seconds()),
		workers:  i,
		done:     make(chan bool),
	}

	return service
}

func (service *service) AddFeed(feed domain.Feeds) error {
	err := service.postgres.AddFeed(feed)

	return err
}

func (service *service) Worker(jobs <-chan domain.Feeds) {
	go func() {
		for feed := range jobs {
			response, err := http.Get(feed.Url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return
			}
			defer response.Body.Close()

			res, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}

			var result domain.RSSFeed

			err = xml.Unmarshal(res, &result)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}

			articles, err := service.postgres.GetArticles(feed)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}

			if len(articles) == 0 {
				err = service.postgres.WriteArticles(result.Channel.Item, feed)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			} else {
				var resultWrite []domain.RSSItem
				var has bool

				for i := range result.Channel.Item {
					has = false
					for j := range articles {
						if result.Channel.Item[i].Title == articles[j].Title {
							has = true
							break
						}
					}
					if !has {
						resultWrite = append(resultWrite, result.Channel.Item[i])
					}
				}

				err = service.postgres.WriteArticles(resultWrite, feed)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			}
		}
	}()
}

func (service *service) Start() {
	go func() {
		for {
			select {
			case <-service.done:
				return
			case <-service.ticker.C:
				feeds, err := service.postgres.GetFeeds()
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}

				jobs := make(chan domain.Feeds)

				for i := 0; i < service.workers; i++ {
					service.Worker(jobs)
				}

				for j := range feeds {
					jobs <- feeds[j]
				}
				close(jobs)
			}
		}
	}()
}

func (service *service) Stop() {
	service.ticker.Stop()
	close(service.done)
}

func (service *service) Fetch() {
	service.ticker = time.NewTicker(service.interval * time.Second)

	service.Start()
}

func (service *service) SetInterval(interval time.Duration) time.Duration {
	result := service.interval

	service.interval = interval
	service.ticker.Reset(service.interval * time.Second)

	return result
}

func (service *service) SetWorkers(workers int) int {
	result := service.workers

	service.workers = workers

	return result
}

func (service *service) DeleteFeed(feed domain.Feeds) error {
	err := service.postgres.DeleteFeed(feed)

	return err
}

func (service *service) GetFeeds() ([]domain.Feeds, error) {
	feeds, err := service.postgres.GetFeeds()
	if err != nil {
		return feeds, err
	}

	return feeds, nil
}

func (service *service) GetArticles(name string) ([]domain.Articles, error) {
	articles, err := service.postgres.GetArticlesByName(name)
	if err != nil {
		return articles, err
	}

	return articles, nil
}
