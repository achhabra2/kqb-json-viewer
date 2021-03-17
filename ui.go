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
	"github.com/achhabra2/kqb-json-viewer/icons"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

// Colors for Gold and Blue Labels
var goldColor = color.RGBA{255, 179, 0, 200}
var blueColor = color.RGBA{43, 93, 255, 200}

// Main App Struct
type KQBApp struct {
	files          []string
	selectedData   stats.StatsJSON
	a              fyne.App
	w              fyne.Window
	mainContainer  *container.Split
	splitContainer *fyne.Container
	u              *Uploader
	submission     bgl.Result
	subData        []stats.SetResult
	bglMap         bgl.BGLMap
	selectedFiles  map[string]int
	selectedFile   string
	fileDropDown   *widget.Select
}

// Main function to perform app setup and show the main window
func (k *KQBApp) ShowMainWindow() {
	k.u = &Uploader{}

	// Initialize the header components
	// timeWidget is the timestamp label for the current match
	timeWidget := widget.NewLabel(getTimeString(k.files[0]))

	// Within the time widget we are also going to show if this file is selected for upload
	checkIconWidget := getStatLogo("Check")
	selectedWidget := container.NewCenter(checkIconWidget)
	k.selectedFiles = make(map[string]int)
	matchLabelWidget := widget.NewLabel("Match: ")
	timeContainer := container.NewHBox(layout.NewSpacer(), matchLabelWidget, timeWidget, selectedWidget, layout.NewSpacer())

	about := fyne.NewMenuItem("About", func() {
		aboutMessage := fmt.Sprintf("kqb-json-viewer version %s \n by Prosive", version)
		dialog := dialog.NewInformation("About", aboutMessage, k.w)
		dialog.Show()
	})
	openDirectory := fyne.NewMenuItem("Open Stats Folder", func() {
		stats.OpenStatDirectory()
	})

	upload := fyne.NewMenuItem("Add Set to Match", func() {
		k.ShowUploadWindow()
	})

	fileMenu := fyne.NewMenu("File", about, openDirectory)
	bglMenu := fyne.NewMenu("BGL", upload)
	mainMenu := fyne.NewMainMenu(fileMenu, bglMenu)

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
		k.selectedFile = value
		if k.u.BGLToken != "" {
			k.u.data = k.selectedData
		}

		if k.selectedFiles[value] == 1 {
			selectedWidget.Show()
		} else {
			selectedWidget.Hide()
		}
		timeWidget.Text = getTimeString(trimmedMap[value])
		cont.Hide()
		cont.Objects[2] = k.BuildPlayerUI()
		cont.Objects[3] = k.BuildMapTable()
		cont.Show()
	})

	k.fileDropDown = combo

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
	//	cont.Add(upload)

	combo.SetSelectedIndex(0)
	trailingContainer := container.NewVBox(layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer())
	vSplitContainer := container.NewVScroll(trailingContainer)
	vSplitLeft := container.NewVScroll(cont)
	mainContainer := container.NewHSplit(vSplitLeft, vSplitContainer)
	k.mainContainer = mainContainer
	k.splitContainer = trailingContainer
	k.w.SetContent(k.mainContainer)
	k.w.SetMainMenu(mainMenu)
	k.w.SetPadded(true)
	k.w.Resize(fyne.NewSize(500, 900))

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

	nameCont := container.New(layout.NewGridLayoutWithColumns(1))
	cont := container.New(layout.NewGridLayoutWithColumns(5))
	nameCont.Add(widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	// cont.Add(widget.NewLabelWithStyle("Kills", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	// cont.Add(widget.NewLabelWithStyle("Deaths", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	// cont.Add(widget.NewLabelWithStyle("Berries", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	// cont.Add(widget.NewLabelWithStyle("Snail", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))

	cont.Add(getStatLogo("Kills"))
	cont.Add(getStatLogo("Deaths"))
	cont.Add(getStatLogo("Berries"))
	cont.Add(getStatLogo("Snail"))

	cont.Add(widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
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
		// entity := getEntity(player.EntityType)
		nameLabel := canvas.NewText(name, col)
		if k.selectedData.Winner() == team {
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
		}
		nameCont.Add(nameLabel)
		cont.Add(widget.NewLabelWithStyle(kills, fyne.TextAlignCenter, fyne.TextStyle{}))
		cont.Add(widget.NewLabelWithStyle(deaths, fyne.TextAlignCenter, fyne.TextStyle{}))
		cont.Add(widget.NewLabelWithStyle(berries, fyne.TextAlignCenter, fyne.TextStyle{}))
		cont.Add(widget.NewLabelWithStyle(snail, fyne.TextAlignCenter, fyne.TextStyle{}))

		// cont.Add(widget.NewLabelWithStyle(entity, fyne.TextAlignCenter, fyne.TextStyle{}))
		cont.Add(getEntityLogo(player.EntityType))
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

		u := &Uploader{
			a:          k.a,
			w:          k.w,
			Players:    players,
			BGLPlayers: make([]string, 0),
			BGLTeams:   make([]string, 0),
			PlayerMap:  make(map[string]string),
			TeamMap:    make(map[string]string),
			BGLMatches: make([]string, 0),
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
	k.w.Resize(fyne.NewSize(900, 900))
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
	k.bglMap = bgl.BGLMap{
		PlayerIDs:   k.u.GetPlayerMapByID(),
		TeamIDs:     k.u.GetTeamMapByID(),
		PlayerNames: k.u.PlayerMap,
		TeamNames:   k.u.TeamMap,
	}
	k.subData = append(k.subData, k.u.data.GetSetResult())
	k.selectedFiles[k.selectedFile] = 1
	k.fileDropDown.SetSelected(k.selectedFile)
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

	loadingWidget := widget.NewProgressBarInfinite()
	loadingDiag := dialog.NewCustom("Match Results Upload", "", loadingWidget, k.w)
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
		// Sets:    k.subData,
		BGLMap: k.bglMap,
		Result: stats.GetMatchResult(k.subData...),
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

	cont := container.New(layout.NewGridLayoutWithColumns(3), mapLabel, winConLabel, winnerLabel)
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
		mLabel := widget.NewLabelWithStyle(mapList[idx], fyne.TextAlignCenter, fyne.TextStyle{})
		conLabel := widget.NewLabelWithStyle(winCon, fyne.TextAlignCenter, fyne.TextStyle{})
		wonLabel := canvas.NewText(team, col)
		wonLabel.Alignment = fyne.TextAlignCenter

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
	// k.subData = []bgl.SetMap{}
	k.subData = []stats.SetResult{}
	k.selectedFiles = make(map[string]int)
	k.w.Resize(fyne.NewSize(500, 900))
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

func getEntityLogo(entity int) *fyne.Container {
	var res *fyne.StaticResource
	switch entity {
	case 3:
		res = fyne.NewStaticResource("Queen_Logo_Large.png", icons.Queen_Logo)
	default:
		res = fyne.NewStaticResource("Worker_Logo_large.png", icons.Worker_Logo)
	}
	iconCanvas := canvas.NewImageFromResource(res)
	iconCanvas.FillMode = canvas.ImageFillContain
	iconCanvas.SetMinSize(fyne.NewSize(36, 36))
	cont := container.NewPadded(iconCanvas)
	return cont
}

func getStatLogo(stat string) *fyne.Container {
	var res *fyne.StaticResource
	switch stat {
	case "Kills":
		res = fyne.NewStaticResource("Icon_Kills.png", icons.Icon_Kills)
	case "Deaths":
		res = fyne.NewStaticResource("Icon_Deaths.png", icons.Icon_Deaths)
	case "Berries":
		res = fyne.NewStaticResource("Icon_Berries.png", icons.Icon_Berries)
	case "Snail":
		res = fyne.NewStaticResource("Icon_Snail.png", icons.Icon_Snail)
	case "Check":
		res = fyne.NewStaticResource("check.png", icons.Check)
	}
	iconCanvas := canvas.NewImageFromResource(res)
	iconCanvas.FillMode = canvas.ImageFillContain
	iconCanvas.SetMinSize(fyne.NewSize(16, 16))
	cont := container.NewPadded(iconCanvas)
	//cont.Resize(fyne.NewSize(16, 16))
	return cont
}

func formatTeamName(input string) string {
	words := strings.Fields(input)
	output := ""
	for _, word := range words {
		output += strings.ToUpper(string(word[0]))
	}
	return output
}
