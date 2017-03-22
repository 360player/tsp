package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/360player/tsp/caller"
	"os"
	"strconv"
	"strings"
)

type Event struct {
	GroupID     int    `json:"groupId"`
	StartsAt    int64  `json:"startsAt"`
	EndsAt      int64  `json:"endsAt"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (event *Event) setGroup() {
	var getGroup func()

	getGroup = func() {
		fmt.Print("Group ID: ")
		groupId, _ := reader.ReadString('\n')

		if strings.TrimRight(groupId, "\n") == "" {
			fmt.Println("Group is mandatory, try again")
			getGroup()
			return
		}

		event.GroupID, _ = strconv.Atoi(strings.TrimRight(groupId, "\n"))

		groupJson, userError := caller.Get(fmt.Sprintf(caller.EP_GROUP, event.GroupID))

		if userError != nil {
			fmt.Println("That group doesn't exist, try again")
			getGroup()
		} else {
			group := make(map[string]interface{})

			jsonError := json.Unmarshal(groupJson, &group)

			if jsonError == nil {
				fmt.Println("In", group["name"])
			}
		}
	}

	getGroup()
}

func (event *Event) setStartTime() {
	fmt.Println("Set start time, UTC timezone")

	var setTime func()

	setTime = func() {
		startDate := inputDate()
		event.StartsAt = startDate.Unix()
	}

	setTime()
}

func (event *Event) setEndTime() {
	fmt.Print("Does the event have an end time? (y/n): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")

	if input != "y" {
		event.EndsAt = event.StartsAt
		return
	}

	fmt.Println("Set end time, UTC timezone")

	var setTime func()

	setTime = func() {
		endDate := inputDate()
		event.EndsAt = endDate.Unix()

		if event.EndsAt < event.StartsAt {
			fmt.Println("\nEnd date has to be after start date")
			setTime()
		}
	}

	setTime()
}

func (event *Event) setTitle() {
	fmt.Print("Event title: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	title := input

	if title == "" {
		title = "No title"
	}

	event.Title = title
}

func (event *Event) setDescription() {
	fmt.Print("Event description: ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	description := input

	event.Description = description
}

func (event *Event) Create() {
	fmt.Println("Create event")

	reader = bufio.NewReader(os.Stdin)

	fmt.Println("")
	event.setGroup()

	fmt.Println("")
	event.setStartTime()

	fmt.Println("")
	event.setEndTime()

	fmt.Println("")
	event.setTitle()

	fmt.Println("")
	event.setDescription()

	_, _ = caller.Post(fmt.Sprintf("/v1/events"), event)
}
