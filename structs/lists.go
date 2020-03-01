package structs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
)

// List holds a list of maps for the user
type List struct {
	User  discordgo.User
	Lists []SubList
}

// SubList holds the different sections for the user's full list
type SubList struct {
	Name   string
	Scores []Score
}

// Score holds the data for the score
type Score struct {
	MapInfo     string
	BeatmapID   int
	Accuracy    float64
	Goods       int
	Mehs        int
	Combo       int
	Misses      int
	Mods        string
	UseAccuracy bool
}

// PPScore is an extension of Score where the PP is also stored
type PPScore struct {
	Score
	PP float64
}

// SRScore is an extension of Score where the SR is also stored
type SRScore struct {
	Score
	SR     float64
	Aim    float64
	Tap    float64
	Finger float64
}

// GetList gets either a new list for the user, or the list saved for the user, the bool tells if it is new or not.
func GetList(u *discordgo.User) (List, bool) {
	newList := List{
		User: *u,
		Lists: []SubList{
			SubList{
				Name: "Untitled",
			},
		},
	}
	list := newList

	_, err := os.Stat("./lists/" + u.ID + ".json")
	if err != nil {
		return newList, true
	}

	f, err := ioutil.ReadFile("./lists/" + u.ID + ".json")
	if err != nil {
		return newList, true
	}

	err = json.Unmarshal(f, &list)
	if err != nil {
		return newList, true
	}

	return list, false
}

// LegacyScore is the structure of the old list format's score
type LegacyScore struct {
	MapInfo   string  `json:"mapinfo"`
	BeatmapID int     `json:"beatmapid"`
	Accuracy  float64 `json:"accuracy"`
	Combo     int     `json:"combo"`
	Misses    int     `json:"misses"`
	Mods      string  `json:"mods"`
}
