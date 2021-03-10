package bgl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type BGLData struct {
	Token       string
	Matches     map[string]int
	Teams       map[string]int
	Players     map[string]int
	HomeID      int
	AwayID      int
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
			b.AwayID = result.Away.ID
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
