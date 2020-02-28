package functions

import (
	"github.com/bwmarrin/discordgo"
)

// BuildHandler builds osu-tools
func BuildHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := Build()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "An error occurred in building!")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Built osu-tools!")
}
