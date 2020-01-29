package functions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"../structs"
	"../values"

	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
)

// ListHandler gives the user theirs or someone else's list
func ListHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check for mention
	user := m.Author
	if len(m.Mentions) >= 1 {
		user = m.Mentions[0]
	}

	// Get file
	list, new := structs.GetList(user)

	if new {
		if m.Author.ID == user.ID {
			s.ChannelMessageSend(m.ChannelID, "Could not find a list with your user ID!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Could not find a list for the given user!")
		}
		return
	}
	log.Println(m.Author.String() + " has requested for " + user.Username + "'s list.")

	// Create tables
	mainText := user.Username + "'s list:\n\n"

	for _, subList := range list.Lists {
		// Create sublist table
		var tableData [][]string
		var builder strings.Builder
		table := tablewriter.NewWriter(&builder)
		table.SetRowLine(true)
		table.SetAutoWrapText(false)

		table.SetHeader([]string{"Num", "Map Info", "Mods", "Combo", "Acc", "100s", "50s", "Misses"})

		// Add score info to tabledata
		for i, score := range subList.Scores {
			scoreData := []string{
				strconv.Itoa(i + 1),
				score.MapInfo,
				score.Mods,
				strconv.Itoa(score.Combo),
				"",
				"",
				"",
				"",
			}

			// Exceptions
			if score.Mods == "" {
				scoreData[2] = "NM"
			}
			if score.Combo == 0 {
				scoreData[3] = "FC"
			}
			if score.UseAccuracy {
				scoreData[4] = strconv.FormatFloat(score.Accuracy, 'f', 2, 64) + "%"
			} else {
				scoreData[5] = strconv.Itoa(score.Goods)
				scoreData[6] = strconv.Itoa(score.Mehs)
			}
			if score.Misses != 0 {
				scoreData[7] = strconv.Itoa(score.Misses)
			}

			tableData = append(tableData, scoreData)
		}

		// Add to table and write to string
		table.AppendBulk(tableData)
		table.Render()
		mainText += "List: " + subList.Name + "\n" + builder.String() + "\n\n"
	}

	go sendPaste(s, m, structs.NewPasteData(user.Username+"'s List", mainText), "When deleting, make sure to keep check of their NUM to make sure you do not accidentally delete the wrong score!")
}

// ListAddHandler lets the user add a map to their list
func ListAddHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get file and sublist name
	list, _ := structs.GetList(m.Author)
	subListName := "Untitled"
	if values.Optionregex.MatchString(m.Content) {
		subListName = values.Optionregex.FindStringSubmatch(m.Content)[1]
	}

	// Find sublist to make sure it exists
	index := -1
	for i, subList := range list.Lists {
		if strings.ToLower(subList.Name) == strings.ToLower(subListName) {
			index = i
			break
		}
	}

	osuType, mapInfo, err := MapHandler(s, m)
	if err != nil {
		return
	}
	msg, err := s.ChannelMessageSend(m.ChannelID, "Adding `"+mapInfo+"` to the **"+subListName+"** list...")
	if err != nil {
		return
	}
	log.Println(m.Author.String() + " has requested to add " + mapInfo + " to their " + subListName + " list.")
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	beatmapID, _ := strconv.Atoi(strings.Split(mapInfo, " ")[0])

	// Cmd foundation
	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll", "simulate", "osu", "./cache/" + strings.Split(mapInfo, " ")[0] + ".osu"}
	switch osuType {
	case "joz":
		args[0] = "./osu-tools/joz/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	case "live":
		args[0] = "./osu-tools/live/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	}

	// Create score object and get score specs, as well as add to cmd args
	score := structs.Score{
		MapInfo:   mapInfo,
		BeatmapID: beatmapID,
	}
	if values.Modregex.MatchString(m.Content) {
		score.Mods = values.Modregex.FindStringSubmatch(m.Content)[1]
		for i := 0; i < len(score.Mods); i += 2 {
			args = append(args, "-m", strings.ToUpper(string(score.Mods[i])+string(score.Mods[i+1])))
		}
	}

	if !values.Accregex.MatchString(m.Content) {
		if values.Goodregex.MatchString(m.Content) {
			goodVal, err := strconv.Atoi(values.Goodregex.FindStringSubmatch(m.Content)[1])
			if err == nil && goodVal > 0 {
				score.Goods = goodVal
				args = append(args, "-G", values.Goodregex.FindStringSubmatch(m.Content)[1])
			}
		}

		if values.Mehregex.MatchString(m.Content) {
			mehVal, err := strconv.Atoi(values.Mehregex.FindStringSubmatch(m.Content)[1])
			if err == nil && mehVal > 0 {
				score.Mehs = mehVal
				args = append(args, "-M", values.Mehregex.FindStringSubmatch(m.Content)[1])
			}
		}
	} else {
		accVal, err := strconv.ParseFloat(values.Accregex.FindStringSubmatch(m.Content)[1], 64)
		if err == nil && accVal > 0 && accVal < 100 {
			args = append(args, "-a", values.Accregex.FindStringSubmatch(m.Content)[1])
			score.Accuracy = accVal
			score.UseAccuracy = true
		}
	}

	if values.Comboregex.MatchString(m.Content) {
		comboVal, err := strconv.Atoi(values.Comboregex.FindStringSubmatch(m.Content)[1])
		if err == nil && comboVal > 0 {
			score.Combo = comboVal
			args = append(args, "--combo", values.Comboregex.FindStringSubmatch(m.Content)[1])
		}
	}

	if values.Missregex.MatchString(m.Content) {
		missVal, err := strconv.Atoi(values.Missregex.FindStringSubmatch(m.Content)[1])
		if err == nil && missVal > 0 {
			score.Misses = missVal
			args = append(args, "-X", values.Missregex.FindStringSubmatch(m.Content)[1])
		}
	}

	if score.Goods == 0 && score.Mehs == 0 && !score.UseAccuracy {
		score.Accuracy = 100
		score.UseAccuracy = true
	}

	// Add score to sublist
	if index != -1 {
		for _, listScore := range list.Lists[index].Scores {
			if listScore.MapInfo == score.MapInfo &&
				listScore.BeatmapID == score.BeatmapID &&
				listScore.Accuracy == score.Accuracy &&
				listScore.Goods == score.Goods &&
				listScore.Mehs == score.Mehs &&
				listScore.Combo == score.Combo &&
				listScore.Misses == score.Misses &&
				listScore.Mods == score.Mods {
				s.ChannelMessageSend(m.ChannelID, "Score already exists for the list given!")
				return
			}
		}
		list.Lists[index].Scores = append(list.Lists[index].Scores, score)
	} else {
		list.Lists = append(list.Lists, structs.SubList{
			Name:   subListName,
			Scores: []structs.Score{score},
		})
	}

	// Run command
	process := exec.Command("dotnet", args...)
	res, err := process.Output()
	if err != nil {
		process.Process.Kill()
		s.ChannelMessageSend(m.ChannelID, "An error occured in processing the score. Please make sure your args are correct.")
		return
	}
	process.Process.Kill()

	// Save updated list
	b, err := json.Marshal(list)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error parsing list.")
		return
	}
	err = ioutil.WriteFile("./lists/"+m.Author.ID+".json", b, 0644)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error saving list to file.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Score saved! Here's the pp value of the given score based on the pp system given:\n```"+string(res)+"```")
}

// ListMoveHandler lets the move scores between lists
func ListMoveHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !values.Numregex.MatchString(m.Content) || !values.Newregex.MatchString(m.Content) {
		s.ChannelMessageSend(m.ChannelID, "Please provide index number and target list!")
		return
	}

	// Get file
	list, new := structs.GetList(m.Author)

	if new {
		s.ChannelMessageSend(m.ChannelID, "Could not find a list with your user ID!")
		return
	}

	log.Println(m.Author.String() + " has requested to move a map.")

	var (
		err                      error
		num                      int
		oldListName, newListName string
	)
	if values.Optionregex.MatchString(m.Content) {
		oldListName = values.Optionregex.FindStringSubmatch(m.Content)[1]
	}

	if values.Numregex.MatchString(m.Content) {
		oldListName = strings.Replace(oldListName, values.Numregex.FindStringSubmatch(m.Content)[0], "", -1)
		num, _ = strconv.Atoi(values.Numregex.FindStringSubmatch(m.Content)[1])
		num--
	}

	if values.Newregex.MatchString(m.Content) {
		oldListName = strings.Replace(oldListName, values.Newregex.FindStringSubmatch(m.Content)[0], "", -1)
		newListName = values.Newregex.FindStringSubmatch(m.Content)[1]
	}

	if oldListName == "" {
		oldListName = "Untitled"
	} else {
		oldListName = strings.TrimSpace(oldListName)
	}

	var (
		oldIndex, newIndex int = -1, -1
	)

	for i, subList := range list.Lists {
		if strings.ToLower(subList.Name) == strings.ToLower(oldListName) {
			oldIndex = i
			if num > len(subList.Scores)-1 {
				s.ChannelMessageSend(m.ChannelID, "Value provided is out of bounds!")
				return
			}
		} else if strings.ToLower(subList.Name) == strings.ToLower(newListName) {
			newIndex = i
		}
	}

	if oldIndex == -1 { // No list found to move the score from
		s.ChannelMessageSend(m.ChannelID, "No list with the given name found!")
		return
	}

	// Move score
	if newIndex == -1 { // Create new list
		list.Lists = append(list.Lists, structs.SubList{
			Name: newListName,
			Scores: []structs.Score{
				list.Lists[oldIndex].Scores[num],
			},
		})
	} else {
		score := list.Lists[oldIndex].Scores[num]
		for _, listScore := range list.Lists[newIndex].Scores { // Check if score already exists there
			if listScore.MapInfo == score.MapInfo &&
				listScore.BeatmapID == score.BeatmapID &&
				listScore.Accuracy == score.Accuracy &&
				listScore.Goods == score.Goods &&
				listScore.Mehs == score.Mehs &&
				listScore.Combo == score.Combo &&
				listScore.Misses == score.Misses &&
				listScore.Mods == score.Mods {
				s.ChannelMessageSend(m.ChannelID, "Score already exists for the list given!")
				return
			}
		}
		list.Lists[newIndex].Scores = append(list.Lists[newIndex].Scores, score)
	}

	// Remove from old list, delete list if that was the only score there
	if len(list.Lists[oldIndex].Scores) == 1 {
		list.Lists = append(list.Lists[:oldIndex], list.Lists[oldIndex+1:]...)
	} else {
		list.Lists[oldIndex].Scores = append(list.Lists[oldIndex].Scores[:num], list.Lists[oldIndex].Scores[num+1:]...)
	}

	// Save updated list
	b, err := json.Marshal(list)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error parsing list.")
		return
	}
	err = ioutil.WriteFile("./lists/"+m.Author.ID+".json", b, 0644)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error saving list to file.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Score moved from `"+oldListName+"` to `"+newListName+"`!")
}

// ListDeleteHandler lets the user delete a map from a sublist
func ListDeleteHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get list
	list, new := structs.GetList(m.Author)
	if new {
		s.ChannelMessageSend(m.ChannelID, "No list found for you!")
		return
	}

	val := values.Delregex.FindStringSubmatch(m.Content)[1]
	if val == "-all" { // Delete a sublist or the whole list
		log.Println(m.Author.String() + " has requested to delete a list/sublist.")
		if values.Optionregex.MatchString(m.Content) { // If they want to only delete a sublist
			subListName := values.Optionregex.FindStringSubmatch(m.Content)[1]
			for i, subList := range list.Lists {
				if strings.ToLower(subList.Name) == strings.ToLower(subListName) {
					list.Lists = append(list.Lists[:i], list.Lists[i+1:]...)
					break
				}
			}
			// Save updated list
			b, err := json.Marshal(list)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error parsing list.")
				return
			}
			err = ioutil.WriteFile("./lists/"+m.Author.ID+".json", b, 0644)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error saving list to file.")
				return
			}

			s.ChannelMessageSend(m.ChannelID, "Removed the "+subListName+" list!")
			return

		}

		// If they want to delete the whole list
		f, err := os.Open("./lists/" + m.Author.ID + ".json")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error in opening your list file!")
			return
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: "Since you are deleting your lists completely, here is a copy of your file so you can reimport it back whenever you want, by attaching it to a message stating `!import`.",
			Files: []*discordgo.File{
				&discordgo.File{
					Name:   m.Author.ID + ".json",
					Reader: f,
				},
			},
		})
		f.Close()

		err = os.Remove("./lists/" + m.Author.ID + ".json")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error in removing your list file!")
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Completely removed "+m.Author.Username+"'s lists!")
		return
	}
	log.Println(m.Author.String() + " has requested to delete a map.")
	// If they only want to delete a score
	index, _ := strconv.Atoi(val)
	index--
	subListName := "Untitled"
	if values.Optionregex.MatchString(m.Content) {
		subListName = values.Optionregex.FindStringSubmatch(m.Content)[1]
	}

	// Find subList
	for i, subList := range list.Lists {
		if strings.ToLower(subList.Name) == strings.ToLower(subListName) {

			if index > len(list.Lists[i].Scores)-1 {
				s.ChannelMessageSend(m.ChannelID, "Number out of bounds!")
				return
			} else if index == 0 && len(list.Lists[i].Scores) == 1 {
				list.Lists = append(list.Lists[:i], list.Lists[i+1:]...)
			} else if index == len(list.Lists[i].Scores)-1 {
				list.Lists[i].Scores = list.Lists[i].Scores[:index]
			} else if index == 0 {
				list.Lists[i].Scores = list.Lists[i].Scores[1:]
			} else {
				list.Lists[i].Scores = append(list.Lists[i].Scores[:index], list.Lists[i].Scores[index+1:]...)
			}

			// Save updated list
			b, err := json.Marshal(list)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error parsing list.")
				return
			}
			err = ioutil.WriteFile("./lists/"+m.Author.ID+".json", b, 0644)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error saving list to file.")
				return
			}

			s.ChannelMessageSend(m.ChannelID, "Removed score!")
			return
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Could not find the specified list.")
}

// ListRunHandler lets the user run their list
func ListRunHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check for mention
	user := m.Author
	if len(m.Mentions) >= 1 {
		user = m.Mentions[0]
	}

	// Get file
	list, new := structs.GetList(user)

	if new {
		if m.Author.ID == user.ID {
			s.ChannelMessageSend(m.ChannelID, "Could not find a list with your user ID!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Could not find a list for the given user!")
		}
		return
	}

	// Create args foundation
	args := []string{"./osu-tools/delta/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll", "simulate", "osu"}

	// Check for other osu! versions
	if values.Jozregex.MatchString(m.Content) {
		args[0] = "./osu-tools/joz/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	} else if values.Liveregex.MatchString(m.Content) {
		args[0] = "./osu-tools/live/osu-tools/PerformanceCalculator/bin/Release/netcoreapp2.0/PerformanceCalculator.dll"
	}

	var (
		scores      []structs.Score
		subListName string
	)

	if values.Optionregex.MatchString(m.Content) {
		subListName = values.Optionregex.FindStringSubmatch(m.Content)[1]
	}

	if subListName != "" { // Get specific sublist
		for _, subList := range list.Lists {
			if strings.ToLower(subList.Name) == strings.ToLower(subListName) {
				scores = subList.Scores
				break
			}
		}
	} else { // Get all lists
		for _, subList := range list.Lists {
			scores = append(scores, subList.Scores...)
		}
	}

	msg, err := s.ChannelMessageSend(m.ChannelID, "Running "+m.Author.Username+"'s score request...")
	if err != nil {
		return
	}
	defer s.ChannelMessageDelete(m.ChannelID, msg.ID)
	log.Println(m.Author.String() + " has requested to run " + user.Username + "'s list.")

	// Remove dupes
	keys := make(map[structs.Score]bool)
	scoresNoDupe := []structs.Score{}
	for _, score := range scores {
		if _, exists := keys[score]; !exists {
			keys[score] = true
			scoresNoDupe = append(scoresNoDupe, score)
		}
	}

	// Run scores
	var PPScoreList []structs.PPScore
	for _, score := range scoresNoDupe {
		tempArgs := append(args, "./cache/"+strconv.Itoa(score.BeatmapID)+".osu")

		if score.UseAccuracy && score.Accuracy > 0 && score.Accuracy < 100 {
			tempArgs = append(tempArgs, "-a", strconv.FormatFloat(score.Accuracy, 'f', 2, 64))
		} else {
			if score.Goods > 0 {
				tempArgs = append(tempArgs, "-G", strconv.Itoa(score.Goods))
			}
			if score.Mehs > 0 {
				tempArgs = append(tempArgs, "-M", strconv.Itoa(score.Mehs))
			}
		}

		if score.Misses > 0 {
			tempArgs = append(tempArgs, "-X", strconv.Itoa(score.Misses))
		}

		if score.Combo > 0 {
			tempArgs = append(tempArgs, "--combo", strconv.Itoa(score.Combo))
		}

		for i := 0; i < len(score.Mods); i += 2 {
			tempArgs = append(tempArgs, "-m", string(score.Mods[i])+string(score.Mods[i+1]))
		}

		process := exec.Command("dotnet", tempArgs...)
		res, err := process.Output()
		if err != nil {
			process.Process.Kill()
			s.ChannelMessageSend(m.ChannelID, "Error in obtaining the pp calc for the score on beatmap ID "+strconv.Itoa(score.BeatmapID))
			continue
		}
		process.Process.Kill()

		ppVal, err := strconv.ParseFloat(values.PPparseregex.FindStringSubmatch(string(res))[1], 64)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error in parsing the pp calc for the score on beatmap ID "+strconv.Itoa(score.BeatmapID))
			continue
		}

		PPScoreList = append(PPScoreList, structs.PPScore{
			Score: score,
			PP:    ppVal,
		})
	}

	sort.Slice(PPScoreList, func(i, j int) bool { return PPScoreList[i].PP > PPScoreList[j].PP })

	// Create sublist table
	var tableData [][]string
	var builder strings.Builder
	table := tablewriter.NewWriter(&builder)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	table.SetHeader([]string{"Map Info", "PP", "Mods", "Combo", "Acc", "100s", "50s", "Misses"})

	// Add score info to tabledata
	var mainText string
	if subListName != "" {
		mainText += "Scores for the sublist: " + subListName + "\n"
	}
	for _, score := range PPScoreList {
		scoreData := []string{
			score.MapInfo,
			strconv.FormatFloat(score.PP, 'f', 2, 64),
			score.Mods,
			strconv.Itoa(score.Combo),
			"",
			"",
			"",
			"",
		}

		// Exceptions
		if score.Mods == "" {
			scoreData[2] = "NM"
		}
		if score.Combo == 0 {
			scoreData[3] = "FC"
		}
		if score.UseAccuracy {
			scoreData[4] = strconv.FormatFloat(score.Accuracy, 'f', 2, 64) + "%"
		} else {
			scoreData[5] = strconv.Itoa(score.Goods)
			scoreData[6] = strconv.Itoa(score.Mehs)
		}
		if score.Misses != 0 {
			scoreData[7] = strconv.Itoa(score.Misses)
		}

		tableData = append(tableData, scoreData)
	}

	// Add to table and write to string
	table.AppendBulk(tableData)
	table.Render()
	mainText += builder.String()

	go sendPaste(s, m, structs.NewPasteData(user.Username+"'s Scores", mainText), "")
}

// ListWhoHandler lets the bot owner see whose list is who
func ListWhoHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	var files []string
	err := filepath.Walk("./lists", func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error obtaining list files!")
		return
	}

	text := "List of users:\n"
	for _, file := range files {
		if strings.HasSuffix(file, ".json") {
			ID := strings.Replace(strings.Replace(file, "lists\\", "", -1), ".json", "", -1)
			user, err := s.User(ID)
			if err != nil {
				continue
			}
			list, _ := structs.GetList(user)

			listLengths := "("
			for _, subList := range list.Lists {
				listLengths += strconv.Itoa(len(subList.Scores)) + " "
			}
			listLengths = strings.TrimSpace(listLengths) + ")"

			text += "`" + user.ID + "` **" + user.Username + "**: " + strconv.Itoa(len(list.Lists)) + " sublist(s) " + listLengths + "\n"
		}
	}

	s.ChannelMessageSend(m.ChannelID, text)
}

// ListImportHandler lets you import a list
func ListImportHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Get new list
	list, new := structs.GetList(m.Author)
	if !new {
		s.ChannelMessageSend(m.ChannelID, "User already has a list!")
		return
	}

	log.Println(m.Author.String() + " has requested to import a list.")

	// Get legacy list's data
	res, err := http.Get(m.Attachments[0].URL)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Cannot access the discord file!")
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in parsing discord file!")
		return
	}
	var legacyList []structs.LegacyScore
	err = json.Unmarshal(b, &legacyList)
	if err != nil {
		// Attempt new list format
		var newList structs.List
		err = json.Unmarshal(b, &newList)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error in parsing JSON contents!")
			return
		}
		newList.User = *m.Author

		ioutil.WriteFile("./lists/"+m.Author.ID+".json", b, 0644)
		s.ChannelMessageSend(m.ChannelID, "Saved under `"+m.Author.ID+".json`!")
		return
	}

	// Put into new list
	for _, score := range legacyList {
		list.Lists[0].Scores = append(list.Lists[0].Scores, structs.Score{
			MapInfo:     score.MapInfo,
			BeatmapID:   score.BeatmapID,
			Accuracy:    score.Accuracy,
			Combo:       score.Combo,
			Misses:      score.Misses,
			Mods:        score.Mods,
			UseAccuracy: true,
		})
	}

	b, err = json.Marshal(list)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in stringifying new list format!")
		return
	}
	err = ioutil.WriteFile("./lists/"+m.Author.ID+".json", b, 0644)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error in saving new list!")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Saved under `"+m.Author.ID+".json`!")
}
