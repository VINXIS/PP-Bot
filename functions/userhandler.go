package functions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
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
	res, err := exec.Command("dotnet", args...).Output()
	if err != nil || res[0] == 83 {
		s.ChannelMessageSend(m.ChannelID, "Could not run command!")
		return
	}

	// Send paste data
	pasteData := structs.NewPasteData(user.Username, string(res))
	ioutil.WriteFile("test.json", pasteData.Marshal(), 0644)
	req, _ := http.NewRequest("POST", "https://api.paste.ee/v1/pastes?key="+values.Conf.PasteAPIKey, bytes.NewBuffer(pasteData.Marshal()))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "No response found from paste.ee.")
	}
	defer resp.Body.Close()

	// Parse result
	bod, _ := ioutil.ReadAll(resp.Body)
	pasteResult := structs.PasteResult{}
	json.Unmarshal(bod, &pasteResult)
	if !pasteResult.Success {
		s.ChannelMessageSend(m.ChannelID, "An error occurred in sending the user calc to paste.ee!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"\n<"+pasteResult.Link+">\nUser calc for **"+user.Username+"** done in "+time.Now().Sub(timeTaken).String())
}
