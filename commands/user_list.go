package commands

import (
	"encoding/json"
	"github.com/360player/tsp/caller"
)

type User struct {
	ID        int           `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	Meta      []interface{} `json:"meta"`
	UpdatedAt int           `json:"updatedAt"`
	CreatedAt int           `json:"createdAt"`
}

type UserList struct {
	Collection
	Users []User `json:"records"`
}

func (userList *UserList) List(page string) {
	userResponse, userError := caller.Get(caller.EP_USER_LIST, caller.QueryParam{"page", page})

	if userError != nil {
		panic(userError)
	}

	userList.Unmarshal(userResponse)
}

func (userList *UserList) Unmarshal(jsonData []byte) {
	err := json.Unmarshal(jsonData, &userList)

	if err != nil {
		panic(err)
	}
}
