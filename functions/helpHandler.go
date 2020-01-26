package functions

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// HelpHandler handles the help function
func HelpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println(m.Author.String() + " has requested help.")
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    "PP Bot functions",
			IconURL: s.State.User.AvatarURL("2048"),
		},
		Description: "For help in map link/attachment pp/sr calculations, use `help map`\n" +
			"For help in user profile calculations, use `help user`",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "General commands:",
				Value: "`map`, " +
					"`user`, " +
					"`acc`",
			},
			&discordgo.MessageEmbedField{
				Name: "List commands:",
				Value: "`add`, " +
					"`run`, " +
					"`list`, " +
					"`import`, " +
					"`move`, " +
					"`delete`",
			},
		},
	}

	argRegex := regexp.MustCompile(`help\s+(.+)`)
	if argRegex.MatchString(m.Content) {
		arg := argRegex.FindStringSubmatch(m.Content)[1]

		switch arg { // Refer to functions below
		case "map":
			embed = helpMap(embed)
		case "user":
			embed = helpUser(embed)
		case "acc":
			embed = helpAcc(embed)

		case "add":
			embed = helpAdd(embed)
		case "run":
			embed = helpRun(embed)
		case "list":
			embed = helpList(embed)
		case "import":
			embed = helpImport(embed)
		case "move":
			embed = helpMove(embed)
		case "delete":
			embed = helpDelete(embed)
		}
	}

	if !strings.HasPrefix(embed.Description, "For") && embed.Fields[0].Name == "General commands:" {
		embed.Fields = []*discordgo.MessageEmbedField{}
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "Make sure to run `list` anytime you use `move` or `delete` whenever you want to continue moving / deleting to make sure you don't accidentally delete the wrong score!",
		Embed:   embed,
	})
}

// helpMap explains params for linking a map
func helpMap(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Description = "Sending a map link / attaching a .osu file will provide the star rating for it with no mods. There are parameters below which you can provide to change the value however you like."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "-pp",
			Value:  "Will provide the pp instead of star rating of the map.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "(-j|-l)",
			Value:  "Will calculate using either the joz system if `-j`, or live if `-l`. If both are provided, it will default to joz.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-m <mods>",
			Value:  "Will calculate using the given mods. They can only be a combination of EZ, NF, HT, HR, DT, NC, HD, FL, and SO.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-100 <100s>",
			Value:  "Will calculate PP using the given goods / 100s. It's recommend to use `-100` and `-50` instead of `-a` when using the delta PP system.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-50 <50s>",
			Value:  "Will calculate PP using the given mehs / 50s. It's recommend to use `-100` and `-50` instead of `-a` when using the delta PP system.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-a <accuracy>",
			Value:  "Will calculate PP using the given acc. It's recommend to use `-100` and `-50` instead of `-a` when using the delta PP system.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-c <combo>",
			Value:  "Will calculate PP using the given combo.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-x <misses>",
			Value:  "Will calculate PP using the given misses.",
			Inline: true,
		},
	}
	return embed
}

// helpUser explains params for linking a user
func helpUser(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Description = "Sending a user link will calculate the user's top 100 plays. If you wish to calculate your plays outside of the top 100, consider using a list."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "(-j|-l)",
			Value:  "Will calculate using either the joz system if `-j`, or live if `-l`. If both are provided, it will default to joz.",
		},
	}
	return embed
}

// helpUser explains params for linking a user
func helpAcc(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!acc"
	embed.Description = "`acc` will create an acc graph. You may customize using the parameters below."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "(-j|-l)",
			Value:  "Will calculate using either the joz system if `-j`, or live if `-l`. If both are provided, it will default to joz.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-m <mods>",
			Value:  "Will create an acc graph using the given mods. They can only be a combination of EZ, NF, HT, HR, DT, NC, HD, FL, and SO.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-c <combo>",
			Value:  "Will create an acc graph using the given combo.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-x <misses>",
			Value:  "Will create an acc graph using the given misses.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-a <low end> <high end>",
			Value:  "Will create an acc graph from <low end> acc to <high end> acc.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-i <value>",
			Value:  "Will calculate the pp every <value>%+.",
			Inline: true,
		},
	}
	return embed
}

// helpAdd explains params for linking a user
func helpAdd(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!add"
	embed.Description = "`add` will add a map to a list."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "(-j|-l)",
			Value:  "Will calculate using either the joz system if `-j`, or live if `-l`. If both are provided, it will default to joz. This is only for calculating the pp value if the map is successfully added to your list. You may calculate your list however you like later using `run`.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-o <list name>",
			Value:  "The list to add to. If this parameter is not used, the map will be added to a list called 'Untitled'. If the specified list doesn't exist, it will be created.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-m <mods>",
			Value:  "Will add the given mods. They can only be a combination of EZ, NF, HT, HR, DT, NC, HD, FL, and SO.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-100 <100s>",
			Value:  "Will add the given goods / 100s. It's recommend to use `-100` and `-50` instead of `-a` when using the delta PP system.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-50 <50s>",
			Value:  "Will add the given mehs / 50s. It's recommend to use `-100` and `-50` instead of `-a` when using the delta PP system.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-a <accuracy>",
			Value:  "Will add the given acc. It's recommend to use `-100` and `-50` instead of `-a` when using the delta PP system.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-c <combo>",
			Value:  "Will add the given combo.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-x <misses>",
			Value:  "Will add the given misses.",
			Inline: true,
		},
	}
	return embed
}

// helpRun explains params for running a list
func helpRun(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!run"
	embed.Description = "`run` will run a sublist or all lists combined."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "(-j|-l)",
			Value:  "Will calculate using either the joz system if `-j`, or live if `-l`. If both are provided, it will default to joz.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-o <list name>",
			Value:  "The list to run. If this parameter is not used, all lists will be combined, and then run.",
			Inline: true,
		},
	}
	return embed
}

// helpList explains params for showing a list
func helpList(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!list"
	embed.Description = "`list` will show all sublists."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "@user",
			Value:  "The user's list to show. If nooone is mentioned, it will show their own list instead.",
		},
	}
	return embed
}

// helpImport explains params for showing a list
func helpImport(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!import"
	embed.Description = "`import` will import a list given a JSON file. The JSON file should be in one of the following formats below.\n**Import will not work if you already have a list.**"
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "JSON file format 1",
			Value:  "```JSON\n" + 
				"[\n" + 
				"\t{\n" + 
				"\t\t\"mapinfo\": \"artist - title [diff name]\",\n" + 
				"\t\t\"beatmapid\": beatmapID\n" + 
				"\t\t\"accuracy\": accuracy\n" + 
				"\t\t\"combo\": combo\n" + 
				"\t\t\"misses\": misses\n" + 
				"\t\t\"mods\": \"mods\"\n" + 
				"\t},\n" +
				"\t{...}, ...\n" +
				"]```",
		},
		&discordgo.MessageEmbedField{
			Name:   "JSON file format 2",
			Value:  "```JSON\n" + 
				"{\n" + 
				"\t\"User\": {},\n" + 
				"\t\"Lists\": [\n" + 
				"\t\t{\n" + 
				"\t\t\t\"Name\": \"Name\",\n" + 
				"\t\t\t\"Scores\": [\n" +
				"\t\t\t\t{\n" + 
				"\t\t\t\t\t\"MapInfo\": \"artist - title [diff name]\",\n" +
				"\t\t\t\t\t\"BeatmapID\": beatmapID,\n" +
				"\t\t\t\t\t\"Accuracy\": accuracy,\n" +
				"\t\t\t\t\t\"Goods\": 100s,\n" +
				"\t\t\t\t\t\"Mehs\": 50s,\n" +
				"\t\t\t\t\t\"Combo\": combo,\n" +
				"\t\t\t\t\t\"Misses\": misses,\n" +
				"\t\t\t\t\t\"Mods\": \"mods\",\n" +
				"\t\t\t\t\t\"UseAccuracy\": true / false\n" +
				"\t\t\t\t},\n" + 
				"\t\t\t\t{...}, ...\n" + 
				"\t\t\t]\n" + 
				"\t\t},\n" +
				"\t\t{...}, ...\n" +
				"\t]\n" + 
				"}```",
		},
	}
	return embed
}

// helpMove explains params for moving a score between lists
func helpMove(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!move"
	embed.Description = "`move` will move a score from 1 list to another."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "-o <list to move from>",
			Value:  "The list to move the score from. If this parameter is not used, it will look in the `Untitled` list.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-n <index>",
			Value:  "The index of the map in the list to move from. You can obtain the index number via `list`",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "-t <list to move to>",
			Value:  "The list to move the score to. If the list doesn't exist, it will be created.",
			Inline: true,
		},
	}
	return embed
}

// helpDelete explains params for moving a score between lists
func helpDelete(embed *discordgo.MessageEmbed) *discordgo.MessageEmbed {
	embed.Author.Name = "!delete"
	embed.Description = "`delete` will remove a score, a sublist, or the whole list."
	embed.Fields = []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:   "-all -o <sublist name>",
			Value:  "Remove a sublist. If `-o <sublist name>` is not given (so only `-all` is given), the full list will be deleted.",
			Inline: true,
		},
		&discordgo.MessageEmbedField{
			Name:   "<index> -o <sublist name>",
			Value:  "The index of the map in the list to move from. If `-o <sublist name>` is not given, it will look in the `Untitled` list. You can obtain the index number via `list`",
			Inline: true,
		},
	}
	return embed
}