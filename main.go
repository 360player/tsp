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

func auth() {

}

func printHelp() {
	fmt.Println("Tsp is an interface for the 360Player API.")
}

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

		userResponse, userError := caller.Get(caller.EP_USER_LIST, caller.QueryParam{"page", *page})

		if userError != nil {
			panic(userError)
		}

		userList := &commands.UserList{}

		userList.Unmarshal(userResponse)
		fmt.Println("Total users:", userList.RecordCount)
		fmt.Println("Showing page", userList.Page, "out of", userList.PageCount)

		for _, user := range userList.Users {
			fmt.Println(user.ID, user.FirstName, user.LastName)
		}
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
	case "config":
		handleConfig()
		break
	case "users":
		handleUsers()
		break
	}
}
