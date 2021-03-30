package stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
)

var WindowsDirectory string = "/AppData/LocalLow/Liquid Bit, LLC/Killer Queen Black/match-stats/"
var MacDirectory string = "/Library/Application Support/com.LiquidBit.KillerQueenBlack/match-stats/"

func ListStatFiles() []string {
	homeDir, _ := os.UserHomeDir()
	var statsPath string
	switch runtime.GOOS {
	case "windows":
		statsPath = filepath.Join(homeDir, WindowsDirectory)
	case "darwin":
		statsPath = filepath.Join(homeDir, MacDirectory)
	default:
		statsPath = filepath.Join(homeDir, WindowsDirectory)
	}
	files, err := ioutil.ReadDir(statsPath)
	if err != nil {
		// could not read KQB stats directory
		log.Println(err)

		// As an error path check the current directory for json files
		statsPath, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		files, err = ioutil.ReadDir(statsPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	files = sortFiles(files)
	var fileNames []string
	for _, file := range files {
		// Only care about json files
		if filepath.Ext(file.Name()) == ".json" {
			fullPath := filepath.Join(statsPath, file.Name())
			fileNames = append(fileNames, fullPath)
		}
	}
	return fileNames
}

func sortFiles(files []os.FileInfo) []os.FileInfo {
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})
	return files
}

func ReadJson(file string) (StatsJSON, error) {
	var statsJSON StatsJSON
	data, err := ioutil.ReadFile(file)
	_, fileName := filepath.Split(file)
	if err != nil {
		jsonError := fmt.Errorf("could not read json file\n (%v): %v", fileName, err)
		log.Println(jsonError.Error())
		return statsJSON, jsonError
	}
	err = json.Unmarshal(data, &statsJSON)
	if err != nil {
		marshalError := fmt.Errorf("could not properly decode json file\n (%v): %v", fileName, err)
		log.Println(marshalError.Error())
		return statsJSON, marshalError
	}

	if len(statsJSON.PlayerMatchStats) == 0 {
		invalidError := fmt.Errorf("invalid KQB Stats JSON File (%v)", fileName)
		log.Println(invalidError.Error())
		return statsJSON, invalidError
	}
	return statsJSON, nil
}

func OpenStatDirectory() {
	homeDir, _ := os.UserHomeDir()
	var statsPath string
	switch runtime.GOOS {
	case "windows":
		statsPath = filepath.Join(homeDir, WindowsDirectory)
	case "darwin":
		statsPath = filepath.Join(homeDir, MacDirectory)
	default:
		statsPath = filepath.Join(homeDir, WindowsDirectory)
	}

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", statsPath).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", statsPath).Start()
	case "darwin":
		err = exec.Command("open", statsPath).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println(err)
	}

}
