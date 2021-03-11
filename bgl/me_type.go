package bgl

import "time"

type User struct {
	ID         int       `json:"id"`
	IsActive   bool      `json:"is_active"`
	FirstName  string    `json:"first_name"`
	DateJoined time.Time `json:"date_joined"`
	LastLogin  time.Time `json:"last_login"`
	Player     struct {
		ID              int       `json:"id"`
		Name            string    `json:"name"`
		NamePhonetic    string    `json:"name_phonetic"`
		Pronouns        string    `json:"pronouns"`
		DiscordUsername string    `json:"discord_username"`
		TwitchUsername  string    `json:"twitch_username"`
		Bio             string    `json:"bio"`
		Emoji           string    `json:"emoji"`
		AvatarURL       string    `json:"avatar_url"`
		Modified        time.Time `json:"modified"`
		Created         time.Time `json:"created"`
		Teams           []struct {
			ID         int       `json:"id"`
			Name       string    `json:"name"`
			InviteCode string    `json:"invite_code"`
			Modified   time.Time `json:"modified"`
			Created    time.Time `json:"created"`
			Circuit    struct {
				ID     int    `json:"id"`
				Region string `json:"region"`
				Tier   string `json:"tier"`
				Name   string `json:"name"`
				Season int    `json:"season"`
			} `json:"circuit"`
			Captain struct {
				ID              int         `json:"id"`
				Name            string      `json:"name"`
				NamePhonetic    string      `json:"name_phonetic"`
				Pronouns        string      `json:"pronouns"`
				DiscordUsername string      `json:"discord_username"`
				AvatarURL       string      `json:"avatar_url"`
				TwitchUsername  string      `json:"twitch_username"`
				Bio             string      `json:"bio"`
				Emoji           string      `json:"emoji"`
				Modified        time.Time   `json:"modified"`
				Created         time.Time   `json:"created"`
				User            interface{} `json:"user"`
			} `json:"captain"`
			Dynasty interface{} `json:"dynasty"`
			Members []struct {
				ID              int         `json:"id"`
				Name            string      `json:"name"`
				NamePhonetic    string      `json:"name_phonetic"`
				Pronouns        string      `json:"pronouns"`
				DiscordUsername string      `json:"discord_username"`
				AvatarURL       string      `json:"avatar_url"`
				TwitchUsername  string      `json:"twitch_username"`
				Bio             string      `json:"bio"`
				Emoji           string      `json:"emoji"`
				Modified        time.Time   `json:"modified"`
				Created         time.Time   `json:"created"`
				User            interface{} `json:"user"`
			} `json:"members"`
		} `json:"teams"`
		AwardSummary []interface{} `json:"award_summary"`
	} `json:"player"`
}
