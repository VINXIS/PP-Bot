package functions

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	values "../values"
	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// MapHandler handles with map commands
func MapHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		beatmapid int
		skill     string = "aimcontrol"
		osuType   string = "joz"
		mapinfo   string
	)

	// See if a specific skill was requested

	// Get beatmap info if a map was linked/attached
	if values.Mapregex.MatchString(m.Content) { // If a map was linked
		submatches := values.Mapregex.FindStringSubmatch(m.Content)
		if submatches[5] != "" {
			beatmapid, _ = strconv.Atoi(submatches[5])
			beatmaps, err := values.OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
				BeatmapID: beatmapid,
			})
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "There was an issue in fetching beatmap info from the osu! API! Try again and/or see if osu! is down here: https://status.ppy.sh/")
				return
			}
			beatmap := beatmaps[0]
			mapinfo = submatches[5] + ": " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
		} else {
			switch submatches[2] {
			case "b", "beatmaps":
				beatmapid, _ = strconv.Atoi(submatches[3])
				beatmaps, err := values.OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
					BeatmapID: beatmapid,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "There was an issue in fetching beatmap info from the osu! API! Try again and/or see if osu! is down here: https://status.ppy.sh/")
					return
				}
				beatmap := beatmaps[0]
				mapinfo = submatches[3] + ": " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
			case "s":
				setid, _ := strconv.Atoi(submatches[3])
				beatmaps, err := values.OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
					BeatmapSetID: setid,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "There was an issue in fetching beatmap info from the osu! API! Try again and/or see if osu! is down here: https://status.ppy.sh/")
					return
				}
				sort.Slice(beatmaps, func(i, j int) bool { return beatmaps[i].DifficultyRating > beatmaps[j].DifficultyRating })
				beatmap := beatmaps[0]
				beatmapid = beatmap.BeatmapID
				mapinfo = strconv.Itoa(beatmapid) + ": " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
			}
		}

		resp, err := http.Get("https://osu.ppy.sh/osu/" + strconv.Itoa(beatmapid))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to reach .osu file.")
			resp.Body.Close()
			return
		}
		out, err := os.Create("./cache/" + strconv.Itoa(beatmapid) + ".osu")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to create .osu file.")
			resp.Body.Close()
			return
		}
		io.Copy(out, resp.Body)
		resp.Body.Close()
		out.Close()
	} else { // If a map was attached
		resp, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to reach discord attachment URL.")
			resp.Body.Close()
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to read from response.")
			resp.Body.Close()
			return
		}
		stringbody := strings.Split(string(body), "\n")
		var artist, title, version string = "", "", ""
		for _, line := range stringbody {
			if values.Titleregex.MatchString(line) {
				title = values.Titleregex.FindStringSubmatch(line)[1]
			}
			if values.Artistregex.MatchString(line) {
				artist = values.Artistregex.FindStringSubmatch(line)[1]
			}
			if values.Versionregex.MatchString(line) {
				version = values.Versionregex.FindStringSubmatch(line)[1]
			}
			if values.BeatmapIDregex.MatchString(line) {
				beatmapid, _ = strconv.Atoi(values.BeatmapIDregex.FindStringSubmatch(line)[1])
			}

		}
		if artist != "" && title != "" && version != "" {
			mapinfo = strconv.Itoa(beatmapid) + ": " + artist + " - " + title + " [" + version + "]"
		}

		file, err := os.Create(m.Attachments[0].Filename)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable create local file from discord attachment.")
			resp.Body.Close()
			file.Close()
			return
		}

		resp.Body.Close()
		file.Close()
	}
	fmt.Println(skill)
	fmt.Println(osuType)
	fmt.Println(mapinfo)

}

// MapDifficultyHandler handles with the difficulty graph of a map
func MapDifficultyHandler() {

}

// MapPPHandler handles with the pp of a map
func MapPPHandler() {

}
