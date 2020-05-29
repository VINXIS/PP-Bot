package values

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Conf holds the keys for the APIs
var Conf Config

// Config is the struct for configuration
type Config struct {
	DiscordAPIKey     string
	RoleLogChannel    string
	MessageLogChannel string
	JoinLogChannel    string
	CalcChannels      []string
	WhitelistRoles    []string
	OnionRoles        []string
	OsuAPIKey         string
	PasteAPIKey       string
	UserID            string
	ServerID          string
}

// GetConfig gets the config information
func GetConfig() {
	config := Config{}
	f, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("Error obtaining config information: " + err.Error())
	}
	_ = json.Unmarshal(f, &config)
	Conf = config
}
