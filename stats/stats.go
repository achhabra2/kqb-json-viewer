package stats

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
		log.Fatal(err)

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

func ReadJson(file string) StatsJSON {
	var statsJSON StatsJSON
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Could not read json file", err)
	}
	err = json.Unmarshal(data, &statsJSON)
	if err != nil {
		log.Fatal("Could not parse json file", err)
	}

	return statsJSON
}
