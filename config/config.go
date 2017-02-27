package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/360player/tsp/caller"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
)

var reader *bufio.Reader

const CONFIG_FILE_NAME = "~/.tsp.json"

type Config struct {
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
}

type userAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:token`
}

func (config *Config) Load() {
	rawConfig, configErr := ioutil.ReadFile(configFilePath())

	if configErr == nil {
		json.Unmarshal(rawConfig, &config)
	}
}

func (config *Config) Auth() {
	if config.ApiUrl == "" {
		fmt.Println("Config.ApiUrl needs to be set to log in")
		config.SetApiUrl()
	}

	fmt.Println("Enter your 360Player credentials.")

	reader = bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimRight(username, "\n")

	fmt.Print("Password: ")
	password, _ := terminal.ReadPassword(0)

	resp, respErr := caller.Post(caller.EP_AUTH, userAuth{
		Username: username,
		Password: string(password[:]),
	})

	if respErr != nil {
		panic(respErr)
	}

	var authData *authResponse
	json.Unmarshal(resp, &authData)

	config.ApiKey = authData.Token
	config.Write()
}

func (config *Config) SetApiUrl() {
	fmt.Print("360Player API URL: ")

	reader = bufio.NewReader(os.Stdin)
	url, _ := reader.ReadString('\n')

	config.ApiUrl = strings.TrimRight(url, "\n")
	config.Write()
}

func (config *Config) Write() {
	configJson, configError := json.Marshal(&config)

	if configError != nil {
		panic(configError)
	}

	writeError := ioutil.WriteFile(configFilePath(), configJson, 0644)

	if writeError != nil {
		panic(writeError)
	}
}

// Expand CONFIG_FILE_NAME, ~ is normally expanded by the shell.
// Built in os/user only work on Darwin systems, so we use go-homedir.
func configFilePath() string {
	absPath, err := homedir.Expand(CONFIG_FILE_NAME)

	if err != nil {
		panic(err)
	}

	return absPath
}
