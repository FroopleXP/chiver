package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type fileMoveInstruction struct {
	filename string
	location string
}

func executeMoveCommands(commands *[]fileMoveInstruction) {
	for _, cmd := range *commands {
		fmt.Printf("Moving %s to %s\n", cmd.filename, cmd.location)
		e := exec.Command("mv", cmd.filename, cmd.location)
		e.Run()
	}
}

func getExtensionFromFilename(filename string) (string, error) {
	return filepath.Ext(filename)[1:], nil
}

func existsInStringArr(arr *[]string, str string) bool {
	found := false
	for _, i := range *arr {
		if i == str {
			found = true
			break
		}
	}
	return found
}

func createFilePaths(paths *[]string, start string) {
	for _, path := range *paths {
		dirName := fmt.Sprintf("%s/%s", start, path)
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created folder %s\n", dirName)
	}
}

func main() {

	const inputDirectory = "../mockdata/files"
	const outputDir = "../out"

	// Registering files to exclude
	notAllowed := []string{"DS_Store"}

	// Registering file types
	fileTypes := make(map[string]string)
	fileTypes["jpg"] = "image"
	fileTypes["png"] = "image"
	fileTypes["mp4"] = "video"
	fileTypes["gif"] = "image"

	// File move instructions
	fileMoveInstructions := make([]fileMoveInstruction, 0)

	/*

		Output:

		/<type>/<filetype>/<year>/<month>/<day>

		Example:
		images/jpeg/2021/3/21/

	*/

	files, err := ioutil.ReadDir(inputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// Holds paths to create
	pathsToCreate := make([]string, 0)

	for _, file := range files {

		// Checking if the file is a folder --- ignore for now
		if file.IsDir() {
			continue
		}

		ext, err := getExtensionFromFilename(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		// Check if we should ignore the file
		if existsInStringArr(&notAllowed, ext) {
			continue
		}

		dateCreated := file.ModTime()

		yearStr := dateCreated.Year()
		monthStr := dateCreated.Month().String()
		dayStr := dateCreated.Day()

		// Checking the type of file from the extension ie, jpg = image etc...
		fileType := fileTypes[ext]
		if fileType == "" {
			continue
		}

		newPath := fmt.Sprintf("%s/%d/%s/%d", fileType, yearStr, monthStr, dayStr)

		// Check if that file path has already been set to create
		if !existsInStringArr(&pathsToCreate, newPath) {
			pathsToCreate = append(pathsToCreate, newPath)
		}

		moveIns := fileMoveInstruction{inputDirectory + "/" + file.Name(), outputDir + "/" + newPath}

		fileMoveInstructions = append(fileMoveInstructions, moveIns)

	}

	// Create the new folders for the files
	createFilePaths(&pathsToCreate, outputDir)

	executeMoveCommands(&fileMoveInstructions)

}
