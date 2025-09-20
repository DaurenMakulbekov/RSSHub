package handlers

import (
	"RSSHub/internal/core/domain"
	"RSSHub/internal/core/ports"
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
				//fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}

			go func() {
				defer conn.Close()

				var buf = make([]byte, 512)

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
					var interval = handler.service.SetInterval(result.SetInterval.Duration)

					log.Printf("Interval of fetching feeds changed from %s minutes to %s minutes\n", interval, result.SetInterval.Duration)
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

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

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
