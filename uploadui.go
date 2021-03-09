package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/achhabra2/kqb-json-viewer/bgl"
)

type Uploader struct {
	BGLToken      string
	a             fyne.App
	w             fyne.Window
	c             *fyne.Container
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
	cont := container.NewVBox()
	u.c = cont
	tokenForm := u.BuildTokenForm()
	cont.Add(tokenForm)
	u.w.SetContent(u.c)
	cont.Resize(fyne.NewSize(300, 300))
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
			u.bgl = bgl.BGLData{Token: "ABCDEF"}
			u.ShowLoadingIndicator()
			u.bgl.LoadCurrentMatches()
			u.BGLMatches = u.bgl.GetMatchNames()
			u.c.Objects[0] = u.BuildMatchForm()
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
			fmt.Println(u.TeamMap)
			u.bgl.LoadPlayersForMatch(u.selectedMatch)
			u.BGLPlayers = u.bgl.GetPlayerNames()
			u.c.Objects[0] = u.BuildPlayerForm()
		},
		OnCancel: func() {
			u.w.Close()
		},
	}

	return form
}

func (u *Uploader) BuildPlayerForm() *widget.Form {

	formItems := make([]*widget.FormItem, 4)
	for idx, name := range u.Players {
		combo := widget.NewSelect(u.BGLPlayers, func(value string) {
			u.PlayerMap[name] = value
		})
		item := widget.NewFormItem(name, combo)
		formItems[idx] = item
	}
	form := &widget.Form{
		Items: formItems,
		OnSubmit: func() {
			fmt.Println(u.PlayerMap)
		},
		OnCancel: func() {
			u.w.Close()
		},
	}

	return form
}

func (u *Uploader) BuildMatchForm() *widget.Form {
	combo := widget.NewSelect(u.BGLMatches, func(value string) {
		u.selectedMatch = value
	})
	item := widget.NewFormItem("Select Match", combo)
	form := &widget.Form{
		Items: []*widget.FormItem{item},
		OnSubmit: func() {
			u.bgl.LoadTeamsForMatch(u.selectedMatch)
			u.BGLTeams = u.bgl.GetTeamNames()
			u.c.Objects[0] = u.BuildTeamForm()
		},
		OnCancel: func() {
			u.w.Close()
		},
	}
	return form
}

func (u *Uploader) ShowLoadingIndicator() {
	progress := widget.NewProgressBarInfinite()
	u.c.Objects[0] = progress
}

func ShowUploadWindow(a fyne.App) {
	window := a.NewWindow("BGL Uploader")
	// playerMap := make(map[string]string)
	players := []string{"Player 1", "Player 2", "Player 3", "Player 4"}
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
	}
	u.ShowUploadWindow()
}
