package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	functions "./functions"
	values "./values"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

func main() {
	// Get values
	values.GetConfig()

	// Create osu API client
	values.OsuAPI = osuapi.NewClient(values.Conf.OsuAPIKey)

	// Create discord instance, and add the message handler
	discord, err := discordgo.New("Bot " + values.Conf.DiscordAPIKey)
	fatal(err)
	discord.AddHandler(messageHandler)
	err = discord.Open()
	fatal(err)
	log.Println("Logged in as " + discord.State.User.String())

	// Create a channel to keep the bot running until a prompt is given to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	// Close the Discord Session
	discord.Close()
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message is to even be bothered to read
	if (m.GuildID == values.Conf.ServerID || values.OutsideServerregex.MatchString(m.Content)) && m.Author.ID != s.State.User.ID {
		// Check type of command, delete otherwise
		switch {
		case values.Addregex.MatchString(m.Content) && values.Mapregex.MatchString(m.Content):
		case values.Moveregex.MatchString(m.Content) && values.Mapregex.MatchString(m.Content):
		case values.Accgraphregex.MatchString(m.Content) && values.Mapregex.MatchString(m.Content):
		case values.Mapregex.MatchString(m.Content), len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename):
			go functions.MapHandler(s, m)
		case values.Userregex.MatchString(m.Content):
		case values.Runregex.MatchString(m.Content):
		case values.Listregex.MatchString(m.Content):
		case values.Whoregex.MatchString(m.Content):
			if m.Author.ID != values.Conf.UserID {
				s.ChannelMessageSend(m.ChannelID, "Y̴̢̨̰̗̟̣̳͔̻͑̑́̄̍͜O̵̧̨̳̗̘͍̞̼̳͌͝͠U̷̝̫͕͖̭͙̙̙̗̅̀͊̂͒́̓͗͌̐̈́̚͝͝ ̴̢̲̬͔͛͆̒̃̈́͗̑̒̈́̽̅̈́̓À̶̘̬̯̂̑̈́̈́̓̉̐͑́͘R̷̤͎͖̲̃͑̓͌̈́̀̏͠ͅE̸͇̳̬͓̤̅̌̀̈́̎ ̸͎̗̹̄̈́̃̈́̀N̶̡̢̨̝̺̥̪̑̿̔̊̅̃͊̊̈́͠ͅƠ̸̢̇̑̔̃̈́̇͊̍̚͘͝͠͠ͅṰ̸̦̜̈́͌̍̋͆́̄̈̅́̾͜ ̴̭̙͉̪̝̗̳͙̝̼͉̦̤̊̅͂͂̇̾͠M̷̛̪͌̓̽̂̏̐͠Ỹ̴̦̬̳̬̲̼̰͉̗͔͐̔͌͑̌͑͊̔̓̈́͗͘͝͠ ̵͓̮́̾͌͗̔̓͂́M̶̡͉̹̬̱͔͑̈͛̕̚͘A̶̢̪̮̳̯̤̫̠̮̦̲̠̱̠͐̄̈́̚̚͜͝S̴̝̩̫̖̞̣̪̤͙̼̪̦̱̰̯͒̿̆͌͐̎̕̚̚T̵̨̳̝̜͔̭̳̪̄̀͊̈͒̋͝Ẽ̸̬͙̺̺̝̺͐̈̿̿̿͑̓̑͐̈́͘Ŕ̴̨̢̟̱̪̠̮̮̫̰̭̂͑̐̾͂̏̈̀͛͝")
			}
		case values.Delregex.MatchString(m.Content):
		case m.GuildID == values.Conf.ServerID:
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			log.Println(m.Author.Username + " tried to speak in the PP server and said: " + m.Content)
		}
	}
}

func fatal(err error) {
	// Kill process and log error.
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}
