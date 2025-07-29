package Parse

import (
	"fmt"
	"sort"
)

func start(rooms map[string]*Room, numberOfAnts int) {
	var start *Room
	var end *Room

	for _, room := range rooms {
		if room.IsStart == true {
			start = room
		} else if room.IsEnd == true {
			end = room
		}
	}

	if start == nil || end == nil {
		fmt.Println("Error: Start or End room not found")
		return
	}

	var allPaths []path
	visited := make(map[*Room]bool)
	var currentPath []string

	findAllPaths(start, end, visited, &currentPath, &allPaths)

	if len(allPaths) > 0 {
		sort.Slice(allPaths, func(i, j int) bool {
			return len(allPaths[i].roomNames) < len(allPaths[j].roomNames)
		})
		simulateAnts(allPaths, numberOfAnts)
	} else {
		fmt.Println("ERROR")
	}
}

func findAllPaths(current *Room, end *Room, visited map[*Room]bool, currentPath *[]string, allPaths *[]path) {
	*currentPath = append(*currentPath, current.Name)

	if current == end {
		newPath := path{
			roomNames: make([]string, len(*currentPath)),
			ISPath:    true,
		}
		copy(newPath.roomNames, *currentPath)
		*allPaths = append(*allPaths, newPath)
		*currentPath = (*currentPath)[:len(*currentPath)-1]
		return
	}

	visited[current] = true

	for _, neighbor := range current.Neighbors {
		if !visited[neighbor] {
			findAllPaths(neighbor, end, visited, currentPath, allPaths)
		}
	}

	visited[current] = false
	*currentPath = (*currentPath)[:len(*currentPath)-1]
}
