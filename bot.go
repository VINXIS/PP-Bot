package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"./functions"
	"./values"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

var pings = make(map[string]time.Time)

func main() {
	// Get args
	build := false
	channelLog := false
	quiet := false
	for _, arg := range os.Args {
		if arg == "-b" || arg == "--build" {
			build = true
		} else if arg == "-l" || arg == "--log" {
			channelLog = true
		} else if arg == "-q" || arg == "--quiet" {
			quiet = true
		}
	}

	// Get values
	values.GetConfig()

	if !channelLog {
		// Create cache folder
		_, err := os.Stat("./cache")
		if os.IsNotExist(err) {
			err = os.MkdirAll("./cache", 0755)
		}

		// Create lists folder
		_, err = os.Stat("./lists")
		if os.IsNotExist(err) {
			err = os.MkdirAll("./lists", 0755)
			fatal(err)
		}

		// Change console type for proper output
		_, err = exec.Command("chcp", "65001").Output()
		fatal(err)
	} else if values.Conf.RoleLogChannel == "" || values.Conf.MessageLogChannel == "" || values.Conf.JoinLogChannel == "" {
		fatal(errors.New("Please add log channel IDs"))
	}

	if build {
		// Build osu-tools
		log.Println("Building osu-tools...")
		err := functions.Build()
		fatal(err)
		log.Println("Built osu-tools!")
	}

	// Create osu API client
	values.OsuAPI = osuapi.NewClient(values.Conf.OsuAPIKey)

	// Create discord instance, and add the message handler
	discord, err := discordgo.New("Bot " + values.Conf.DiscordAPIKey)
	fatal(err)
	if channelLog {
		discord.AddHandler(logMessageHandler)
		discord.AddHandler(logMessageEditHandler)
		discord.AddHandler(roleHandler)
		discord.AddHandler(joinHandler)
		discord.AddHandler(leaveHandler)
		log.Println("Added logging!")
	} else {
		discord.AddHandler(normalMessageHandler)
	}
	err = discord.Open()
	for err != nil {
		err = discord.Open()
	}
	log.Println("Logged in as " + discord.State.User.String())
	if !channelLog && !quiet {
		for _, ch := range values.Conf.CalcChannels {
			go discord.ChannelMessageSend(ch, "osu! calculations are now up. ("+strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1)+")")
		}
	}

	// Create a channel to keep the bot running until a prompt is given to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc

	if !channelLog && !quiet {
		for _, ch := range values.Conf.CalcChannels {
			go discord.ChannelMessageSend(ch, "osu! calculations are now going down. ("+strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1)+")")
		}
	}

	// Close the Discord Session
	discord.Close()
}

func normalMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message is to even be bothered to read
	if (m.GuildID == values.Conf.ServerID || values.OutsideServerregex.MatchString(m.Content)) && m.Author.ID != s.State.User.ID {
		// Check if it's a command
		switch {
		case values.Helpregex.MatchString(m.Content):
			go functions.HelpHandler(s, m)
		case values.Addregex.MatchString(m.Content) && values.Mapregex.MatchString(m.Content): // Add score to list
			go functions.ListAddHandler(s, m)
		case values.Accgraphregex.MatchString(m.Content) && (values.Mapregex.MatchString(m.Content) || len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename)): // Get accuracy graph for a map
			go functions.AccGraphHandler(s, m)
		case values.Attachregex.MatchString(m.Content), values.Mapregex.MatchString(m.Content), len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename): // Get map SR / PP
			if values.PPregex.MatchString(m.Content) {
				go functions.MapPPHandler(s, m)
			} else {
				go functions.MapDifficultyHandler(s, m)
			}
		case values.Userregex.MatchString(m.Content): // Run user profile
			go functions.UserHandler(s, m)
		case values.Runregex.MatchString(m.Content): // Run user list
			go functions.ListRunHandler(s, m)
		case values.Listregex.MatchString(m.Content): // Show list
			go functions.ListHandler(s, m)
		case values.Moveregex.MatchString(m.Content): // Move score between lists
			go functions.ListMoveHandler(s, m)
		case values.Delregex.MatchString(m.Content): // Delete map from list
			go functions.ListDeleteHandler(s, m)
		case values.Whoregex.MatchString(m.Content): // See user IDs and who has what list
			if m.Author.ID != values.Conf.UserID {
				s.ChannelMessageSend(m.ChannelID, "Y̴̢̨̰̗̟̣̳͔̻͑̑́̄̍͜O̵̧̨̳̗̘͍̞̼̳͌͝͠U̷̝̫͕͖̭͙̙̙̗̅̀͊̂͒́̓͗͌̐̈́̚͝͝ ̴̢̲̬͔͛͆̒̃̈́͗̑̒̈́̽̅̈́̓À̶̘̬̯̂̑̈́̈́̓̉̐͑́͘R̷̤͎͖̲̃͑̓͌̈́̀̏͠ͅE̸͇̳̬͓̤̅̌̀̈́̎ ̸͎̗̹̄̈́̃̈́̀N̶̡̢̨̝̺̥̪̑̿̔̊̅̃͊̊̈́͠ͅƠ̸̢̇̑̔̃̈́̇͊̍̚͘͝͠͠ͅṰ̸̦̜̈́͌̍̋͆́̄̈̅́̾͜ ̴̭̙͉̪̝̗̳͙̝̼͉̦̤̊̅͂͂̇̾͠M̷̛̪͌̓̽̂̏̐͠Ỹ̴̦̬̳̬̲̼̰͉̗͔͐̔͌͑̌͑͊̔̓̈́͗͘͝͠ ̵͓̮́̾͌͗̔̓͂́M̶̡͉̹̬̱͔͑̈͛̕̚͘A̶̢̪̮̳̯̤̫̠̮̦̲̠̱̠͐̄̈́̚̚͜͝S̴̝̩̫̖̞̣̪̤͙̼̪̦̱̰̯͒̿̆͌͐̎̕̚̚T̵̨̳̝̜͔̭̳̪̄̀͊̈͒̋͝Ẽ̸̬͙̺̺̝̺͐̈̿̿̿͑̓̑͐̈́͘Ŕ̴̨̢̟̱̪̠̮̮̫̰̭̂͑̐̾͂̏̈̀͛͝")
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				return
			}
			go functions.ListWhoHandler(s, m)
		case values.Buildregex.MatchString(m.Content): // Build osu-tools
			if m.Author.ID != values.Conf.UserID {
				s.ChannelMessageSend(m.ChannelID, "Y̴̢̨̰̗̟̣̳͔̻͑̑́̄̍͜O̵̧̨̳̗̘͍̞̼̳͌͝͠U̷̝̫͕͖̭͙̙̙̗̅̀͊̂͒́̓͗͌̐̈́̚͝͝ ̴̢̲̬͔͛͆̒̃̈́͗̑̒̈́̽̅̈́̓À̶̘̬̯̂̑̈́̈́̓̉̐͑́͘R̷̤͎͖̲̃͑̓͌̈́̀̏͠ͅE̸͇̳̬͓̤̅̌̀̈́̎ ̸͎̗̹̄̈́̃̈́̀N̶̡̢̨̝̺̥̪̑̿̔̊̅̃͊̊̈́͠ͅƠ̸̢̇̑̔̃̈́̇͊̍̚͘͝͠͠ͅṰ̸̦̜̈́͌̍̋͆́̄̈̅́̾͜ ̴̭̙͉̪̝̗̳͙̝̼͉̦̤̊̅͂͂̇̾͠M̷̛̪͌̓̽̂̏̐͠Ỹ̴̦̬̳̬̲̼̰͉̗͔͐̔͌͑̌͑͊̔̓̈́͗͘͝͠ ̵͓̮́̾͌͗̔̓͂́M̶̡͉̹̬̱͔͑̈͛̕̚͘A̶̢̪̮̳̯̤̫̠̮̦̲̠̱̠͐̄̈́̚̚͜͝S̴̝̩̫̖̞̣̪̤͙̼̪̦̱̰̯͒̿̆͌͐̎̕̚̚T̵̨̳̝̜͔̭̳̪̄̀͊̈͒̋͝Ẽ̸̬͙̺̺̝̺͐̈̿̿̿͑̓̑͐̈́͘Ŕ̴̨̢̟̱̪̠̮̮̫̰̭̂͑̐̾͂̏̈̀͛͝")
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				return
			}
			go functions.BuildHandler(s, m)
		case values.Importregex.MatchString(m.Content) && len(m.Attachments) > 0 && strings.HasSuffix(m.Attachments[0].Filename, ".json"): // Import a map list
			go functions.ListImportHandler(s, m)
		}
	}
}

func logMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if message in calc channel
	found := false
	for _, ch := range values.Conf.CalcChannels {
		if ch == m.ChannelID {
			found = true
			break
		}
	}

	// Check if the message is to even be bothered to read
	if found && (m.GuildID == values.Conf.ServerID || values.OutsideServerregex.MatchString(m.Content)) && m.Author.ID != s.State.User.ID {

		// Check if user has whitelist role
		for _, role := range m.Member.Roles {
			for _, whitelist := range values.Conf.WhitelistRoles {
				if whitelist == role {
					return
				}
			}
		}

		// Check for pings
		if len(m.Mentions) > 0 || len(m.MentionRoles) > 0 {
			go s.ChannelMessageDelete(m.ChannelID, m.ID)
			lastPing, exist := pings[m.Author.ID]
			if !exist || time.Now().Sub(lastPing).Hours() > 24 {
				go s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Do not ping others, or else you will be **banned.**")
				pings[m.Author.ID] = time.Now()
			} else if time.Now().Sub(lastPing).Hours() < 24 {
				go s.GuildBanCreateWithReason(m.GuildID, m.Author.ID, "Pinging twice within 24 hours.", -1)
				go delete(pings, m.Author.ID)
				if values.Conf.MessageLogChannel != "" {
					ch, err := s.Channel(m.ChannelID)
					if err == nil {
						go s.ChannelMessageSend(values.Conf.MessageLogChannel, "**"+m.Author.String()+"** tried to speak in **#"+ch.Name+"** has now been **banned.** ("+strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1)+")")
					}
				} else {
					ch, err := s.Channel(m.ChannelID)
					if err == nil {
						go log.Println("**" + m.Author.String() + "** tried to speak in **#" + ch.Name + "** and has now been **banned.** (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
					} else {
						go log.Println("**" + m.Author.String() + "** tried to speak and has now been **banned.** (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
					}
				}
			}
		}

		help := values.Helpregex.MatchString(m.Content)
		add := values.Addregex.MatchString(m.Content) && values.Mapregex.MatchString(m.Content)
		acc := values.Accgraphregex.MatchString(m.Content) && (values.Mapregex.MatchString(m.Content) || len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename))
		beatmap := values.Attachregex.MatchString(m.Content) || values.Mapregex.MatchString(m.Content) || (len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename))
		user := values.Userregex.MatchString(m.Content)
		run := values.Runregex.MatchString(m.Content)
		list := values.Listregex.MatchString(m.Content)
		move := values.Moveregex.MatchString(m.Content)
		delete := values.Delregex.MatchString(m.Content)
		who := values.Whoregex.MatchString(m.Content)
		build := values.Buildregex.MatchString(m.Content)
		listImport := values.Importregex.MatchString(m.Content) && len(m.Attachments) > 0 && strings.HasSuffix(m.Attachments[0].Filename, ".json")
		inServer := m.GuildID == values.Conf.ServerID

		// Delete messages that are not commands
		if inServer && !help && !add && !acc && !beatmap && !user && !run && !list && !move && !delete && !who && !build && !listImport {
			go s.ChannelMessageDelete(m.ChannelID, m.ID)
			if values.Conf.MessageLogChannel != "" {
				ch, err := s.Channel(m.ChannelID)
				if err == nil {
					go s.ChannelMessageSend(values.Conf.MessageLogChannel, "**"+m.Author.String()+"** tried to speak in **#"+ch.Name+"** and said: "+m.Content+" ("+strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1)+")")
				}
			} else {
				ch, err := s.Channel(m.ChannelID)
				if err == nil {
					go log.Println("**" + m.Author.String() + "** tried to speak in **#" + ch.Name + "** and said: " + m.Content + " (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
				} else {
					go log.Println("**" + m.Author.String() + "** tried to speak and said: " + m.Content + " (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
				}
			}
		}
	}
}

// for edits, idk how to kill this message duplication
func logMessageEditHandler(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author == nil { // Embed loading
		return
	}

	// Check if message in calc channel
	found := false
	for _, ch := range values.Conf.CalcChannels {
		if ch == m.ChannelID {
			found = true
			break
		}
	}

	// Check if the message is to even be bothered to read
	if found && (m.GuildID == values.Conf.ServerID || values.OutsideServerregex.MatchString(m.Content)) && m.Author.ID != s.State.User.ID {

		// Check if user has whitelist role
		for _, role := range m.Member.Roles {
			for _, whitelist := range values.Conf.WhitelistRoles {
				if whitelist == role {
					return
				}
			}
		}

		// Check for pings
		if len(m.Mentions) > 0 || len(m.MentionRoles) > 0 {
			go s.ChannelMessageDelete(m.ChannelID, m.ID)
			lastPing, exist := pings[m.Author.ID]
			if !exist || time.Now().Sub(lastPing).Hours() > 24 {
				go s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Do not ping others, or else you will be **banned.**")
				pings[m.Author.ID] = time.Now()
			} else if time.Now().Sub(lastPing).Hours() < 24 {
				go s.GuildBanCreateWithReason(m.GuildID, m.Author.ID, "Pinging twice within 24 hours.", -1)
				go delete(pings, m.Author.ID)
				if values.Conf.MessageLogChannel != "" {
					ch, err := s.Channel(m.ChannelID)
					if err == nil {
						go s.ChannelMessageSend(values.Conf.MessageLogChannel, "**"+m.Author.String()+"** tried to speak in **#"+ch.Name+"** has now been **banned.** ("+strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1)+")")
					}
				} else {
					ch, err := s.Channel(m.ChannelID)
					if err == nil {
						go log.Println("**" + m.Author.String() + "** tried to speak in **#" + ch.Name + "** and has now been **banned.** (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
					} else {
						go log.Println("**" + m.Author.String() + "** tried to speak and has now been **banned.** (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
					}
				}
			}
		}

		help := values.Helpregex.MatchString(m.Content)
		add := values.Addregex.MatchString(m.Content) && values.Mapregex.MatchString(m.Content)
		acc := values.Accgraphregex.MatchString(m.Content) && (values.Mapregex.MatchString(m.Content) || len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename))
		beatmap := values.Attachregex.MatchString(m.Content) || values.Mapregex.MatchString(m.Content) || (len(m.Attachments) > 0 && values.Fileregex.MatchString(m.Attachments[0].Filename))
		user := values.Userregex.MatchString(m.Content)
		run := values.Runregex.MatchString(m.Content)
		list := values.Listregex.MatchString(m.Content)
		move := values.Moveregex.MatchString(m.Content)
		delete := values.Delregex.MatchString(m.Content)
		who := values.Whoregex.MatchString(m.Content)
		build := values.Buildregex.MatchString(m.Content)
		listImport := values.Importregex.MatchString(m.Content) && len(m.Attachments) > 0 && strings.HasSuffix(m.Attachments[0].Filename, ".json")
		inServer := m.GuildID == values.Conf.ServerID

		// Delete messages that are not commands
		if inServer && !help && !add && !acc && !beatmap && !user && !run && !list && !move && !delete && !who && !build && !listImport {
			go s.ChannelMessageDelete(m.ChannelID, m.ID)
			if values.Conf.MessageLogChannel != "" {
				ch, err := s.Channel(m.ChannelID)
				if err == nil {
					go s.ChannelMessageSend(values.Conf.MessageLogChannel, "**"+m.Author.String()+"** tried to speak in **#"+ch.Name+"** and said: "+m.Content+" ("+strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1)+")")
				}
			} else {
				ch, err := s.Channel(m.ChannelID)
				if err == nil {
					go log.Println("**" + m.Author.String() + "** tried to speak in **#" + ch.Name + "** and said: " + m.Content + " (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
				} else {
					go log.Println("**" + m.Author.String() + "** tried to speak and said: " + m.Content + " (" + strings.Replace(time.Now().UTC().Format(time.RFC822Z), "+0000", "UTC", -1) + ")")
				}
			}
		}
	}
}

func roleHandler(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	// Make sure if it's in the proper server
	if values.Conf.ServerID != m.GuildID {
		return
	}

	auditLog, err := s.GuildAuditLog(values.Conf.ServerID, values.Conf.UserID, "", 25, -1)
	if err != nil {
		s.ChannelMessageSend(values.Conf.RoleLogChannel, "A role was changed, but there was an error in obtaining the audit log!")
		return
	}

	roleLog := auditLog.AuditLogEntries[0]

	// Check the time cuz this just spams the latest rn
	timeStamp, err := discordgo.SnowflakeTimestamp(roleLog.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if time.Now().Sub(timeStamp).Minutes() > 1 {
		return
	}

	roleAffectedUser, err := s.User(roleLog.TargetID)
	if err != nil {
		s.ChannelMessageSend(values.Conf.RoleLogChannel, "A role was changed, but there was an error in obtaining the user!")
		return
	}
	roleAction := roleLog.Changes[0].Key
	roleName := roleLog.Changes[0].NewValue.([]interface{})[0].(map[string]interface{})["name"].(string)
	if *roleAction == discordgo.AuditLogChangeKeyRoleAdd {
		s.ChannelMessageSend(values.Conf.RoleLogChannel, "The user **"+roleAffectedUser.String()+"** has been given the **"+roleName+"** role!")
	} else if *roleAction == discordgo.AuditLogChangeKeyRoleRemove {
		s.ChannelMessageSend(values.Conf.RoleLogChannel, "The user **"+roleAffectedUser.String()+"** has lost the **"+roleName+"** role!")
	}
}

func joinHandler(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Make sure if it's in the proper server
	if values.Conf.ServerID != m.GuildID {
		return
	}

	s.ChannelMessageSend(values.Conf.JoinLogChannel, "**"+m.User.String()+"** has joined!")
}

func leaveHandler(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	// Make sure if it's in the proper server
	if values.Conf.ServerID != m.GuildID {
		return
	}

	s.ChannelMessageSend(values.Conf.JoinLogChannel, "**"+m.User.String()+"** has left!")
}

func fatal(err error) {
	// Kill process and log error.
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}
