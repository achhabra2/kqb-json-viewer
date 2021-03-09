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
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

type KQBApp struct {
	files        []string
	selectedData stats.StatsJSON
	a            fyne.App
	w            fyne.Window
}

func (k *KQBApp) ShowMainWindow() {
	timeWidget := widget.NewLabel(getTimeString(k.files[0]))
	// a.SetIcon(resourceLogoPng)
	about := fyne.NewMenuItem("About", func() {
		aboutMessage := fmt.Sprintf("kqb-json-viewer version %s \n by Prosive", version)
		dialog := dialog.NewInformation("About", aboutMessage, k.w)
		dialog.Show()
	})
	fileMenu := fyne.NewMenu("File", about)
	mainMenu := fyne.NewMainMenu(fileMenu)
	mapsWon := k.selectedData.MapsWon()
	blueMapsLabel := widget.NewLabelWithStyle("Blue Maps", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	goldMapsLabel := widget.NewLabelWithStyle("Gold Maps", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	blueMaps := widget.NewLabelWithStyle(strconv.Itoa(mapsWon["blue"]), fyne.TextAlignCenter, fyne.TextStyle{})
	goldMaps := widget.NewLabelWithStyle(strconv.Itoa(mapsWon["gold"]), fyne.TextAlignCenter, fyne.TextStyle{})

	mapsContainer := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(2), goldMapsLabel, blueMapsLabel, goldMaps, blueMaps)
	players := k.BuildPlayerUI()
	trimmedMap := make(map[string]string)
	var trimmed []string
	for _, file := range k.files {
		_, name := filepath.Split(file)
		trimmedMap[name] = file
		trimmed = append(trimmed, name)
	}
	cont := container.NewVBox()

	combo := widget.NewSelect(trimmed, func(value string) {
		log.Println("Select file", value)
		k.selectedData = stats.ReadJson(trimmedMap[value])
		timeWidget = widget.NewLabel(getTimeString(trimmedMap[value]))
		players = k.BuildPlayerUI()
		mapsWon = k.selectedData.MapsWon()
		blueMaps.SetText(strconv.Itoa(mapsWon["blue"]))
		goldMaps.SetText(strconv.Itoa(mapsWon["gold"]))
		cont.Objects[1] = timeWidget
		cont.Objects[2] = players
		cont.Refresh()
	})

	upload := widget.NewButton("Upload", func() {
		ShowUploadWindow(k.a, k.selectedData)
	})

	advancedStatsButton := widget.NewButtonWithIcon("Adv. Stats", theme.FileImageIcon(), func() {
		k.selectedData = stats.ReadJson(trimmedMap[combo.Selected])
		k.ShowAdvancedStats()
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

	controls := container.NewHBox(layout.NewSpacer(), openButton, combo, prevButton, nextButton, advancedStatsButton, layout.NewSpacer())
	cont.Add(controls)
	cont.Add(timeWidget)
	cont.Add(players)
	cont.Add(mapsContainer)
	cont.Add(upload)

	combo.SetSelectedIndex(0)

	k.w.SetContent(cont)
	k.w.SetMainMenu(mainMenu)
	k.w.CenterOnScreen()
	go k.UpdateCheckUI()
	k.w.ShowAndRun()
}

func (k *KQBApp) UpdateCheckUI() {
	w := k.w
	shouldUpdate, latestVersion := checkForUpdate()
	if shouldUpdate {
		updateMessage := fmt.Sprintf("New Version Available, would you like to update to v%s", latestVersion)
		confirmDialog := dialog.NewConfirm("Update Checker", updateMessage, func(action bool) {
			if action {
				log.Println("Update clicked")
				updated := doSelfUpdate()
				if updated {
					updatedDialog := dialog.NewInformation("Update Status", "Update Succeeded, please restart", w)
					updatedDialog.Show()
				} else {
					updatedDialog := dialog.NewInformation("Update Status", "Update Failed", w)
					updatedDialog.Show()
				}
			}
		}, w)
		confirmDialog.Show()
	}
}

func (k *KQBApp) ShowAdvancedStats() {
	a := k.a
	data := k.selectedData
	w := a.NewWindow("Advanced Stats")

	advStatsPlot := stats.PlotStats(data)
	objStatsPlot := stats.PlotObjectiveStats(data)
	advStatsCanvas := canvas.NewImageFromImage(advStatsPlot)
	advStatsCanvas.SetMinSize(fyne.NewSize(1280, 720))
	objStatsCanvas := canvas.NewImageFromImage(objStatsPlot)
	objStatsCanvas.SetMinSize(fyne.NewSize(1280, 720))

	advStatsCanvas.FillMode = canvas.ImageFillContain
	objStatsCanvas.FillMode = canvas.ImageFillContain
	cont := container.NewVBox()
	nextButton := widget.NewButtonWithIcon("Military", theme.MediaSkipNextIcon(), func() {
		cont.Objects[0] = advStatsCanvas
		cont.Refresh()
	})
	prevButton := widget.NewButtonWithIcon("Objective", theme.MediaSkipPreviousIcon(), func() {
		cont.Objects[0] = objStatsCanvas
		cont.Refresh()
	})
	controls := container.NewHBox(layout.NewSpacer(), prevButton, nextButton, layout.NewSpacer())
	cont.Add(advStatsCanvas)
	cont.Add(controls)
	w.CenterOnScreen()
	w.SetContent(cont)
	w.Show()
}

func (k *KQBApp) BuildPlayerUI() *fyne.Container {
	data := k.selectedData
	nameCont := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(1))
	cont := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(6))
	nameCont.Add(widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Kills", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Deaths", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Berries", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Snail", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Team", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Entity", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
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
		entity := getEntity(player.EntityType)
		nameCont.Add(widget.NewLabel(name))
		cont.Add(widget.NewLabel(kills))
		cont.Add(widget.NewLabel(deaths))
		cont.Add(widget.NewLabel(berries))
		cont.Add(widget.NewLabel(snail))
		cont.Add(widget.NewLabel(team))
		cont.Add(widget.NewLabel(entity))
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

func getEntity(entity int) string {
	switch entity {
	case 3:
		return "Queen"
	default:
		return "Worker"
	}
}
