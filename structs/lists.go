package structs

import (
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
