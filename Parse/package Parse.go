package Parse

import (
    "bufio"
    "fmt"
    "os"
    "sort"
    "strconv"
    "strings"
)

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

    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    for i := 1; i < len(lines); i++ {
        input := lines[i]
        if strings.Contains(input, "-") {
            continue
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