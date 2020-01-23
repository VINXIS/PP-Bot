package functions

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"../values"

	"github.com/bwmarrin/discordgo"
)

// MapDifficultyHandler handles with the difficulty graph of a map
func MapDifficultyHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	osuType, mapInfo, err := MapHandler(s, m)
	if err != nil {
		return
	}
	msg, err := s.ChannelMessageSend(m.ChannelID, "Obtaining diff calc for "+mapInfo+"...")
	if err != nil {
		return
	}
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	defer removeFiles(mapInfo, osuType)
	mapID := strings.Split(mapInfo, " ")[0]

	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll", "difficulty", "./" + mapInfo + ".osu"}

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

	if !strings.HasPrefix(mapInfo, "-1") {
		args[2] = "./cache/" + strings.Split(mapInfo, " ")[0] + ".osu"
	}

	res, err := exec.Command("dotnet", args...).Output()
	if err != nil || res[0] == 83 {
		s.ChannelMessageSend(m.ChannelID, "Could not run command!")
		return
	}

	// Finish live diff calc here
	if osuType == "live" {
		s.ChannelMessageSend(m.ChannelID, "```\n"+strings.ToValidUTF8(string(res), "")+"```")
		return
	}

	// Get skill choice if applicable
	skill := "aimcontrol"
	if values.Skillregex.MatchString(m.Content) {
		skill = strings.ReplaceAll(strings.ToLower(values.Skillregex.FindStringSubmatch(m.Content)[1]), " ", "")
	}

	// Get time range
	times := []int{}
	if values.Timeregex.MatchString(m.Content) {
		time1, _ := strconv.Atoi(values.Timeregex.FindStringSubmatch(m.Content)[1])
		time2, _ := strconv.Atoi(values.Timeregex.FindStringSubmatch(m.Content)[2])
		times = []int{time1, time2}
	}

	// Get graph content
	var graphContent []byte

	switch osuType {
	case "delta":
		graphContent, err = ioutil.ReadFile("./cache/graph_" + mapID + ".txt")
		if err != nil {
			graphContent, err = ioutil.ReadFile("./cache/graph_.txt")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Could not find graph data!")
				return
			}
		}
	case "joz":
		graphContent, err = ioutil.ReadFile("./" + mapID + skill + ".txt")
		if err != nil {
			graphContent, err = ioutil.ReadFile("./" + skill + ".txt")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Could not find graph data!")
				return
			}
		}
		res, err = ioutil.ReadFile("./" + mapID + "values.txt")
		if err != nil {
			res, err = ioutil.ReadFile("./values.txt")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Could not get values!")
				return
			}
		}
	}
	lines := strings.Split(string(graphContent), "\n")

	// Get graph axes configuration
	var (
		x                      []int
		difference, start, end int
	)
	for _, line := range lines {
		line = strings.Replace(strings.Replace(line, "(", "", -1), ")", "", -1)
		val, _ := strconv.ParseFloat(strings.Split(line, " ")[0], 64)
		if osuType == "joz" {
			val /= 1000.0
		}
		x = append(x, int(val))
	}
	if len(times) < 2 { // Time values weren't given in message
		times = x
		start = times[0]
		end = times[len(times)-2]
	} else { // Time values were given in message
		// Check for values going out of bounds
		if times[0] < x[0] {
			times[0] = x[0]
		}
		if times[1] > x[len(x)-2] {
			times[1] = x[len(x)-2]
		}

		start = times[0]
		end = times[1]
	}
	difference = end - start

	args = []string{"plot.py", skill, mapID, strconv.Itoa(start), strconv.Itoa(end), strconv.Itoa(difference), "\"" + mapInfo + "\"", "delta"}
	if osuType == "joz" {
		args[len(args)-1] = "joz"
	}

	// Generate graph using python script
	_, err = exec.Command("python", args...).Output()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error generating graph for the map!")
		return
	}

	// Send value and delete files
	img, err := ioutil.ReadFile("./" + mapID + ".png")
	imgBytes := bytes.NewBuffer(img)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "```\n" + strings.ToValidUTF8(string(res), "") + "```",
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   mapID + ".png",
				Reader: imgBytes,
			},
		},
	})
}

// MapPPHandler handles with the pp of a map
func MapPPHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	osuType, mapInfo, err := MapHandler(s, m)
	if err != nil {
		return
	}
	msg, err := s.ChannelMessageSend(m.ChannelID, "Obtaining pp calc for "+mapInfo+"...")
	if err != nil {
		return
	}
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	defer removeFiles(mapInfo, osuType)

	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll", "simulate", "osu", mapInfo + ".osu"}

	// Get score specs (acc, combo, e.t.c)
	if values.Modregex.MatchString(m.Content) {
		mods := values.Modregex.FindStringSubmatch(m.Content)[1]
		for i := 0; i < len(mods); i += 2 {
			args = append(args, "-m", string(mods[i])+string(mods[i+1]))
		}
	}

	if values.Accregex.MatchString(m.Content) {
		args = append(args, "-a", values.Accregex.FindStringSubmatch(m.Content)[1])
	} else {
		if values.Goodregex.MatchString(m.Content) {
			args = append(args, "-G", values.Goodregex.FindStringSubmatch(m.Content)[1])
		}
		if values.Mehregex.MatchString(m.Content) {
			args = append(args, "-M", values.Mehregex.FindStringSubmatch(m.Content)[1])
		}
	}

	if values.Comboregex.MatchString(m.Content) {
		args = append(args, "--combo", values.Comboregex.FindStringSubmatch(m.Content)[1])
	}

	if values.Missregex.MatchString(m.Content) {
		args = append(args, "-X", values.Missregex.FindStringSubmatch(m.Content)[1])
	}

	// Change osu! version
	switch osuType {
	case "joz":
		args[0] = "./osu-tools/joz/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	case "live":
		args[0] = "./osu-tools/live/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	}

	if !strings.HasPrefix(mapInfo, "-1") {
		args[3] = "./cache/" + strings.Split(mapInfo, " ")[0] + ".osu"
	}

	// Run command
	res, err := exec.Command("dotnet", args...).Output()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not run command!")
		return
	}

	// Send pp value and remove files
	s.ChannelMessageSend(m.ChannelID, "```\n"+string(res)+"```")
}

func removeFiles(mapInfo, osuType string) {
	if strings.HasPrefix(mapInfo, "-1") { // Non-submitted beatmap
		os.Remove(mapInfo + ".osu")
	}

	mapID := strings.Split(mapInfo, " ")[0]
	if osuType == "joz" {
		os.Remove(mapID + "aimcontrol.txt")
		os.Remove(mapID + "fingercontrol.txt")
		os.Remove(mapID + "jumpaim.txt")
		os.Remove(mapID + "speed.txt")
		os.Remove(mapID + "stamina.txt")
		os.Remove(mapID + "streamaim.txt")
		os.Remove(mapID + "test.txt")
		os.Remove(mapID + "values.txt")
	}
	os.Remove(mapID + ".png")

	// Remove blanks if they still exist:
	if mapInfo != "" {
		removeFiles("", osuType)
	}
}
