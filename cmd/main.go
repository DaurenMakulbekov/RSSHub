package main

import (
	"RSSHub/internal/adapters/handlers"
	"RSSHub/internal/adapters/repositories/postgres"
	"RSSHub/internal/core/services"
	"RSSHub/internal/infrastructure/config"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
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

	var args = os.Args
	var config = config.NewAppConfig()

	var postgresRepository = postgres.NewPostgresRepository(config.DB)
	var service = services.NewService(config.Config, postgresRepository)
	var handler = handlers.NewHandler(service)

	command, err := handler.GetCommand(args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch command.Name {
	case "add":
		fmt.Println("Add: ", command.Add)
		handler.AddFeedHandler(command.Add)
	case "set-interval":
		handler.SetIntervalHandler(command)
	case "set-workers":
		fmt.Println("SetWorkers: ", command.SetWorkers)
		log.Printf("Number of workers changed from () to %s\n", command.SetWorkers.Count)
	case "list":
		fmt.Println("List: ", command.List)
	case "delete":
		fmt.Println("Delete: ", command.Delete)
	case "articles":
		fmt.Println("Articles: ", command.ArticlesCommand)
	case "fetch":
		handler.FetchHandler()
	}
}
