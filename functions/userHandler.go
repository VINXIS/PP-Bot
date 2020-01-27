package functions

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"../structs"
	"../values"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// UserHandler handles with user commands
func UserHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	osuType := "delta"
	if values.Jozregex.MatchString(m.Content) {
		osuType = "joz"
	} else if values.Liveregex.MatchString(m.Content) {
		osuType = "live"
	}

	// Get user
	user := new(osuapi.User)
	username := values.Userregex.FindStringSubmatch(m.Content)[3]
	userID, err := strconv.Atoi(username)
	if err == nil {
		user, err = values.OsuAPI.GetUser(osuapi.GetUserOpts{
			UserID: userID,
		})
		if err != nil {
			user, _ = values.OsuAPI.GetUser(osuapi.GetUserOpts{
				Username: username,
			})
		}
	} else {
		user, _ = values.OsuAPI.GetUser(osuapi.GetUserOpts{
			Username: username,
		})
	}

	// Check if user was found
	if user == nil || user.UserID == 0 {
		s.ChannelMessageSend(m.ChannelID, "No user found")
		return
	}
	msg, err := s.ChannelMessageSend(m.ChannelID, "Obtaining profile calc for "+user.Username+"...")
	if err != nil {
		return
	}
	log.Println(m.Author.String() + " has requested a profile calc for " + user.Username + ".")
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	timeTaken := time.Now()

	// Create args
	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll", "profile", strconv.Itoa(user.UserID), values.Conf.OsuAPIKey}

	// Get score specs (acc, combo, e.t.c)
	if values.Modregex.MatchString(m.Content) {
		mods := values.Modregex.FindStringSubmatch(m.Content)[1]
		for i := 0; i < len(mods); i += 2 {
			args = append(args, "-m", string(mods[i])+string(mods[i+1]))
		}
	}

	switch osuType {
	case "joz":
		args[0] = "./osu-tools/joz/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	case "live":
		args[0] = "./osu-tools/live/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	}

	// Run command
	process := exec.Command("dotnet", args...)
	res, err := process.Output()
	if err != nil || res[0] == 83 {
		process.Process.Kill()
		s.ChannelMessageSend(m.ChannelID, "Could not run command!")
		return
	}
	process.Process.Kill()

	go sendPaste(s, m, structs.NewPasteData(user.Username, string(res)), "User calc for **"+user.Username+"** done in "+time.Now().Sub(timeTaken).String())

	if osuType == "joz" {
		// Remove filename spam
		var fileNames []string

		err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fileNames = append(fileNames, path)
			return nil
		})

		for _, fileName := range fileNames {
			if values.Spamfileregex.MatchString(fileName) {
				os.Remove(fileName)
			}
		}
	}
}
