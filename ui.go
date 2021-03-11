package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/achhabra2/kqb-json-viewer/bgl"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

var goldColor = color.RGBA{255, 179, 0, 1}
var blueColor = color.RGBA{43, 93, 255, 1}

type KQBApp struct {
	files          []string
	selectedData   stats.StatsJSON
	a              fyne.App
	w              fyne.Window
	mainContainer  *container.Split
	splitContainer *fyne.Container
	u              *Uploader
	submission     bgl.Result
	subData        []bgl.SetMap
}

func (k *KQBApp) ShowMainWindow() {
	k.u = &Uploader{}
	timeWidget := widget.NewLabel(getTimeString(k.files[0]))
	timeContainer := container.NewHBox(layout.NewSpacer(), timeWidget, layout.NewSpacer())
	// a.SetIcon(resourceLogoPng)
	about := fyne.NewMenuItem("About", func() {
		aboutMessage := fmt.Sprintf("kqb-json-viewer version %s \n by Prosive", version)
		dialog := dialog.NewInformation("About", aboutMessage, k.w)
		dialog.Show()
	})
	fileMenu := fyne.NewMenu("File", about)
	mainMenu := fyne.NewMainMenu(fileMenu)

	mapsContainer := k.BuildMapTable()
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
		if k.u.BGLToken != "" {
			k.u.data = k.selectedData
		}
		timeWidget = widget.NewLabel(getTimeString(trimmedMap[value]))
		players = k.BuildPlayerUI()

		timeContainer.Objects[1] = timeWidget
		cont.Objects[2] = players
		cont.Objects[3] = k.BuildMapTable()
		cont.Refresh()
	})

	upload := widget.NewButton("Add Set to Match Result", func() {
		k.ShowUploadWindow()
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
	cont.Add(timeContainer)
	cont.Add(players)
	cont.Add(mapsContainer)
	cont.Add(upload)

	combo.SetSelectedIndex(0)
	trailingContainer := container.NewVBox(layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer())
	vSplitContainer := container.NewVScroll(trailingContainer)
	mainContainer := container.NewHSplit(cont, vSplitContainer)
	k.mainContainer = mainContainer
	k.splitContainer = trailingContainer
	k.w.SetContent(k.mainContainer)
	k.w.SetMainMenu(mainMenu)
	k.w.SetPadded(true)
	k.w.Resize(fyne.NewSize(600, 600))
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
		//		cont.Refresh()
	})
	prevButton := widget.NewButtonWithIcon("Objective", theme.MediaSkipPreviousIcon(), func() {
		cont.Objects[0] = objStatsCanvas
		//		cont.Refresh()
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
	cont := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(5))
	nameCont.Add(widget.NewLabelWithStyle("Name", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Kills", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Deaths", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Berries", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	cont.Add(widget.NewLabelWithStyle("Snail", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	cont.Add(widget.NewLabelWithStyle("Type", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	sort.Slice(data.PlayerMatchStats, func(i, j int) bool {
		return data.PlayerMatchStats[i].Team < data.PlayerMatchStats[j].Team
	})
	for _, player := range data.PlayerMatchStats {
		var col color.RGBA
		team := getTeam(player.Team)
		if team == "Blue" {
			col = blueColor
		} else {
			col = goldColor
		}
		name := player.Nickname
		kills := strconv.Itoa(player.Kills)
		deaths := strconv.Itoa(player.Deaths)
		berries := strconv.Itoa(player.Berries)
		snail := strconv.FormatFloat(player.Snail, 'f', 0, 64)
		entity := getEntity(player.EntityType)
		nameLabel := canvas.NewText(name, col)
		if k.selectedData.Winner() == team {
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
		}
		nameCont.Add(nameLabel)
		cont.Add(widget.NewLabel(kills))
		cont.Add(widget.NewLabel(deaths))
		cont.Add(widget.NewLabel(berries))
		cont.Add(widget.NewLabel(snail))

		cont.Add(widget.NewLabel(entity))
	}
	playerContainer := container.NewHBox(layout.NewSpacer(), nameCont, cont, layout.NewSpacer())
	headerLabel := widget.NewLabelWithStyle("Player Info", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	headerSeparator := widget.NewSeparator()
	wrapperContainer := container.NewVBox(headerLabel, headerSeparator, playerContainer)
	return wrapperContainer
}

func (k *KQBApp) ShowUploadWindow() {
	if k.u.BGLToken == "" {
		players := k.selectedData.Players()
		BGLPlayers := []string{"BGL 1", "BGL 2", "BGL 3", "BGL 4"}
		BGLTeams := []string{"BGL Team 1", "BGL Team 2"}
		BGLMatches := []string{"Match 1", "Match 2", "Match 3"}
		u := &Uploader{
			a:          k.a,
			w:          k.w,
			Players:    players,
			BGLPlayers: BGLPlayers,
			BGLTeams:   BGLTeams,
			PlayerMap:  make(map[string]string),
			TeamMap:    make(map[string]string),
			BGLMatches: BGLMatches,
			data:       k.selectedData,
			OnSuccess:  k.OnSetSuccess,
			OnFail:     k.OnSetFail,
		}

		uploadContainer := u.ShowUploadWindow()
		k.splitContainer.Objects[0] = uploadContainer
		k.u = u
	} else {
		k.u.Players = k.selectedData.Players()
		uploadContainer := k.u.ShowUploadWindow()
		k.splitContainer.Objects[0] = uploadContainer
		k.splitContainer.Objects[1] = layout.NewSpacer()
	}
	k.w.Resize(fyne.NewSize(900, 600))
}

func (k *KQBApp) OnSetSuccess() {
	if k.submission.ID == 0 {
		matchID := k.u.bgl.Matches[k.u.selectedMatch]
		k.submission = bgl.Result{
			ID:     matchID,
			Status: "Completed",
			Sets:   []bgl.Set{k.u.set},
		}
		k.submission.Sets[0].Number = 1
	} else {
		k.submission.Sets = append(k.submission.Sets, k.u.set)
		sLen := len(k.submission.Sets)
		k.submission.Sets[sLen-1].Number = sLen
	}
	k.subData = append(k.subData,
		bgl.SetMap{
			BGLMap: bgl.BGLMap{
				PlayerIDs:   k.u.GetPlayerMapByID(),
				TeamIDs:     k.u.GetTeamMapByID(),
				PlayerNames: k.u.PlayerMap,
				TeamNames:   k.u.TeamMap,
			},
			Raw: k.u.data,
		})
	k.splitContainer.Objects[0] = widget.NewLabelWithStyle("Select another set...", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	k.splitContainer.Objects[1] = k.ShowInputSets()
}

func (k *KQBApp) OnSetCompletion() {
	setCount := bgl.SetCount{
		Total: len(k.submission.Sets),
	}
	homeSets := 0
	for _, set := range k.submission.Sets {
		if set.Winner.ID == k.u.bgl.HomeID {
			homeSets++
		}
	}
	setCount.Away = setCount.Total - homeSets
	setCount.Home = homeSets

	if setCount.Home > setCount.Away {
		k.submission.Winner.ID = k.u.bgl.HomeID
		k.submission.Loser.ID = k.u.bgl.AwayID
	} else {
		k.submission.Loser.ID = k.u.bgl.HomeID
		k.submission.Winner.ID = k.u.bgl.AwayID
	}
	k.submission.SetCount = setCount

	loadingDiag := dialog.NewProgressInfinite("Match Results Upload", "Sending results to BGL", k.w)
	loadingDiag.Show()
	err := k.u.bgl.HandleMatchUpdate(k.submission)
	if err != nil {
		loadingDiag.Hide()
		dialog := dialog.NewInformation("Error", err.Error(), k.w)
		dialog.Show()
		return
	}

	finalOuput := bgl.FinalOutput{
		MatchID: k.u.bgl.Matches[k.u.selectedMatch],
		Sets:    k.subData,
	}
	k.u.bgl.SaveRawOutput(finalOuput)
	if err != nil {
		loadingDiag.Hide()
		dialog := dialog.NewInformation("Error", err.Error(), k.w)
		dialog.Show()
		return
	}
	loadingDiag.Hide()
	successDiag := dialog.NewInformation("Upload Success", "Match Upload Successful", k.w)
	successDiag.Show()
	k.ResetUploader()
}

func (k *KQBApp) OnSetFail() {
	k.splitContainer.Objects[0] = widget.NewLabelWithStyle("Select a set to continue...", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
}

func (k *KQBApp) ShowInputSets() *fyne.Container {
	base := container.NewVBox()
	if len(k.submission.Sets) > 0 {
		base.Add(widget.NewLabelWithStyle("Queued for Upload", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
		base.Add(widget.NewSeparator())
		cont := container.NewGridWithColumns(3,
			widget.NewLabel("Set"),
			widget.NewLabel("Winner"),
			widget.NewLabel("Loser"),
		)
		for idx, set := range k.submission.Sets {
			label := widget.NewLabel(strconv.Itoa(idx + 1))
			winner := widget.NewLabel(formatTeamName(set.Winner.Name))
			loser := widget.NewLabel(formatTeamName(set.Loser.Name))
			cont.Add(label)
			cont.Add(winner)
			cont.Add(loser)
		}

		base.Add(cont)
		if len(k.submission.Sets) >= 3 {
			uploadAction := widget.NewButton("Submit Match Results", func() {
				k.OnSetCompletion()
			})
			uploadAction.Importance = widget.HighImportance
			base.Add(uploadAction)
		}
		resetAction := widget.NewButton("Reset Upload Form", func() {
			k.ResetUploader()
		})
		base.Add(resetAction)
	}
	return base
}

func (k *KQBApp) BuildMapTable() *fyne.Container {
	mapLabel := widget.NewLabelWithStyle("Map", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	winConLabel := widget.NewLabelWithStyle("Win Con", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	winnerLabel := widget.NewLabelWithStyle("Winner", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	cont := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(3), mapLabel, winConLabel, winnerLabel)
	teamWinners := k.selectedData.TeamWinners()
	mapList := k.selectedData.MapList()

	for idx, winCon := range k.selectedData.WinCons() {
		var col color.RGBA
		team := teamWinners[idx]
		if team == "Blue" {
			col = blueColor
		} else {
			col = goldColor
		}
		mLabel := widget.NewLabel(mapList[idx])
		conLabel := widget.NewLabel(winCon)
		wonLabel := canvas.NewText(team, col)
		if k.selectedData.Winner() == team {
			wonLabel.TextStyle = fyne.TextStyle{Bold: true}
		}
		cont.Add(mLabel)
		cont.Add(conLabel)
		cont.Add(wonLabel)
	}
	mapWrapper := container.NewHBox(layout.NewSpacer(), cont, layout.NewSpacer())
	headerLabel := widget.NewLabelWithStyle("Map Details", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	headerSeparator := widget.NewSeparator()
	wrapperContainer := container.NewVBox(headerLabel, headerSeparator, mapWrapper)
	return wrapperContainer
}

func (k *KQBApp) ResetUploader() {
	for idx, _ := range k.splitContainer.Objects {
		k.splitContainer.Objects[idx] = layout.NewSpacer()
	}
	k.u = &Uploader{}
	k.submission = bgl.Result{}
	k.subData = []bgl.SetMap{}
	k.w.Resize(fyne.NewSize(600, 600))
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

func formatTeamName(input string) string {
	words := strings.Fields(input)
	output := ""
	for _, word := range words {
		output += strings.ToUpper(string(word[0]))
	}
	return output
}
