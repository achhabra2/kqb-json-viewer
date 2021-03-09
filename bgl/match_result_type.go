package bgl

// MatchResult This should be the data returned from
// https://kqb.buzz/api/matches/
type MatchResult struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		ID   int `json:"id"`
		Home struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			IsActive bool        `json:"is_active"`
			Wins     int         `json:"wins"`
			Losses   int         `json:"losses"`
			Members  []string    `json:"members"`
			Dynasty  interface{} `json:"dynasty"`
		} `json:"home"`
		Away struct {
			ID       int         `json:"id"`
			Name     string      `json:"name"`
			IsActive bool        `json:"is_active"`
			Wins     int         `json:"wins"`
			Losses   int         `json:"losses"`
			Members  []string    `json:"members"`
			Dynasty  interface{} `json:"dynasty"`
		} `json:"away"`
		Circuit struct {
			ID     int `json:"id"`
			Season struct {
				ID     int    `json:"id"`
				Name   string `json:"name"`
				League struct {
					Name string `json:"name"`
					Href string `json:"_href"`
				} `json:"league"`
			} `json:"season"`
			Region      string      `json:"region"`
			Tier        string      `json:"tier"`
			Name        interface{} `json:"name"`
			VerboseName string      `json:"verbose_name"`
		} `json:"circuit"`
		Round struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"round"`
		StartTime     interface{} `json:"start_time"`
		TimeUntil     interface{} `json:"time_until"`
		Scheduled     bool        `json:"scheduled"`
		PrimaryCaster struct {
			Name       string `json:"name"`
			BioLink    string `json:"bio_link"`
			StreamLink string `json:"stream_link"`
		} `json:"primary_caster"`
		SecondaryCasters []interface{} `json:"secondary_casters"`
		Result           struct {
			Status string `json:"status"`
			Winner struct {
				ID            int    `json:"id"`
				Name          string `json:"name"`
				IsActive      bool   `json:"is_active"`
				Wins          int    `json:"wins"`
				Losses        int    `json:"losses"`
				CircuitAbbrev string `json:"circuit_abbrev"`
			} `json:"winner"`
			Loser struct {
				ID            int    `json:"id"`
				Name          string `json:"name"`
				IsActive      bool   `json:"is_active"`
				Wins          int    `json:"wins"`
				Losses        int    `json:"losses"`
				CircuitAbbrev string `json:"circuit_abbrev"`
			} `json:"loser"`
			Sets []struct {
				Number int `json:"number"`
				Winner struct {
					ID            int    `json:"id"`
					Name          string `json:"name"`
					IsActive      bool   `json:"is_active"`
					Wins          int    `json:"wins"`
					Losses        int    `json:"losses"`
					CircuitAbbrev string `json:"circuit_abbrev"`
				} `json:"winner"`
				Loser struct {
					ID            int    `json:"id"`
					Name          string `json:"name"`
					IsActive      bool   `json:"is_active"`
					Wins          int    `json:"wins"`
					Losses        int    `json:"losses"`
					CircuitAbbrev string `json:"circuit_abbrev"`
				} `json:"loser"`
			} `json:"sets"`
			SetCount struct {
				Home  int `json:"home"`
				Away  int `json:"away"`
				Total int `json:"total"`
			} `json:"set_count"`
		} `json:"result"`
		VodLink string `json:"vod_link"`
	} `json:"results"`
}

// ResultSubmission is the data we are sending back to
// https://kqb.buzz/api/matches/
type ResultSubmission struct {
	ID     int `json:"ID"`
	Result struct {
		Status string `json:"status"`
		Winner struct {
			ID int `json:"id"`
		} `json:"winner"`
		Loser struct {
			ID int `json:"id"`
		} `json:"loser"`
		Sets []struct {
			Number int `json:"number"`
			Winner struct {
				ID int `json:"id"`
			} `json:"winner"`
			Loser struct {
				ID int `json:"id"`
			} `json:"loser"`
		} `json:"sets"`
		SetCount struct {
			Home  int `json:"home"`
			Away  int `json:"away"`
			Total int `json:"total"`
		} `json:"set_count"`
	} `json:"result"`
}
