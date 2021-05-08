package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Println(err)
	}

}

func FileNameToTime(basePath string, filename string) time.Time {
	arr := strings.Split(filename, "-")
	timeArr := make([]int, 0)
	if len(arr) == 7 {
		for _, val := range arr[1:] {
			intVal, _ := strconv.Atoi(val)
			timeArr = append(timeArr, intVal)
		}
		currentLocation := time.Now().Location()
		t := time.Date(timeArr[0], time.Month(timeArr[1]), timeArr[2], timeArr[3], timeArr[4], timeArr[5], 0, currentLocation)
		return t
	} else {
		fullFilePath := filepath.Join(basePath, filename)
		fInfo, _ := os.Open(fullFilePath)
		info, _ := fInfo.Stat()
		return info.ModTime()
	}
}
