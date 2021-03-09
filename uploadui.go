package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/achhabra2/kqb-json-viewer/bgl"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

type Uploader struct {
	BGLToken      string
	a             fyne.App
	w             fyne.Window
	c             *fyne.Container
	data          stats.StatsJSON
	BGLPlayers    []string
	Players       []string
	BGLTeams      []string
	BGLMatches    []string
	PlayerMap     map[string]string
	TeamMap       map[string]string
	bgl           bgl.BGLData
	selectedMatch string
}

func (u *Uploader) ShowUploadWindow() {
	header := widget.NewLabelWithStyle("BGL Stats Uploader", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	separator := widget.NewSeparator()
	cont := container.NewVBox(header, separator)
	u.c = cont
	tokenForm := u.BuildTokenForm()
	cont.Add(tokenForm)
	cont.Add(layout.NewSpacer())
	u.w.SetContent(u.c)
	u.w.Resize(fyne.NewSize(500, 500))
	u.w.CenterOnScreen()
	u.w.Show()
}

func (u *Uploader) BuildTokenForm() *widget.Form {
	entry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "BGL Token", Widget: entry}},
		OnSubmit: func() {
			u.BGLToken = entry.Text
			if u.BGLToken == "" {
				errorDialog := dialog.NewInformation("Input Error", "Please enter a valid Token", u.w)
				errorDialog.Show()
				return
			} else {
				u.bgl = bgl.BGLData{Token: u.BGLToken}
				u.ShowLoadingIndicator()
				if u.IsValidToken() {
					u.bgl.LoadCurrentMatches()
					u.BGLMatches = u.bgl.GetMatchNames()
					u.c.Objects[2] = u.BuildMatchForm()
				} else {
					u.c.Objects[2] = u.BuildTokenForm()
					errorDialog := dialog.NewInformation("Token Error", "Invalid Token, Please try again", u.w)
					errorDialog.Show()
					return
				}
			}
		},
		OnCancel: func() {
			u.w.Close()
		},
	}

	return form
}

func (u *Uploader) BuildTeamForm() *widget.Form {
	formItems := make([]*widget.FormItem, 2)
	// teamLabels := []string{"Gold", "Blue"}
	// for idx, name := range teamLabels {
	// combo := widget.NewSelect(u.BGLTeams, func(value string) {
	// 	u.TeamMap[teamLabels[idx]] = value
	// })
	// 	item := widget.NewFormItem(name, combo)
	// 	formItems[idx] = item
	// }

	blueCombo := widget.NewSelect(u.BGLTeams, func(value string) {
		u.TeamMap["Blue"] = value
	})
	formItems[1] = widget.NewFormItem("Blue Team", blueCombo)
	goldCombo := widget.NewSelect(u.BGLTeams, func(value string) {
		u.TeamMap["Gold"] = value
	})
	formItems[0] = widget.NewFormItem("Gold Team", goldCombo)

	form := &widget.Form{
		Items: formItems,
		OnSubmit: func() {
			if u.ValidateParams() {
				errorDialog := dialog.NewInformation("Input Error", "Duplicate entries found, please correct the information and try again. ", u.w)
				errorDialog.Show()
			} else {
				log.Println(u.TeamMap)
				u.bgl.LoadPlayersForMatch(u.selectedMatch)
				u.BGLPlayers = u.bgl.GetPlayerNames()
				u.c.Objects[2] = u.BuildPlayerForm()
			}
		},
		OnCancel: func() {
			u.w.Close()
		},
	}

	return form
}

func (u *Uploader) BuildPlayerForm() *widget.Form {

	formItems := make([]*widget.FormItem, 8)
	for idx, name := range u.Players {
		fmt.Println(idx, name)
		combo := widget.NewSelect(u.BGLPlayers, u.playerCallback(name))
		item := widget.NewFormItem(name, combo)
		formItems[idx] = item
	}
	form := &widget.Form{
		Items: formItems,
		OnSubmit: func() {
			fmt.Println(u.PlayerMap)
			if u.ValidateParams() {
				errorDialog := dialog.NewInformation("Input Error", "Duplicate entries found, please correct the information and try again. ", u.w)
				errorDialog.Show()
			} else {
				u.HandleUpload()
			}
		},
		OnCancel: func() {
			u.w.Close()
		},
	}

	return form
}

func (u *Uploader) playerCallback(name string) func(string) {
	return func(value string) {
		u.PlayerMap[name] = value
	}
}

func (u *Uploader) BuildMatchForm() *widget.Form {
	combo := widget.NewSelect(u.BGLMatches, func(value string) {
		u.selectedMatch = value
	})
	item := widget.NewFormItem("Select Match", combo)
	form := &widget.Form{
		Items: []*widget.FormItem{item},
		OnSubmit: func() {
			if u.selectedMatch == "" {
				errorDialog := dialog.NewInformation("Input Error", "Please select a valid match ", u.w)
				errorDialog.Show()
			} else {
				u.bgl.LoadTeamsForMatch(u.selectedMatch)
				u.BGLTeams = u.bgl.GetTeamNames()
				u.c.Objects[2] = u.BuildTeamForm()
			}
		},
		OnCancel: func() {
			u.w.Close()
		},
	}
	return form
}

func (u *Uploader) ShowLoadingIndicator() {
	progress := widget.NewProgressBarInfinite()
	u.c.Objects[2] = progress
}

func ShowUploadWindow(a fyne.App, data stats.StatsJSON) {
	window := a.NewWindow("BGL Uploader")
	players := data.Players()
	BGLPlayers := []string{"BGL 1", "BGL 2", "BGL 3", "BGL 4"}
	BGLTeams := []string{"BGL Team 1", "BGL Team 2"}
	BGLMatches := []string{"Match 1", "Match 2", "Match 3"}
	u := &Uploader{
		a:          a,
		w:          window,
		Players:    players,
		BGLPlayers: BGLPlayers,
		BGLTeams:   BGLTeams,
		PlayerMap:  make(map[string]string),
		TeamMap:    make(map[string]string),
		BGLMatches: BGLMatches,
		data:       data,
	}
	u.ShowUploadWindow()
}

func (u *Uploader) ValidateParams() bool {
	selectedPlayers := make([]string, 0)
	selectedTeams := make([]string, 0)

	for _, name := range u.PlayerMap {
		selectedPlayers = append(selectedPlayers, name)
	}

	for _, team := range u.TeamMap {
		selectedTeams = append(selectedTeams, team)
	}

	teamDuplicates := findDuplicates(selectedTeams)
	playerDuplicates := findDuplicates(selectedPlayers)
	if teamDuplicates || playerDuplicates {
		return true
	} else {
		return false
	}
}

func (u *Uploader) HandleUpload() {

}

func (u *Uploader) IsValidToken() bool {
	err := u.bgl.GetMe()
	if err != nil {
		return false
	} else {
		return true
	}
}

func findDuplicates(array []string) bool {
	matchFound := false
	for _, needle := range array {
		matches := 0
		for _, haystack := range array {
			if needle == haystack {
				matches++
			}
		}
		if matches > 1 {
			matchFound = true
			break
		}
	}
	return matchFound
}
