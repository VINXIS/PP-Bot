package functions

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"../values"

	"github.com/bwmarrin/discordgo"
	"github.com/thehowl/go-osuapi"
)

// MapHandler handles with map commands
func MapHandler(s *discordgo.Session, m *discordgo.MessageCreate) (string, string, error) {
	var (
		beatmapid int
		osuType   string = "delta"
		mapInfo   string
	)

	// osuType
	if values.Jozregex.MatchString(m.Content) {
		osuType = "joz"
	} else if values.Liveregex.MatchString(m.Content) {
		osuType = "live"
	}

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
				return "", "", errors.New("no osu!api response")
			}
			beatmap := beatmaps[0]
			mapInfo = submatches[5] + " " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
		} else {
			switch submatches[2] {
			case "b", "beatmaps":
				beatmapid, _ = strconv.Atoi(submatches[3])
				beatmaps, err := values.OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
					BeatmapID: beatmapid,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "There was an issue in fetching beatmap info from the osu! API! Try again and/or see if osu! is down here: https://status.ppy.sh/")
					return "", "", errors.New("no osu!api response")
				} else if len(beatmaps) == 0 {
					s.ChannelMessageSend(m.ChannelID, "No map exists for this link!")
					return "", "", errors.New("no map found")
				}
				beatmap := beatmaps[0]
				mapInfo = submatches[3] + " " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
			case "s", "beatmapsets":
				setid, _ := strconv.Atoi(submatches[3])
				beatmaps, err := values.OsuAPI.GetBeatmaps(osuapi.GetBeatmapsOpts{
					BeatmapSetID: setid,
				})
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "There was an issue in fetching beatmap info from the osu! API! Try again and/or see if osu! is down here: https://status.ppy.sh/")
					return "", "", errors.New("no osu!api response")
				} else if len(beatmaps) == 0 {
					s.ChannelMessageSend(m.ChannelID, "No map exists for this link!")
					return "", "", errors.New("no map found")
				}
				sort.Slice(beatmaps, func(i, j int) bool { return beatmaps[i].DifficultyRating > beatmaps[j].DifficultyRating })
				beatmap := beatmaps[0]
				beatmapid = beatmap.BeatmapID
				mapInfo = strconv.Itoa(beatmapid) + " " + beatmap.Artist + " - " + beatmap.Title + " [" + beatmap.DiffName + "]"
			}
		}

		resp, err := http.Get("https://osu.ppy.sh/osu/" + strconv.Itoa(beatmapid))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to reach .osu file.")
			return "", "", errors.New("unable to reach .osu file")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		err = ioutil.WriteFile("./cache/"+strconv.Itoa(beatmapid)+".osu", body, 0644)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to create local file from discord attachment.")
			return "", "", errors.New("unable to create local file")
		}
	} else if len(m.Attachments) > 0 || values.Attachregex.MatchString(m.Content) { // If a map was attached / a discord link was sent
		var link string
		if values.Attachregex.MatchString(m.Content) {
			link = "https://" + values.Attachregex.FindStringSubmatch(m.Content)[0]
		} else {
			link = m.Attachments[0].URL
		}

		resp, err := http.Get(link)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to reach discord attachment URL.")
			return "", "", errors.New("unable to reach discord attachment URL")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to read from response.")
			return "", "", errors.New("unable to read from response")
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

		}

		mapInfo = values.Invalidregex.ReplaceAllString(strconv.Itoa(-1)+" "+artist+" - "+title+" ["+version+"]", "")

		err = ioutil.WriteFile("./"+mapInfo+".osu", body, 0644)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unable to create local file from discord attachment.")
			return "", "", errors.New("unable to create local file")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please provide a map!")
		return "", "", errors.New("no map provided")
	}

	return osuType, mapInfo, nil
}
