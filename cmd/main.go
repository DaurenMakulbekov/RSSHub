package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func helpMessage() {
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  rsshub COMMANDS [OPTIONS]")
	fmt.Println("  rsshub --help")
	fmt.Println("")
	fmt.Println("Common Commands:")
	fmt.Println("  add           add New RSS feed")
	fmt.Println("  set-interval  set RSS fetch interval")
	fmt.Println("  set-workers   set number of workers")
	fmt.Println("  list          list available RSS feeds")
	fmt.Println("  delete        delete RSS feed")
	fmt.Println("  articles      show latest articles")
	fmt.Println("  fetch         starts the background process that periodically fetches and processes RSS feeds using a worker pool")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = helpMessage
	flag.Parse()

	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGTSTP)
	defer signalCtxStop()

	//go func() {

	//}()

	<-signalCtx.Done()

	log.Println("Shutting down process...")
	time.Sleep(5 * time.Second)

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	log.Println("Graceful shutdown: aggregator stopped")
}
