package Parse

func Parsing(File string) {
	rooms, numberOfAnts, lines := parseFile(File)
	parseRooms(lines, rooms)
	parseTunnels(lines, rooms)
	start(rooms, numberOfAnts)
}

type path struct {
	roomNames []string
	ISPath    bool
}

type Room struct {
	Name      string
	x         int
	y         int
	IsStart   bool
	IsEnd     bool
	Neighbors []*Room
}

func Parsing(File string) {
	rooms := make(map[string]*Room)
	file, err := os.Open(File)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(0)
	}
	defer file.Close()

	var Isstart bool
	var Isend bool

	scanner := bufio.NewScanner(file)
	var lines []string

	// Read all lines first
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Parse rooms first
	for i := 1; i < len(lines); i++ {
		input := lines[i]
		if strings.Contains(input, "-") {
			continue // Skip tunnel lines for now
		}

		parts := strings.Split(input, " ")
		if len(parts) == 1 {
			if input == "##start" {
				Isstart = true
			} else if input == "##end" {
				Isend = true
			}
			continue
		} else if len(parts) >= 3 {
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			room := Room{
				Name:    parts[0],
				x:       x,
				y:       y,
				IsStart: Isstart,
				IsEnd:   Isend,
			}
			rooms[room.Name] = &room
			Isend = false
			Isstart = false
		}
	}

	// Parse tunnels/connections
	for i := 0; i < len(lines); i++ {
		input := lines[i]
		if strings.Contains(input, "-") {
			parts := strings.Split(input, "-")
			if len(parts) != 2 {
				fmt.Println("error: invalid tunnel format")
				continue
			}
			fromRoom := parts[0]
			toRoom := parts[1]

			FromRoom, exists1 := rooms[fromRoom]
			ToRoom, exists2 := rooms[toRoom]

			if !exists1 || !exists2 {
				fmt.Println("error: room not found")
				continue
			}

			FromRoom.Neighbors = append(FromRoom.Neighbors, ToRoom)
			ToRoom.Neighbors = append(ToRoom.Neighbors, FromRoom)

		}
	}

	numberOfAnts, err := strconv.Atoi(lines[0])
	if err != nil {
		fmt.Println("Error converting first line to number:", err)
		os.Exit(0)
	}
	if numberOfAnts <= 0 {
		fmt.Println("The number of ants must be greater than 0.")
		os.Exit(0)
	}
	start(rooms, numberOfAnts)
}

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

	// Find ALL paths from start to end
	var allPaths []path
	visited := make(map[*Room]bool)
	var currentPath []string

	findAllPaths(start, end, visited, &currentPath, &allPaths)

	if len(allPaths) > 0 {
		// Sort paths by length (shortest first)
		sort.Slice(allPaths, func(i, j int) bool {
			return len(allPaths[i].roomNames) < len(allPaths[j].roomNames)
		})

		// Use the paths for simulation
		simulateAnts(allPaths, numberOfAnts)
	} else {
		fmt.Println("ERROR")
	}
}

func findAllPaths(current *Room, end *Room, visited map[*Room]bool, currentPath *[]string, allPaths *[]path) {
	// Add current room to the path
	*currentPath = append(*currentPath, current.Name)

	// If we reached the end room, save this path
	if current == end {
		newPath := path{
			roomNames: make([]string, len(*currentPath)),
			ISPath:    true,
		}
		copy(newPath.roomNames, *currentPath)
		*allPaths = append(*allPaths, newPath)

		// Remove current room and return (don't mark as visited permanently for this branch)
		*currentPath = (*currentPath)[:len(*currentPath)-1]
		return
	}

	// Mark current room as visited for this path
	visited[current] = true

	// Try all unvisited neighbors
	for _, neighbor := range current.Neighbors {
		if !visited[neighbor] {
			findAllPaths(neighbor, end, visited, currentPath, allPaths)
		}
	}

	// Backtrack: unmark current room and remove from path
	visited[current] = false
	*currentPath = (*currentPath)[:len(*currentPath)-1]
}

func simulateAnts(allPaths []path, numberOfAnts int) {
	if len(allPaths) == 0 {
		fmt.Println("No paths available for simulation")
		return
	}

	// For simplicity, let's use the two shortest paths
	// In a full implementation, you'd use more sophisticated path selection
	var usablePaths []path
	for _, p := range allPaths {
		if len(usablePaths) < 2 { // Use up to 2 paths
			usablePaths = append(usablePaths, p)
		}
	}

	// Track ant positions: map[antID] = [pathIndex, positionInPath]
	antPositions := make(map[int][2]int)
	antsFinished := 0
	turn := 1

	for antsFinished < numberOfAnts {
		moves := []string{}

		// Move existing ants forward
		for antID := 1; antID <= numberOfAnts; antID++ {
			if pos, exists := antPositions[antID]; exists {
				pathIndex := pos[0]
				currentPos := pos[1]
				path := usablePaths[pathIndex]

				// If not at end, move forward
				if currentPos < len(path.roomNames)-1 {
					newPos := currentPos + 1
					antPositions[antID] = [2]int{pathIndex, newPos}
					moves = append(moves, fmt.Sprintf("L%d-%s", antID, path.roomNames[newPos]))

					// Check if reached end
					if newPos == len(path.roomNames)-1 {
						antsFinished++
					}
				}
			}
		}

		// Start new ants on available paths
		for antID := 1; antID <= numberOfAnts; antID++ {
			if _, exists := antPositions[antID]; !exists {
				// Find an available path (start position not occupied)
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
				break // Only start one ant per turn
			}
		}

		// Print moves for this turn
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}

		turn++

		// Safety check
		if turn > 50 {
			break
		}
	}
}

func dfs(current *Room, end *Room, visited map[*Room]bool) bool {
	// If we reached the end room
	if current == end {
		return true
	}

	// Mark current room as visited
	visited[current] = true

	// Try all unvisited neighbors
	for _, neighbor := range current.Neighbors {
		if !visited[neighbor] {
			if dfs(neighbor, end, visited) {
				return true
			}
		}
	}

	// No path found through this room
	return false
}

func dfsWithPath(current *Room, end *Room, visited map[*Room]bool, currentPath *[]string, foundPath *path) bool {
	// Add current room to the path
	*currentPath = append(*currentPath, current.Name)

	// If we reached the end room, save the path
	if current == end {
		foundPath.roomNames = make([]string, len(*currentPath))
		copy(foundPath.roomNames, *currentPath)
		foundPath.ISPath = true
		return true
	}

	// Mark current room as visited
	visited[current] = true

	// Try all unvisited neighbors
	for _, neighbor := range current.Neighbors {
		if !visited[neighbor] {
			if dfsWithPath(neighbor, end, visited, currentPath, foundPath) {
				return true
			}
		}
	}

	// Remove current room from path (backtrack)
	*currentPath = (*currentPath)[:len(*currentPath)-1]
	return false
}

func printPath(p path) {
	if p.ISPath {
		fmt.Print("Path found: ")
		for i, roomName := range p.roomNames {
			if i > 0 {
				fmt.Print(" -> ")
			}
			fmt.Print(roomName)
		}
		fmt.Println()
	} else {
		fmt.Println("No path stored")
	}
}
