package Parse

import (
	"fmt"
	"strings"
)

func simulateAnts(allPaths []path, numberOfAnts int) {
	if len(allPaths) == 0 {
		fmt.Println("No paths available for simulation")
		return
	}

	var usablePaths []path
	for _, p := range allPaths {
		if len(usablePaths) < 2 {
			usablePaths = append(usablePaths, p)
		}
	}

	antPositions := make(map[int][2]int)
	antsFinished := 0

	for antsFinished < numberOfAnts {
		moves := []string{}

		for antID := 1; antID <= numberOfAnts; antID++ {
			if pos, exists := antPositions[antID]; exists {
				pathIndex := pos[0]
				currentPos := pos[1]
				path := usablePaths[pathIndex]

				if currentPos < len(path.roomNames)-1 {
					newPos := currentPos + 1
					antPositions[antID] = [2]int{pathIndex, newPos}
					moves = append(moves, fmt.Sprintf("L%d-%s", antID, path.roomNames[newPos]))

					if newPos == len(path.roomNames)-1 {
						antsFinished++
					}
				}
			}
		}

		for antID := 1; antID <= numberOfAnts; antID++ {
			if _, exists := antPositions[antID]; !exists {
				for pathIndex, path := range usablePaths {
					startOccupied := false
					for checkAnt := 1; checkAnt <= numberOfAnts; checkAnt++ {
						if pos, exists := antPositions[checkAnt]; exists {
							if pos[0] == pathIndex && pos[1] == 0 {
								startOccupied = true
								break
							}
						}
					}

					if !startOccupied {
						antPositions[antID] = [2]int{pathIndex, 0}
						moves = append(moves, fmt.Sprintf("L%d-%s", antID, path.roomNames[0]))
						break
					}
				}
				break
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}
