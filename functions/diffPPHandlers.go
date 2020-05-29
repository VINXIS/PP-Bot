package functions

import (
	"bytes"
	"io/ioutil"
	"log"
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
	log.Println(m.Author.String() + " has requested a difficulty calc for " + mapInfo + ".")
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	defer removeFiles(mapInfo, osuType)
	mapID := strings.Split(mapInfo, " ")[0]

	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp3.1/PerformanceCalculator.dll", "difficulty", "./" + mapInfo + ".osu"}

	// Get score specs (acc, combo, e.t.c)
	var mods string
	if values.Modregex.MatchString(m.Content) {
		mods = values.Modregex.FindStringSubmatch(m.Content)[1]
		for i := 0; i < len(mods); i += 2 {
			args = append(args, "-m", string(mods[i])+string(mods[i+1]))
		}
	}

	switch osuType {
	case "joz":
		args[0] = "./osu-tools/joz/osu-tools/PerformanceCalculator/bin/Release/netcoreapp3.1/PerformanceCalculator.dll"
	case "live":
		args[0] = "./osu-tools/live/osu-tools/PerformanceCalculator/bin/Release/netcoreapp3.1/PerformanceCalculator.dll"
	}

	if !strings.HasPrefix(mapInfo, "-1") {
		args[2] = "./cache/" + strings.Split(mapInfo, " ")[0] + ".osu"
	}

	process := exec.Command("dotnet", args...)
	res, err := process.Output()
	if err != nil || string(res) == values.InvalidCommand {
		process.Process.Kill()
		s.ChannelMessageSend(m.ChannelID, "Could not run command!")
		return
	}
	process.Process.Kill()

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
		if values.TapGraphregex.MatchString(m.Content) {
			graphContent, err = ioutil.ReadFile("./cache/graph_" + mapID + "_" + mods + "_tap.txt")
			if err != nil {
				graphContent, err = ioutil.ReadFile("./cache/graph__" + mods + "_tap.txt")
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Could not find tap graph data!")
					return
				}
			}
		} else if values.Fingerregex.MatchString(m.Content) {
			graphContent, err = ioutil.ReadFile("./cache/graph_" + mapID + "_" + mods + "_finger.txt")
			if err != nil {
				graphContent, err = ioutil.ReadFile("./cache/graph__" + mods + "_finger.txt")
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Could not find finger control graph data!")
					return
				}
			}
		} else {
			graphContent, err = ioutil.ReadFile("./cache/graph_" + mapID + "_" + mods + ".txt")
			if err != nil {
				graphContent, err = ioutil.ReadFile("./cache/graph__" + mods + ".txt")
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Could not find graph data!")
					return
				}
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

	args = []string{"plot.py", skill, mapID, strconv.Itoa(start), strconv.Itoa(end), strconv.Itoa(difference), mapInfo, mods, "delta"}
	if osuType == "joz" {
		args[len(args)-1] = "joz"
	} else if values.TapGraphregex.MatchString(m.Content) {
		args[len(args)-1] = "tap"
	} else if values.Fingerregex.MatchString(m.Content) {
		args[len(args)-1] = "finger"
	}

	// Generate graph using python script
	_, err = exec.Command("python", args...).Output()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "An error in generating a graph for the map occurred!")
		return
	}

	// Send value and delete files
	img, err := ioutil.ReadFile("./" + mapID + ".png")
	imgBytes := bytes.NewBuffer(img)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: m.Author.Mention() + "\n```\n" + strings.ToValidUTF8(string(res), "") + "```",
		Files: []*discordgo.File{
			{
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
	log.Println(m.Author.String() + " has requested a pp calc for " + mapInfo + ".")
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	defer removeFiles(mapInfo, osuType)

	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp3.1/PerformanceCalculator.dll", "simulate", "osu", "./" + mapInfo + ".osu"}

	// Get score specs (acc, combo, e.t.c)
	if values.Modregex.MatchString(m.Content) {
		mods := values.Modregex.FindStringSubmatch(m.Content)[1]
		for i := 0; i < len(mods); i += 2 {
			args = append(args, "-m", string(mods[i])+string(mods[i+1]))
		}
	}

	if !values.Accregex.MatchString(m.Content) {
		if values.Goodregex.MatchString(m.Content) {
			args = append(args, "-G", values.Goodregex.FindStringSubmatch(m.Content)[1])
		}
		if values.Mehregex.MatchString(m.Content) {
			args = append(args, "-M", values.Mehregex.FindStringSubmatch(m.Content)[1])
		}
	} else {
		args = append(args, "-a", values.Accregex.FindStringSubmatch(m.Content)[1])
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
		args[0] = "./osu-tools/joz/osu-tools/PerformanceCalculator/bin/Release/netcoreapp3.1/PerformanceCalculator.dll"
	case "live":
		args[0] = "./osu-tools/live/osu-tools/PerformanceCalculator/bin/Release/netcoreapp3.1/PerformanceCalculator.dll"
	}

	if !strings.HasPrefix(mapInfo, "-1") {
		args[3] = "./cache/" + strings.Split(mapInfo, " ")[0] + ".osu"
	}

	// Run command
	process := exec.Command("dotnet", args...)
	res, err := process.Output()
	if err != nil || string(res) == values.InvalidCommand {
		process.Process.Kill()
		s.ChannelMessageSend(m.ChannelID, "Could not run command!")
		return
	}
	process.Process.Kill()

	// Send pp value and remove files
	s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"\n```\n"+string(res)+"```")
}
