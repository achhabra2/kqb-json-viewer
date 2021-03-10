package bgl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/achhabra2/kqb-json-viewer/stats"
)

type BGLData struct {
	Token       string
	Matches     map[string]int
	Teams       map[string]int
	Players     map[string]int
	HomeID      int
	AwayID      int
	HomeName    string
	AwayName    string
	matchResult MatchResult
}

func (b *BGLData) LoadCurrentMatches() {
	dir, _ := os.Getwd()
	file := filepath.Join(dir, "/fixtures/match_result.json")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Could not read json file", err)
	}
	err = json.Unmarshal(data, &b.matchResult)
	if err != nil {
		log.Fatal("Could not parse json file", err)
	}

	b.Matches = make(map[string]int)
	for _, result := range b.matchResult.Results {
		key := result.Away.Name + " @ " + result.Home.Name
		b.Matches[key] = result.ID
	}
}

func (b *BGLData) GetMatchNames() []string {
	matches := make([]string, 0)
	for key, _ := range b.Matches {
		matches = append(matches, key)
	}
	return matches
}

func (b *BGLData) LoadPlayersForMatch(match string) {
	matchID := b.Matches[match]
	players := make(map[string]int)
	for _, result := range b.matchResult.Results {
		if result.ID == matchID {
			for idx, playerName := range result.Home.Members {
				players[playerName] = idx
			}
			for idx, playerName := range result.Away.Members {
				players[playerName] = idx
			}
		}
	}
	b.Players = players
}

func (b *BGLData) GetPlayerNames() []string {
	players := make([]string, 0)
	for playerName, _ := range b.Players {
		players = append(players, playerName)
	}
	return players
}

func (b *BGLData) LoadTeamsForMatch(match string) {
	matchID := b.Matches[match]
	teams := make(map[string]int)
	for _, result := range b.matchResult.Results {
		if result.ID == matchID {
			teams[result.Away.Name] = result.Away.ID
			teams[result.Home.Name] = result.Home.ID
			b.HomeID = result.Home.ID
			b.HomeName = result.Home.Name
			b.AwayID = result.Away.ID
			b.AwayName = result.Away.Name
		}
	}
	b.Teams = teams
}

func (b *BGLData) GetTeamNames() []string {
	teams := make([]string, 0)
	for name, _ := range b.Teams {
		teams = append(teams, name)
	}
	return teams
}

func (b *BGLData) GetMe() error {
	// url := "https://kqb.buzz/api/me/?format=json"
	// method := "GET"

	// client := &http.Client{}
	// req, err := http.NewRequest(method, url, nil)

	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	// req.Header.Add("Authorization", "Token "+b.Token)

	// res, err := client.Do(req)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// if res.StatusCode != 200 {
	// 	return errors.New("Invalid status code")
	// }
	// defer res.Body.Close()

	return nil
}

func (b *BGLData) HandleMatchUpdate(result Result) error {
	wd, _ := os.Getwd()
	outPath := filepath.Join(wd, "/tmp/match_update.json")
	output, err := json.MarshalIndent(result, "  ", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outPath, output, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (b *BGLData) SaveRawOutput(final FinalOutput) error {
	wd, _ := os.Getwd()
	outPath := filepath.Join(wd, "/tmp/match_output.json")
	output, err := json.MarshalIndent(final, "  ", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outPath, output, 0644)
	if err != nil {
		return err
	}
	return nil
}

type BGLMap struct {
	PlayerNames map[string]string `json:"player_names"`
	TeamNames   map[string]string `json:"team_names"`
	PlayerIDs   map[string]int    `json:"player_ids"`
	TeamIDs     map[string]int    `json:"team_ids"`
}

type SetMap struct {
	Raw    stats.StatsJSON `json:"raw"`
	BGLMap BGLMap          `json:"bgl_map"`
}

type FinalOutput struct {
	MatchID int      `json:"match_id"`
	Sets    []SetMap `json:"sets"`
}
