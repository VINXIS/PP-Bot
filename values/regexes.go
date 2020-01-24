package values

import "regexp"

var (
	// OutsideServerregex checks if a regex was from outside the PP server
	OutsideServerregex = regexp.MustCompile(`-g`)

	// Addregex checks if a command was made to add a score to their list
	Addregex = regexp.MustCompile(`!add`)
	// Moveregex checks if a command was made to move a score from one sublist to another
	Moveregex = regexp.MustCompile(`!move`)
	// Accgraphregex checks if a command was made to create an acc graph for the map
	Accgraphregex = regexp.MustCompile(`!acc`)
	// Runregex checks if a command was made to run their list
	Runregex = regexp.MustCompile(`!run`)
	// Listregex checks if a command was made to show their list
	Listregex = regexp.MustCompile(`!list`)
	// Whoregex checks if a command was made to show who has a list
	Whoregex = regexp.MustCompile(`!wholist`)
	// Delregex checks if a command was made to delete a score from their list
	Delregex = regexp.MustCompile(`!(d(el(ete)?)?|r(emove)?)(\s+(\d+)?)`)

	// Mapregex checks if a map was linked
	Mapregex = regexp.MustCompile(`(osu|old)\.ppy\.sh\/(s|b|beatmaps|beatmapsets)\/(\d+)(#osu\/(\d+))?`)
	// Userregex checks if a user was linked
	Userregex = regexp.MustCompile(`(osu|old)\.ppy\.sh\/(u|users)\/(\S+)`)
	// Fileregex checks if a map was attached
	Fileregex = regexp.MustCompile(`\.osu`)

	// PPregex checks if a tag for pp values was used
	PPregex = regexp.MustCompile(`-pp`)

	// Titleregex gets the title from a .osu file
	Titleregex = regexp.MustCompile(`Title:(.*)(\r|\n)`)
	// Artistregex gets the artist from a .osu file
	Artistregex = regexp.MustCompile(`Artist:(.*)(\r|\n)`)
	// Versionregex gets the diff name from a .osu file
	Versionregex = regexp.MustCompile(`Version:(.*)(\r|\n)`)
	// BeatmapIDregex gets the beatmap id from a .osu file
	BeatmapIDregex = regexp.MustCompile(`BeatmapID:(.*)(\r|\n)`)

	// PPparseregex gets the pp from a pp calc result
	PPparseregex = regexp.MustCompile(`pp\s+:\s+((\d+)\.?\d+)`)
	// Aimregex gets the aim pp from a pp calc result
	Aimregex = regexp.MustCompile(`Aim\s+:\s+((\d+)\.?\d+)`)
	// Tapregex gets the tap pp from a pp calc result
	Tapregex = regexp.MustCompile(`Tap\s+:\s+((\d+)\.?\d+)`)
	// accPPregex gets the acc pp from a pp calc result
	AccPPregex = regexp.MustCompile(`Accuracy\s+:\s+((\d+)\.?\d+)`)

	// Spamfileregex checks for the spam files created by the custom osu-tools
	Spamfileregex = regexp.MustCompile(`\d*(jumpaim|speed|stamina|streamaim|aimcontrol|fingercontrol|values|test).txt`)

	// Modregex looks for mods
	Modregex = regexp.MustCompile(`(?i)-m\s+((?:EZ|NF|HT|HR|DT|NC|HD|FL|SO)+)`)
	// Accregex looks for acc
	Accregex = regexp.MustCompile(`-a\s+((\d|\.)+)`)
	// Incrementregex looks for the increment to use for acc graphs
	Incrementregex = regexp.MustCompile(`-i\s+((\d|\.)+)`)
	// Comboregex looks for combo
	Comboregex = regexp.MustCompile(`-c\s+(\d+)`)
	// Missregex looks for misses
	Missregex = regexp.MustCompile(`-x\s+(\d+)`)
	// Goodregex looks for 100s
	Goodregex = regexp.MustCompile(`-100\s+(\d+)`)
	// Mehregex looks for 50s
	Mehregex = regexp.MustCompile(`-50\s+(\d+)`)
	// Skillregex looks for a specific skill
	Skillregex = regexp.MustCompile(`(?i)-s\s+(Aim\s*Control|Jump\s*Aim|Stream\s*Aim|Finger\s*Control|Speed|Stamina)`)
	// Accrangeregex looks for a specific acc range for acc graphs
	Accrangeregex = regexp.MustCompile(`-a\s+(\d+)\s+(\d+)`)
	// Timeregex looks for a specific time range for the graphs
	Timeregex = regexp.MustCompile(`-r\s+(\d+)\s+(\d+)`)

	// Jozregex will run joz instead of delta
	Jozregex = regexp.MustCompile(`-j`)
	// Liveregex will run live instead of delta
	Liveregex = regexp.MustCompile(`-l`)
)
