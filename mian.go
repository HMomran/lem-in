package main

import (
	"Lem-in/Parse"
	"fmt"
	"os"
)

func main() {
	// <---Check if the correct number of arguments is provided---> \\
	file := os.Args[1:]
	if len(file) != 1 {
		fmt.Println("Please provide exactly one file as an argument.")
		os.Exit(0)
	}
	File := []rune(file[0])
	FileTypeCheck := []rune{}
	for i := len(File) - 1; i > len(File)-5; i-- {
		FileTypeCheck = append(FileTypeCheck, File[i])
	}
	if string(FileTypeCheck) == "txt." {
		Parse.Parsing(file[0])
	} else {
		fmt.Println("File type is not supported. Please provide a .txt file.")
		os.Exit(0)
	}

}
