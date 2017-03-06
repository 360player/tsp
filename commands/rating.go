package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/360player/tsp/caller"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Rating struct {
	ID                  int          `json:"-"`
	RaterUserID         int          `json:"raterUserId"`
	UserID              int          `json:"userId"`
	GroupID             int          `json:"groupId"`
	SuggestedPositionID int          `json:"suggestedPositionId"`
	PrimaryFoot         string       `json:"primaryFoot"`
	MaturityAdjustment  int          `json:"maturityAdjustment"`
	RatingItems         []ratingItem `json:"ratingItems"`
}

type apiRatingItems struct {
	Collection
	Items []struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"records"`
}

type apiPositions struct {
	Collection
	Positions []position `json:"records"`
}

type ratingItem struct {
	ID    int    `json:"ratingItemId"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type position struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LongName string `json:"longName"`
}

var reader *bufio.Reader

// Fetch rating items from the API and loop through them all.
// A value between 1 and 100 is required for input.
// Values above 100 will be set 100.
// Values below 1 will prompt input on the same item again.
func (rating *Rating) setRatingItems() {
	ratingItemsJson, riErr := caller.Get(caller.EP_RATING_ITEMS)

	if riErr != nil {
		panic(riErr)
	}

	var itemCollection apiRatingItems
	json.Unmarshal(ratingItemsJson, &itemCollection)

	fmt.Println("Grade skills (1-100)")

	rating.RatingItems = make([]ratingItem, len(itemCollection.Items), len(itemCollection.Items))

	longestName := 0

	for _, item := range itemCollection.Items {
		if len(item.Name) > longestName {
			longestName = len(item.Name)
		}
	}

	longestName += 2

	var gradeSkill func(index int)
	gradeSkill = func(index int) {
		item := itemCollection.Items[index]

		fmt.Printf("%2d/%d %-"+strconv.Itoa(longestName)+"s: ", index+1, len(itemCollection.Items), item.Name)
		value, _ := reader.ReadString('\n')
		ratingItem := ratingItem{}

		ratingItem.Value, _ = strconv.Atoi(strings.TrimRight(value, "\n"))

		if ratingItem.Value > 100 {
			ratingItem.Value = 100
		} else if ratingItem.Value < 1 {
			fmt.Println("Try again")
			gradeSkill(index)
			return
		}

		ratingItem.ID = item.ID
		ratingItem.Name = item.Name

		rating.RatingItems[index] = ratingItem

		if index+1 < len(itemCollection.Items) {
			gradeSkill(index + 1)
		}
	}

	gradeSkill(0)
}

// Set the rating position, fetch available positions from the API.
// This method recurses until a valid position is set.
// A position is valid when it exists in the API response.
func (rating *Rating) setPosition() {
	positionsJson, err := caller.Get(caller.EP_POSITIONS)

	if err != nil {
		panic(err)
	}

	var positionCollection apiPositions
	json.Unmarshal(positionsJson, &positionCollection)

	fmt.Println("Pick a position, choose a number below")

	positionIds := make(map[int]struct{})

	for _, position := range positionCollection.Positions {
		positionIds[position.ID] = struct{}{}
		fmt.Printf("%2d. %s\n", position.ID, position.LongName)
	}

	var pickPosition func()

	pickPosition = func() {
		fmt.Print("Position: ")
		terminalInput, _ := reader.ReadString('\n')
		id, _ := strconv.Atoi(strings.TrimRight(terminalInput, "\n"))

		if id == 0 {
			fmt.Println("Invalid position, try again")
			pickPosition()
		} else if _, ok := positionIds[id]; !ok {
			fmt.Println("Invalid position, try again")
			pickPosition()
		} else {
			rating.SuggestedPositionID = id
		}
	}

	pickPosition()
}

func (rating *Rating) setUser() {
	var getUser func()

	getUser = func() {
		fmt.Print("User ID: ")
		userId, _ := reader.ReadString('\n')
		rating.UserID, _ = strconv.Atoi(strings.TrimRight(userId, "\n"))

		userJson, userError := caller.Get(fmt.Sprintf(caller.EP_USER, rating.UserID))

		if userError != nil {
			fmt.Println("That user doesn't exist, try again")
			getUser()
		} else {
			user := make(map[string]interface{})

			jsonError := json.Unmarshal(userJson, &user)

			if jsonError == nil {
				fmt.Println("Creating rating for", user["firstName"], user["lastName"])
			}
		}
	}

	getUser()
}

func (rating *Rating) setGroup() {
	var getGroup func()

	getGroup = func() {
		fmt.Print("Group ID (optional): ")
		groupId, _ := reader.ReadString('\n')

		if strings.TrimRight(groupId, "\n") == "" {
			return
		}

		rating.GroupID, _ = strconv.Atoi(strings.TrimRight(groupId, "\n"))

		groupJson, userError := caller.Get(fmt.Sprintf(caller.EP_GROUP, rating.GroupID))

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

func (rating *Rating) setMaturityAdjustment() {
	var getMaturityAdjustment func()

	options := map[int]string{
		-2: "Very early",
		-1: "Early",
		0:  "Average (default)",
		1:  "Late",
		2:  "Very late",
	}

	keys := []int{-2, -1, 0, 1, 2}
	sort.Ints(keys)

	fmt.Println("Choose a maturity level (early or late in puberty). Only applicable for players up to the age of 18):")

	for _, key := range keys {
		fmt.Printf("%2d %s\n", key, options[key])
	}

	getMaturityAdjustment = func() {
		fmt.Print("Maturity adjustment: ")
		input, _ := reader.ReadString('\n')

		maturityAdjustment, _ := strconv.Atoi(strings.TrimRight(input, "\n"))

		if maturityAdjustment < -2 || maturityAdjustment > 2 {
			fmt.Println("Invalid maturity adjustment, try again")
			getMaturityAdjustment()
		} else {
			fmt.Println("Setting maturity adjustment to", options[maturityAdjustment])
			rating.MaturityAdjustment = maturityAdjustment
		}
	}

	getMaturityAdjustment()
}

func (rating *Rating) setFoot() {
	fmt.Println("Choose primary foot")

	fmt.Println("0. Right (default)")
	fmt.Println("1. Left")

	fmt.Print("Primary foot: ")

	input, _ := reader.ReadString('\n')
	foot, _ := strconv.Atoi(strings.TrimRight(input, "\n"))

	if foot == 0 {
		rating.PrimaryFoot = "right"
	} else {
		rating.PrimaryFoot = "left"
	}

	fmt.Println("Setting primary foot to", rating.PrimaryFoot)
}

func (rating *Rating) Create() {
	fmt.Println("Creating rating")

	reader = bufio.NewReader(os.Stdin)

	fmt.Println("")
	rating.setUser()

	fmt.Println("")
	rating.setGroup()

	fmt.Println("")
	rating.setFoot()

	fmt.Println("")
	rating.setMaturityAdjustment()

	fmt.Println("")
	rating.setPosition()

	fmt.Println("")
	rating.setRatingItems()

	_, _ = caller.Post(fmt.Sprintf(caller.EP_USER_RATINGS, rating.UserID), rating)
}
