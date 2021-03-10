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
	OnSuccess     func()
	OnFail        func()
	set           bgl.Set
}

func (u *Uploader) ShowUploadWindow() *fyne.Container {
	header := widget.NewLabelWithStyle("BGL Stats Uploader", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	separator := widget.NewSeparator()
	cont := container.NewVBox(header, separator)
	u.c = cont
	if u.BGLToken == "" {
		tokenForm := u.BuildTokenForm()
		cont.Add(tokenForm)
	} else if u.selectedMatch == "" {
		matchForm := u.BuildMatchForm()
		cont.Add(matchForm)
	} else if u.TeamMap["Blue"] == "" || u.TeamMap["Gold"] == "" {
		teamForm := u.BuildTeamForm()
		cont.Add(teamForm)
	} else {
		u.PlayerMap = make(map[string]string)
		playerForm := u.BuildPlayerForm()
		cont.Add(playerForm)
	}
	cont.Add(layout.NewSpacer())
	return cont
}

func (u *Uploader) BuildTokenForm() *widget.Form {
	entry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "BGL Token", Widget: entry}},
		OnSubmit: func() {
			u.BGLToken = entry.Text
			if u.BGLToken == "" {
				errorDialog := dialog.NewInformation("Token Input Error", "Please enter a valid Token", u.w)
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
					errorDialog := dialog.NewInformation("Token Validation Error", "Invalid Token, Please try again", u.w)
					errorDialog.Show()
					return
				}
			}
		},
		OnCancel: func() {
			u.OnFail()
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
			} else if !u.IsTeamFormFilled() {
				errorDialog := dialog.NewInformation("Input Error", "Make a selection for both teams. ", u.w)
				errorDialog.Show()
			} else {
				log.Println(u.TeamMap)
				u.bgl.LoadPlayersForMatch(u.selectedMatch)
				u.BGLPlayers = u.bgl.GetPlayerNames()
				u.c.Objects[2] = u.BuildPlayerForm()
			}
		},
		OnCancel: func() {
			u.OnFail()
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
			} else if !u.IsPlayerFormFilled() {
				errorDialog := dialog.NewInformation("Input Error", "Make a selection for all players. ", u.w)
				errorDialog.Show()
			} else {
				u.HandleSubmit()
			}
		},
		OnCancel: func() {
			u.OnFail()
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
			u.OnFail()
		},
	}
	return form
}

func (u *Uploader) ShowLoadingIndicator() {
	progress := widget.NewProgressBarInfinite()
	u.c.Objects[2] = progress
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

func (u *Uploader) IsPlayerFormFilled() bool {
	if len(u.PlayerMap) < 8 {
		return false
	}
	return true
}

func (u *Uploader) IsTeamFormFilled() bool {
	if len(u.TeamMap) < 2 {
		return false
	}
	return true
}

func (u *Uploader) HandleSubmit() {
	// TODO - Come up with final set JSON
	// matchID := u.bgl.Matches[u.selectedMatch]
	goldTeamName := u.TeamMap["Gold"]
	blueTeamName := u.TeamMap["Blue"]
	goldTeamID := u.bgl.Teams[goldTeamName]
	blueTeamID := u.bgl.Teams[blueTeamName]

	winner := bgl.TeamInfo{}
	loser := bgl.TeamInfo{}

	switch u.data.Winner() {
	case "Blue":
		winner.ID = blueTeamID
		winner.Name = blueTeamName
		loser.ID = goldTeamID
		loser.Name = goldTeamName
	case "Gold":
		winner.ID = goldTeamID
		winner.Name = goldTeamName
		loser.ID = blueTeamID
		loser.Name = blueTeamName
	default:
		break
	}

	submissionSet := bgl.Set{
		Winner: winner,
		Loser:  loser,
	}
	u.set = submissionSet
	u.OnSuccess()
}

func (u *Uploader) IsValidToken() bool {
	err := u.bgl.GetMe()
	if err != nil {
		return false
	} else {
		return true
	}
}

func (u *Uploader) GetPlayerMapByID() map[string]int {
	output := make(map[string]int)
	for key, val := range u.PlayerMap {
		output[key] = u.bgl.Players[val]
	}
	return output
}

func (u *Uploader) GetTeamMapByID() map[string]int {
	output := make(map[string]int)
	for key, val := range u.TeamMap {
		output[key] = u.bgl.Teams[val]
	}
	return output
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
