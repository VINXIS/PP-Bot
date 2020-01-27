package functions

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"../values"
	"github.com/bwmarrin/discordgo"
)

// AccGraphHandler obtains the acc graph of a map
func AccGraphHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	osuType, mapInfo, err := MapHandler(s, m)
	if err != nil {
		return
	}
	defer removeFiles(mapInfo, osuType)

	// Create args foundation
	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll", "simulate", "osu", "./" + mapInfo + ".osu"}

	// Reject other osu! versions
	if osuType == "joz" || osuType == "live" {
		s.ChannelMessageSend(m.ChannelID, "No support for joz or live versions regarding acc graphs yet!")
		return
	}

	if !strings.HasPrefix(mapInfo, "-1") {
		args[3] = "./cache/" + strings.Split(mapInfo, " ")[0] + ".osu"
	}

	var (
		accrange   []float64 = []float64{95, 100}
		difference float64   = 5.0
		increment  float64   = 0.5
	)

	// Get score specs (mods, combo, e.t.c) and fill in above variables
	if values.Modregex.MatchString(m.Content) {
		mods := values.Modregex.FindStringSubmatch(m.Content)[1]
		for i := 0; i < len(mods); i += 2 {
			args = append(args, "-m", string(mods[i])+string(mods[i+1]))
		}
	}

	if values.Comboregex.MatchString(m.Content) {
		args = append(args, "--combo", values.Comboregex.FindStringSubmatch(m.Content)[1])
	}

	if values.Missregex.MatchString(m.Content) {
		args = append(args, "-X", values.Missregex.FindStringSubmatch(m.Content)[1])
	}

	if values.Accrangeregex.MatchString(m.Content) {
		acc1, _ := strconv.ParseFloat(values.Accrangeregex.FindStringSubmatch(m.Content)[1], 64)
		acc2, _ := strconv.ParseFloat(values.Accrangeregex.FindStringSubmatch(m.Content)[2], 64)
		accrange[0] = math.Min(acc1, 100)
		accrange[1] = math.Min(acc2, 100)
		difference = accrange[1] - accrange[0]
		increment = difference / 10.0
	}

	if values.Incrementregex.MatchString(m.Content) {
		increment, _ = strconv.ParseFloat(values.Incrementregex.FindStringSubmatch(m.Content)[1], 64)
		if increment > difference {
			increment = difference / 10.0
		}
	}

	if difference <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a valid acc range with a difference greater than 0.")
		return
	}

	// Begin pp value accumulation
	msg, err := s.ChannelMessageSend(m.ChannelID, "Running accuracies from "+strconv.FormatFloat(accrange[0], 'f', 2, 64)+" to "+strconv.FormatFloat(accrange[1], 'f', 2, 64)+" in increments of "+strconv.FormatFloat(increment, 'f', 2, 64)+" for "+mapInfo)
	if err != nil {
		return
	}
	log.Println(m.Author.String() + " has requested an acc graph for " + mapInfo + ".")
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)

	var accValues string
	for i := accrange[0]; i <= accrange[1]; i += increment {
		tempArgs := append(args, "-a", strconv.FormatFloat(i, 'f', 2, 64))
		process := exec.Command("dotnet", tempArgs...)
		res, err := process.Output()
		if err != nil {
			process.Process.Kill()
			s.ChannelMessageSend(m.ChannelID, "Error in obtaining pp values for the accuracy: "+strconv.FormatFloat(i, 'f', 2, 64)+".")
			return
		}
		process.Process.Kill()
		txt := strings.Replace(string(res), "Accuracy", "", 1)
		aim, _ := strconv.ParseFloat(values.Aimregex.FindStringSubmatch(txt)[1], 64)
		tap, _ := strconv.ParseFloat(values.Tapregex.FindStringSubmatch(txt)[1], 64)
		acc, _ := strconv.ParseFloat(values.AccPPregex.FindStringSubmatch(txt)[1], 64)
		pp, _ := strconv.ParseFloat(values.PPparseregex.FindStringSubmatch(txt)[1], 64)
		accValues += "(" +
			strconv.FormatFloat(i, 'f', 2, 64) + ", " +
			strconv.FormatFloat(aim, 'f', 2, 64) + ", " +
			strconv.FormatFloat(tap, 'f', 2, 64) + ", " +
			strconv.FormatFloat(acc, 'f', 2, 64) + ", " +
			strconv.FormatFloat(pp, 'f', 2, 64) + ")\n"
	}
	err = ioutil.WriteFile(mapInfo+"ppVals.txt", []byte(accValues), 0644)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in saving pp values.")
		return
	}

	// Get acc graph
	_, err = exec.Command("python", "acc.py",
		mapInfo,
		strconv.FormatFloat(accrange[0], 'f', 2, 64),
		strconv.FormatFloat(accrange[1], 'f', 2, 64),
		strconv.FormatFloat(difference, 'f', 2, 64),
		strconv.FormatFloat(increment, 'f', 2, 64)).Output()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in generating the acc graph for the map.")
		return
	}

	// Send value and delete files
	img, err := ioutil.ReadFile("./" + mapInfo + ".png")
	imgBytes := bytes.NewBuffer(img)
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: m.Author.Mention() + "\n" + mapInfo,
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   mapInfo + ".png",
				Reader: imgBytes,
			},
		},
	})
}
