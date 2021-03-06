package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func main() {
	names := stats.ListStatFiles()
	data := stats.ReadJson(names[0])

	timeWidget := widget.NewLabel(getTimeString(names[0]))
	a := app.New()
	w := a.NewWindow("KQB JSON Viewer")
	// a.SetIcon(resourceLogoPng)
	mapsWon := data.MapsWon()
	blueMapsLabel := widget.NewLabelWithStyle("Blue Maps", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	goldMapsLabel := widget.NewLabelWithStyle("Gold Maps", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	blueMaps := widget.NewLabelWithStyle(strconv.Itoa(mapsWon["blue"]), fyne.TextAlignCenter, fyne.TextStyle{})
	goldMaps := widget.NewLabelWithStyle(strconv.Itoa(mapsWon["gold"]), fyne.TextAlignCenter, fyne.TextStyle{})

	mapsContainer := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(2), goldMapsLabel, blueMapsLabel, goldMaps, blueMaps)
	players := BuildPlayerUI(data)
	trimmedMap := make(map[string]string)
	var trimmed []string
	for _, file := range names {
		_, name := filepath.Split(file)
		trimmedMap[name] = file
		trimmed = append(trimmed, name)
	}
	cont := container.NewVBox()

	combo := widget.NewSelect(trimmed, func(value string) {
		log.Println("Select file", value)
		selectedData := stats.ReadJson(trimmedMap[value])
		timeWidget = widget.NewLabel(getTimeString(trimmedMap[value]))
		players = BuildPlayerUI(selectedData)
		mapsWon = selectedData.MapsWon()
		blueMaps.SetText(strconv.Itoa(mapsWon["blue"]))
		goldMaps.SetText(strconv.Itoa(mapsWon["gold"]))
		cont.Objects[1] = timeWidget
		cont.Objects[2] = players
		cont.Refresh()
	})

	upload := widget.NewButton("Upload", func() {
	})

	advancedStatsButton := widget.NewButtonWithIcon("Adv. Stats", theme.FileImageIcon(), func() {
		selectedData := stats.ReadJson(trimmedMap[combo.Selected])
		ShowAdvancedStats(&a, selectedData)
	})

	nextButton := widget.NewButtonWithIcon("Next", theme.MediaSkipNextIcon(), func() {
		idx := combo.SelectedIndex()
		if idx+1 < len(trimmed) {
			combo.SetSelectedIndex(idx + 1)
		}
	})
	prevButton := widget.NewButtonWithIcon("Prev", theme.MediaSkipPreviousIcon(), func() {
		idx := combo.SelectedIndex()
		if idx-1 >= 0 {
			combo.SetSelectedIndex(idx - 1)
		}
	})

	openButton := widget.NewButtonWithIcon("Open", theme.FileIcon(), func() {
		selectedFile := trimmedMap[combo.Selected]
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", selectedFile).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", selectedFile).Start()
		case "darwin":
			err = exec.Command("open", selectedFile).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			log.Fatal(err)
		}
	})

	refreshButton := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() {
	})
	controls := container.NewHBox(layout.NewSpacer(), openButton, refreshButton, combo, prevButton, nextButton, advancedStatsButton, layout.NewSpacer())
	cont.Add(controls)
	cont.Add(timeWidget)
	cont.Add(players)
	cont.Add(mapsContainer)
	cont.Add(upload)

	combo.SetSelectedIndex(0)

	w.SetContent(cont)
	w.CenterOnScreen()
	w.ShowAndRun()
}

func BuildPlayerUI(data stats.StatsJSON) *fyne.Container {
	nameCont := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(1))
	cont := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(5))
	nameCont.Add(widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Kills", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Deaths", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Berries", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Snail", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Team", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	sort.Slice(data.PlayerMatchStats, func(i, j int) bool {
		return data.PlayerMatchStats[i].Team < data.PlayerMatchStats[j].Team
	})
	for _, player := range data.PlayerMatchStats {
		name := player.Nickname
		kills := strconv.Itoa(player.Kills)
		deaths := strconv.Itoa(player.Deaths)
		berries := strconv.Itoa(player.Berries)
		snail := strconv.FormatFloat(player.Snail, 'f', 0, 64)
		team := getTeam(player.Team)
		nameCont.Add(widget.NewLabel(name))
		cont.Add(widget.NewLabel(kills))
		cont.Add(widget.NewLabel(deaths))
		cont.Add(widget.NewLabel(berries))
		cont.Add(widget.NewLabel(snail))
		cont.Add(widget.NewLabel(team))
	}
	playerContainer := container.NewHBox(layout.NewSpacer(), nameCont, cont, layout.NewSpacer())
	return playerContainer
}

func getTimeString(file string) string {
	fInfo, _ := os.Open(file)
	info, _ := fInfo.Stat()
	timeStr := info.ModTime().String()
	return timeStr
}

func getTeam(team int) string {
	switch team {
	case 1:
		return "Gold"
	case 2:
		return "Blue"
	default:
		return ""
	}
}
