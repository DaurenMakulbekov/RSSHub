package handlers

import (
	"RSSHub/internal/core/domain"
	"fmt"
	"strings"
	"time"
)

func GetAdd(args []string) (domain.Commands, error) {
	var command = domain.Commands{
		Name: "add",
		Add:  domain.Add{},
	}

	if len(args) == 1 {
		return command, fmt.Errorf("Required options")
	}

	args = args[1:]
	if len(args) < 4 {
		return command, fmt.Errorf("Incorrect input")
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "--name" && len(args) > index+1 {
			command.Add.Name = args[index+1]
		} else if args[index] == "--url" && len(args) > index+1 {
			command.Add.Url = args[index+1]
		} else {
			return command, fmt.Errorf("Incorrect input")
		}
		index++
	}

	return command, nil
}

func GetFetch(args string) domain.Commands {
	var command = domain.Commands{
		Name: "fetch",
	}

	return command
}

func GetSetInterval(args []string) (domain.Commands, error) {
	var command = domain.Commands{
		Name:        "set-interval",
		SetInterval: domain.SetInterval{},
	}

	if len(args) == 1 {
		return command, fmt.Errorf("Required options")
	}

	args = args[1:]
	if len(args) < 2 {
		return command, fmt.Errorf("Incorrect input")
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "--duration" && len(args) > index+1 {
			s, err := time.ParseDuration(args[index+1])

			if err != nil || strings.Contains(args[index+1], "-") {
				return command, fmt.Errorf("Incorrect input")
			}

			command.SetInterval.Duration = time.Duration(s.Seconds())
		} else {
			return command, fmt.Errorf("Incorrect input")
		}
		index++
	}

	return command, nil
}

func GetSetWorkers(args []string) (domain.Commands, error) {
	var command = domain.Commands{
		Name:       "set-workers",
		SetWorkers: domain.SetWorkers{},
	}

	if len(args) == 1 {
		return command, fmt.Errorf("Required options")
	}

	args = args[1:]
	if len(args) < 2 {
		return command, fmt.Errorf("Incorrect input")
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "--count" && len(args) > index+1 {
			command.SetWorkers.Count = args[index+1]
		} else {
			return command, fmt.Errorf("Incorrect input")
		}
		index++
	}

	return command, nil
}

func GetList(args []string) (domain.Commands, error) {
	var command = domain.Commands{
		Name: "list",
		List: domain.List{},
	}

	if len(args) == 1 {
		return command, nil
	}

	args = args[1:]
	if len(args) < 2 {
		return command, fmt.Errorf("Incorrect input")
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "--num" && len(args) > index+1 {
			command.List.Num = args[index+1]
		} else {
			return command, fmt.Errorf("Incorrect input")
		}
		index++
	}

	return command, nil
}

func GetDelete(args []string) (domain.Commands, error) {
	var command = domain.Commands{
		Name:   "delete",
		Delete: domain.Delete{},
	}

	if len(args) == 1 {
		return command, fmt.Errorf("Required options")
	}

	args = args[1:]
	if len(args) < 2 {
		return command, fmt.Errorf("Incorrect input")
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "--name" && len(args) > index+1 {
			command.Delete.Name = args[index+1]
		} else {
			return command, fmt.Errorf("Incorrect input")
		}
		index++
	}

	return command, nil
}

func GetArticles(args []string) (domain.Commands, error) {
	var command = domain.Commands{
		Name:            "articles",
		ArticlesCommand: domain.ArticlesCommand{},
	}

	if len(args) == 1 {
		return command, fmt.Errorf("Required options")
	}

	args = args[1:]
	if len(args) < 2 {
		return command, fmt.Errorf("Incorrect input")
	}

	for index := 0; index < len(args); index++ {
		if args[index] == "--feed-name" && len(args) > index+1 {
			command.ArticlesCommand.FeedName = args[index+1]
		} else if args[index] == "--num" && len(args) > index+1 {
			command.ArticlesCommand.Num = args[index+1]
		} else {
			return command, fmt.Errorf("Incorrect input")
		}
		index++
	}

	return command, nil
}

func (hd *handler) GetCommand(args []string) (domain.Commands, error) {
	var command domain.Commands
	var err error

	for i := range args {
		if args[i] == "add" {
			command, err = GetAdd(args[i:])
			return command, err
		} else if args[i] == "fetch" {
			command = GetFetch(args[i])
			return command, nil
		} else if args[i] == "set-interval" {
			command, err = GetSetInterval(args[i:])
			return command, err
		} else if args[i] == "set-workers" {
			command, err = GetSetWorkers(args[i:])
			return command, err
		} else if args[i] == "list" {
			command, err = GetList(args[i:])
			return command, err
		} else if args[i] == "delete" {
			command, err = GetDelete(args[i:])
			return command, err
		} else if args[i] == "articles" {
			command, err = GetArticles(args[i:])
			return command, err
		} else {
			return command, fmt.Errorf(args[i] + ": command not found...")
		}
	}

	return command, nil
}
