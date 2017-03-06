package main

import (
	"flag"
	"fmt"
	"github.com/360player/tsp/caller"
	"github.com/360player/tsp/commands"
	"github.com/360player/tsp/config"
	"os"
)

var appConfig *config.Config

func printHelp() {
	fmt.Println("Tsp is an interface for the 360Player API.")
}

// Reconfigure the app, setting a new API URL also requires the user to log back in.
func handleConfig() {
	appConfig.SetApiUrl()
	caller.SetBaseUrl(appConfig.ApiUrl)

	appConfig.Auth()
	caller.SetAuth(appConfig.ApiKey)
}

func handleUsers() {
	if len(os.Args) < 3 || os.Args[2] == "help" {
		//@todo Print user help
		return
	}

	listFlagSet := flag.NewFlagSet("user list", flag.ContinueOnError)
	var page = listFlagSet.String("page", "1", "")

	switch os.Args[2] {
	case "list":
		listFlagSet.Parse(os.Args[3:])

		userList := &commands.UserList{}
		userList.List(*page)

		fmt.Println("Total users:", userList.RecordCount)
		fmt.Println("Showing page", userList.Page, "out of", userList.PageCount)

		for _, user := range userList.Users {
			fmt.Println(user.ID, user.FirstName, user.LastName)
		}
	}
}

func handleRatings() {
	if len(os.Args) < 3 || os.Args[2] == "help" {
		//@todo Print ratings help
		return
	}

	switch os.Args[2] {
	case "create":
		rating := &commands.Rating{}
		rating.Create()
		break
	}
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "help" {
		printHelp()
		return
	}

	appConfig = &config.Config{}
	appConfig.Load()

	if appConfig.ApiUrl == "" || appConfig.ApiKey == "" {
		handleConfig()
	}

	caller.SetBaseUrl(appConfig.ApiUrl)
	caller.SetAuth(appConfig.ApiKey)

	switch os.Args[1] {
	case "login":
		appConfig.Auth()
		break
	case "config":
		handleConfig()
		break
	case "users":
		handleUsers()
		break
	case "ratings":
		handleRatings()
		break
	}
}
