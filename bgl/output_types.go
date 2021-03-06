package bgl

import "github.com/achhabra2/kqb-json-viewer/stats"

type BGLMap struct {
	PlayerNames map[string]string `json:"player_names"`
	TeamNames   map[string]string `json:"team_names"`
	PlayerIDs   map[string]int    `json:"player_ids"`
	TeamIDs     map[string]int    `json:"team_ids"`
}

type SetMap struct {
	Raw   stats.StatsJSON `json:"raw"`
	Stats stats.SetResult `json:"set_result"`
}

type FinalOutput struct {
	MatchID int `json:"match_id"`
	//	Sets    []SetMap `json:"sets"`
	Result stats.MatchResult `json:"match_result"`
	BGLMap BGLMap            `json:"bgl_map"`
}
