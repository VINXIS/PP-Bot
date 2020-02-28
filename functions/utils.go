package functions

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"../structs"
	"../values"
	"github.com/bwmarrin/discordgo"
)

// Build builds osu-tools
func Build() error {
	delta := exec.Command("dotnet", "build", "./osu-tools/delta/osu-tools/PerformanceCalculator", "-c", "Release")
	joz := exec.Command("dotnet", "build", "./osu-tools/joz/osu-tools/PerformanceCalculator", "-c", "Release")
	live := exec.Command("dotnet", "build", "./osu-tools/live/osu-tools/PerformanceCalculator", "-c", "Release")
	_, err := delta.Output()
	if err != nil {
		delta.Process.Kill()
		return err
	}
	delta.Process.Kill()
	_, err = joz.Output()
	if err != nil {
		joz.Process.Kill()
		return err
	}
	joz.Process.Kill()
	_, err = live.Output()
	if err != nil {
		joz.Process.Kill()
		return err
	}
	live.Process.Kill()
	return nil
}

// accGeneration returns 300s, 100s, and misses basde
func accGeneration(target float64, objects, misses int) (greats, goods, mehs int) {
	greats = objects - misses

	target /= 100 // decimals instead of percentages

	newacc := float64(6*greats+2*goods+mehs) / float64(6*objects)

	// LOOP TIME
	if newacc > target {
		for {
			greats--
			goods++
			newacc = float64(6*greats+2*goods+mehs) / float64(6*objects)
			if newacc < target {
				goods--
				mehs++
				newacc = float64(6*greats+2*goods+mehs) / float64(6*objects)
				if newacc < target {
					mehs--
					greats++
					if goods == 0 {
						break
					} else {
						for {
							if goods == 0 {
								break
							}
							goods--
							mehs++
							newacc = float64(6*greats+2*goods+mehs) / float64(6*objects)
							if newacc < target {
								goods++
								mehs--
								break
							}
						}
					}

					break
				}
			}
		}
	}

	return greats, goods, mehs
}

// removeFiles removes extra files created by diff/pp calc
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
	os.Remove(mapInfo + ".png")
	os.Remove(mapInfo + "ppVals.txt")

	// Remove blanks if they still exist:
	if mapInfo != "" {
		removeFiles("", osuType)
	}
}

// sendPaste sends data to paste.ee
func sendPaste(s *discordgo.Session, m *discordgo.MessageCreate, pasteData structs.PasteData, text string) {
	// Send paste data
	req, _ := http.NewRequest("POST", "https://api.paste.ee/v1/pastes?key="+values.Conf.PasteAPIKey, bytes.NewBuffer(pasteData.Marshal()))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "No response found from paste.ee.")
		return
	}
	defer resp.Body.Close()

	// Parse result
	bod, _ := ioutil.ReadAll(resp.Body)
	pasteResult := structs.PasteResult{}
	json.Unmarshal(bod, &pasteResult)
	if !pasteResult.Success {
		s.ChannelMessageSend(m.ChannelID, "An error occurred in sending the user calc to paste.ee!")
		log.Println(string(bod))
		return
	}

	s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"\n<"+pasteResult.Link+">\n"+text)
}
