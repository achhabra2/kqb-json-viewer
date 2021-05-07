package bgl

// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    matchResult, err := UnmarshalMatchResult(bytes)
//    bytes, err = matchResult.Marshal()

import (
	"encoding/json"
	"time"

	"github.com/achhabra2/kqb-json-viewer/stats"
)

func UnmarshalMatchResult(data []byte) (MatchResult, error) {
	var r MatchResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MatchResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MatchResult struct {
	Count    int             `json:"count,omitempty"`
	Next     string          `json:"next,omitempty"`
	Previous interface{}     `json:"previous"`
	Results  []ResultElement `json:"results,omitempty"`
}

type ResultElement struct {
	ID               int           `json:"id,omitempty"`
	Home             TeamInfo      `json:"home,omitempty"`
	Away             TeamInfo      `json:"away,omitempty"`
	Circuit          Circuit       `json:"circuit,omitempty"`
	Round            Round         `json:"round,omitempty"`
	StartTime        interface{}   `json:"start_time"`
	TimeUntil        interface{}   `json:"time_until"`
	Scheduled        bool          `json:"scheduled,omitempty"`
	PrimaryCaster    PrimaryCaster `json:"primary_caster,omitempty"`
	SecondaryCasters []interface{} `json:"secondary_casters,omitempty"`
	Result           Result        `json:"result,omitempty"`
	VODLink          string        `json:"vod_link,omitempty"`
}

type Circuit struct {
	ID          int         `json:"id,omitempty"`
	Season      Season      `json:"season,omitempty"`
	Region      string      `json:"region,omitempty"`
	Tier        string      `json:"tier,omitempty"`
	Name        interface{} `json:"name"`
	VerboseName string      `json:"verbose_name,omitempty"`
}

type Season struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	League League `json:"league,omitempty"`
}

type League struct {
	Name string `json:"name,omitempty"`
	Href string `json:"_href,omitempty"`
}

type PrimaryCaster struct {
	Name       string `json:"name,omitempty"`
	BioLink    string `json:"bio_link,omitempty"`
	StreamLink string `json:"stream_link,omitempty"`
}

type Result struct {
	ID       int      `json:"match,omitempty"`
	Status   string   `json:"status,omitempty"`
	Winner   TeamInfo `json:"winner,omitempty"`
	Loser    TeamInfo `json:"loser,omitempty"`
	Sets     []Set    `json:"sets,omitempty"`
	SetCount SetCount `json:"set_count,omitempty"`
}

type TeamInfo struct {
	ID            int      `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	IsActive      bool     `json:"is_active,omitempty"`
	WINS          int      `json:"wins,omitempty"`
	Losses        int      `json:"losses,omitempty"`
	CircuitAbbrev string   `json:"circuit_abbrev,omitempty"`
	Members       []Member `json:"members,omitempty"`
}

type SetCount struct {
	Home  int `json:"home,omitempty"`
	Away  int `json:"away,omitempty"`
	Total int `json:"total,omitempty"`
}

type Set struct {
	Number int      `json:"number,omitempty"`
	Winner TeamInfo `json:"winner,omitempty"`
	Loser  TeamInfo `json:"loser,omitempty"`
}

type Round struct {
	Number string `json:"number,omitempty"`
	Name   string `json:"name,omitempty"`
}

type Member struct {
	ID              int64    `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	NamePhonetic    string   `json:"name_phonetic,omitempty"`
	Pronouns        Pronouns `json:"pronouns,omitempty"`
	DiscordUsername string   `json:"discord_username,omitempty"`
	TwitchUsername  string   `json:"twitch_username,omitempty"`
	Bio             string   `json:"bio,omitempty"`
	Emoji           string   `json:"emoji,omitempty"`
	AvatarURL       string   `json:"avatar_url,omitempty"`
	Modified        string   `json:"modified,omitempty"`
	Created         string   `json:"created,omitempty"`
}

type Pronouns string

const (
	HeHim    Pronouns = "he/him"
	SheHer   Pronouns = "she/her"
	TheyThem Pronouns = "they/them"
)

type ResultSubmission struct {
	Match         int                   `json:"match,omitempty"`
	Status        string                `json:"status,omitempty"`
	Winner        int                   `json:"winner,omitempty"`
	Loser         int                   `json:"loser,omitempty"`
	Sets          []ResultSubmissionSet `json:"sets,omitempty"`
	PlayerMapping []PlayerObject        `json:"player_mappings,omitempty"`
	TeamMapping   []TeamObject          `json:"team_mappings,omitempty"`
	Source        string                `json:"source"`
	Notes         string                `json:"notes"`
}

type ResultSubmissionSet struct {
	Number    int          `json:"number,omitempty"`
	Winner    int          `json:"winner,omitempty"`
	Loser     int          `json:"loser,omitempty"`
	SetLog    ResultSetLog `json:"log,omitempty"`
	TimeStamp time.Time
}

type ResultSetLog struct {
	FileName string          `json:"filename,omitempty"`
	Body     stats.StatsJSON `json:"body,omitempty"`
}

type PlayerObject struct {
	ID   int    `json:"player,omitempty"`
	Name string `json:"nickname,omitempty"`
}

// type TeamObject struct {
// 	Color string `json:"color,omitempty"`
// 	Team  int    `json:"team,omitempty"`
// }

type TeamObject struct {
	Team int `json:"color,omitempty"`
	ID   int `json:"team,omitempty"`
}

func BglMapToObjects(b BGLMap) ([]PlayerObject, []TeamObject) {
	playerMapping := make([]PlayerObject, 0)
	teamMapping := make([]TeamObject, 0)

	for k, v := range b.PlayerIDs {
		playerObject := PlayerObject{
			Name: k,
			ID:   v,
		}
		playerMapping = append(playerMapping, playerObject)
	}

	for k, v := range b.TeamIDs {
		var team int
		switch k {
		case "Gold":
			team = 1
		case "Blue":
			team = 2
		}
		teamObject := TeamObject{
			Team: team,
			ID:   v,
		}
		teamMapping = append(teamMapping, teamObject)
	}
	return playerMapping, teamMapping
}
