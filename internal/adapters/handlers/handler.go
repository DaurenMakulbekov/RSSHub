package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"RSSHub/internal/core/domain"
	"RSSHub/internal/core/ports"
)

type handler struct {
	service ports.Service
}

func NewHandler(service ports.Service) *handler {
	handler := &handler{
		service: service,
	}

	return handler
}

func (handler *handler) AddFeedHandler(add domain.Add) {
	feed := domain.Feeds{
		Name: add.Name,
		Url:  add.Url,
	}

	err := handler.service.AddFeed(feed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func (handler *handler) FetchHandler() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Background process is already running\n")
		os.Exit(1)
	}
	defer listener.Close()

	handler.service.Fetch()
	log.Println("The background process for fetching feeds has started (interval = 3 minutes, workers = 3)")

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}

			go func() {
				defer conn.Close()

				buf := make([]byte, 512)

				sz, err := conn.Read(buf)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					return
				}

				var result domain.Commands

				decoder := json.NewDecoder(strings.NewReader(string(buf[:sz])))

				err = decoder.Decode(&result)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error decode: ", err)
				}

				if result.Name == "set-interval" {
					interval := handler.service.SetInterval(result.SetInterval.Duration)

					log.Printf("Interval of fetching feeds changed from %s minutes to %s minutes\n", interval, result.SetInterval.Duration)
				} else if result.Name == "set-workers" {
					workers := handler.service.SetWorkers(result.SetWorkers.Count)

					log.Printf("Number of workers changed from %d to %d\n", workers, result.SetWorkers.Count)
				}
			}()
		}
	}()

	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	defer signalCtxStop()

	<-signalCtx.Done()

	log.Println("Shutting down process...")
	handler.service.Stop()
	time.Sleep(5 * time.Second)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	log.Println("Graceful shutdown: aggregator stopped")
}

func (handler *handler) SetIntervalHandler(command domain.Commands) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	result, _ := json.Marshal(command)
	conn.Write(result)
}

func (handler *handler) SetWorkersHandler(command domain.Commands) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	result, _ := json.Marshal(command)
	conn.Write(result)
}

func (handler *handler) DeleteHandler(command domain.Delete) {
	feed := domain.Feeds{
		Name: command.Name,
	}

	err := handler.service.DeleteFeed(feed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func (handler *handler) ListHandler(command domain.List) {
	feeds, err := handler.service.GetFeeds()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	index := len(feeds) - command.Num
	if index >= 1 && index < len(feeds) {
		feeds = feeds[index:]
	}

	fmt.Println("# Available RSS Feeds")
	fmt.Println()

	for i := range feeds {
		fmt.Print(i+1, ".")
		fmt.Println(" Name: ", feeds[i].Name)
		fmt.Println("   URL: ", feeds[i].Url)
		fmt.Println("   Added: ", feeds[i].Created)
		fmt.Println()
	}
}

func (handler *handler) ArticlesHandler(command domain.ArticlesCommand) {
	articles, err := handler.service.GetArticles(command.FeedName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	index := len(articles) - command.Num
	if index >= 1 && index < len(articles) {
		articles = articles[index:]
	}

	fmt.Println("Feed: ", command.FeedName)
	fmt.Println()

	for i := range articles {
		fmt.Print(i+1, ".")
		fmt.Println(" ", "["+articles[i].Published+"]", articles[i].Title)
		fmt.Println("   ", articles[i].Link)
		fmt.Println()
	}
}
