package Parse

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseFile(File string) (map[string]*Room, int, []string) {
	rooms := make(map[string]*Room)
	file, err := os.Open(File)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(0)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
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

	return rooms, numberOfAnts, lines
}

func parseRooms(lines []string, rooms map[string]*Room) {
	var Isstart bool
	var Isend bool

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
}

func parseTunnels(lines []string, rooms map[string]*Room) {
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
}
