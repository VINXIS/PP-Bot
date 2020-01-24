package functions

import (
	"os"
	"strings"
)

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
