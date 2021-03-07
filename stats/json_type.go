package stats

type StatsJSON struct {
	PlayerMatchStats []struct {
		Nickname               string  `json:"nickname"`
		InputID                int     `json:"inputID"`
		ActorNr                int     `json:"actorNr"`
		PlayerID               string  `json:"playerId"`
		PlayerIndex            int     `json:"playerIndex"`
		Dropped                bool    `json:"dropped"`
		Kills                  int     `json:"kills"`
		Ping                   int     `json:"ping"`
		QueenKills             int     `json:"queenKills"`
		Deaths                 int     `json:"deaths"`
		Berries                int     `json:"berries"`
		Glances                int     `json:"glances"`
		Snail                  float64 `json:"snail"`
		SnailDeaths            int     `json:"snailDeaths"`
		BerryThrowIns          int     `json:"berryThrowIns"`
		MostQueenKillsInAMatch int     `json:"mostQueenKillsInAMatch"`
		MostKillsPerLife       int     `json:"mostKillsPerLife"`
		AllBerriesInSingleGame bool    `json:"allBerriesInSingleGame"`
		Team                   int     `json:"team"`
		EntityType             int     `json:"entityType"`
		EntitySkin             int     `json:"entitySkin"`
	} `json:"playerMatchStats"`
	GameWinners   []int `json:"gameWinners"`
	WinConditions []int `json:"winConditions"`
	MapPool       []int `json:"mapPool"`
	MatchType     int   `json:"matchType"`
	Profiles      []struct {
		LiquidID    string `json:"liquidId"`
		PlayerID    string `json:"playerId"`
		DisplayName string `json:"displayName"`
		ExternalIds struct {
			Discord  interface{} `json:"discord"`
			Nintendo interface{} `json:"nintendo"`
			Test     string      `json:"test"`
			Steam    string      `json:"steam"`
			Xboxone  interface{} `json:"xboxone"`
			Ps4      interface{} `json:"ps4"`
			Stadia   interface{} `json:"stadia"`
		} `json:"externalIds"`
		HighestRankingRating float64 `json:"highestRankingRating"`
		RankedRankingData    struct {
			Rating        float64 `json:"rating"`
			Deviation     float64 `json:"deviation"`
			Volatility    float64 `json:"volatility"`
			RoundedRating int     `json:"roundedRating"`
		} `json:"rankedRankingData"`
		UnrankedRankingData struct {
			Rating        float64 `json:"rating"`
			Deviation     float64 `json:"deviation"`
			Volatility    float64 `json:"volatility"`
			RoundedRating int     `json:"roundedRating"`
		} `json:"unrankedRankingData"`
		RankedRecord struct {
			Win  int `json:"win"`
			Loss int `json:"loss"`
		} `json:"rankedRecord"`
		UnrankedRecord struct {
			Win  int `json:"win"`
			Loss int `json:"loss"`
		} `json:"unrankedRecord"`
		PlacementRecord struct {
			Win  int `json:"win"`
			Loss int `json:"loss"`
		} `json:"placementRecord"`
		Stats struct {
			Kills          int     `json:"kills"`
			KillsAsWarrior int     `json:"killsAsWarrior"`
			KillsAsQueen   int     `json:"killsAsQueen"`
			Deaths         int     `json:"deaths"`
			Berries        int     `json:"berries"`
			SnailDistance  float64 `json:"snailDistance"`
			RunnerMinutes  float64 `json:"runnerMinutes"`
			WarriorMinutes float64 `json:"warriorMinutes"`
			QueenMinutes   float64 `json:"queenMinutes"`
		} `json:"stats"`
		RankingAdjustment                int           `json:"rankingAdjustment"`
		PreviousLeague                   int           `json:"previousLeague"`
		CurrentLeague                    int           `json:"currentLeague"`
		TimeoutExpiration                string        `json:"timeoutExpiration"`
		Status                           int           `json:"status"`
		AllowFriendsToJoinParty          bool          `json:"allowFriendsToJoinParty"`
		AllowFriendsOfFriendsToJoinParty bool          `json:"allowFriendsOfFriendsToJoinParty"`
		AllowFriendsToJoinCustomMatch    bool          `json:"allowFriendsToJoinCustomMatch"`
		AllowSpectateCustomMatch         bool          `json:"allowSpectateCustomMatch"`
		CurrentNetworkingPreferences     int           `json:"currentNetworkingPreferences"`
		CrossPlayEnabled                 bool          `json:"crossPlayEnabled"`
		PartyCount                       int           `json:"partyCount"`
		LocalPlayerCount                 int           `json:"localPlayerCount"`
		InParty                          bool          `json:"inParty"`
		PartyLeader                      bool          `json:"partyLeader"`
		AvatarURL                        string        `json:"avatarUrl"`
		RankInLeague                     int           `json:"rankInLeague"`
		DefaultAvatar                    int           `json:"defaultAvatar"`
		PromotionPercentComplete         float64       `json:"promotionPercentComplete"`
		Version                          int           `json:"version"`
		RankedV2Data                     []interface{} `json:"rankedV2Data"`
		RequiredPlacementMatches         float64       `json:"requiredPlacementMatches"`
		Type                             int           `json:"type"`
		LiquidFriend                     bool          `json:"liquidFriend"`
		RequiredPlacementMatchesString   string        `json:"requiredPlacementMatchesString"`
		Online                           bool          `json:"online"`
		InGame                           bool          `json:"inGame"`
		CombinedID                       string        `json:"combinedId"`
	} `json:"profiles"`
	Games []struct {
		PlayerStats []struct {
			Nickname                string  `json:"nickname"`
			ActorNr                 int     `json:"actorNr"`
			InputID                 int     `json:"inputID"`
			PlayerID                string  `json:"playerId"`
			ExternalPlayerID        string  `json:"externalPlayerId"`
			PlayerIndex             int     `json:"playerIndex"`
			EntityType              int     `json:"entityType"`
			EntitySkin              int     `json:"entitySkin"`
			Ping                    int     `json:"ping"`
			Team                    int     `json:"team"`
			Dropped                 bool    `json:"dropped"`
			CurrentWarriorKillCount int     `json:"currentWarriorKillCount"`
			CurrentWorkerKillCount  int     `json:"currentWorkerKillCount"`
			CurrentQueenKillCount   int     `json:"currentQueenKillCount"`
			CurrentBerryDeposits    int     `json:"currentBerryDeposits"`
			CurrentBerryThrowIns    int     `json:"currentBerryThrowIns"`
			CurrentGlanceCount      int     `json:"currentGlanceCount"`
			CurrentSnailDistance    float64 `json:"currentSnailDistance"`
			StartingSnailPos        struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"startingSnailPos"`
			TotalKillCount            int     `json:"totalKillCount"`
			TotalQueenKillCount       int     `json:"totalQueenKillCount"`
			TotalWorkerKillCount      int     `json:"totalWorkerKillCount"`
			TotalWarriorKillCount     int     `json:"totalWarriorKillCount"`
			TotalBerryDeposits        int     `json:"totalBerryDeposits"`
			TotalBerryThrowIns        int     `json:"totalBerryThrowIns"`
			TotalGlanceCount          int     `json:"totalGlanceCount"`
			TotalSnailDistance        float64 `json:"totalSnailDistance"`
			TotalSnailDeaths          int     `json:"totalSnailDeaths"`
			TotalDeathCount           int     `json:"totalDeathCount"`
			WarriorAndQueenDeathCount int     `json:"warriorAndQueenDeathCount"`
			WorkerDeathCount          int     `json:"workerDeathCount"`
			MostKillsPerLife          int     `json:"mostKillsPerLife"`
			TimeSpentAsWarriorSeconds float64 `json:"timeSpentAsWarriorSeconds"`
			CurrentKillCount          int     `json:"currentKillCount"`
		} `json:"playerStats"`
		BlueQueenKillTimes []float64 `json:"blueQueenKillTimes"`
		GoldQueenKillTimes []float64 `json:"goldQueenKillTimes"`
		BlueWarriorsUp     []struct {
			Time  float64 `json:"time"`
			Count int     `json:"count"`
		} `json:"blueWarriorsUp"`
		GoldWarriorsUp []struct {
			Time  float64 `json:"time"`
			Count int     `json:"count"`
		} `json:"goldWarriorsUp"`
		SnailPosition []struct {
			Time  float64 `json:"time"`
			Count int     `json:"count"`
		} `json:"snailPosition"`
		BlueBerryCount []struct {
			Time  float64 `json:"time"`
			Count int     `json:"count"`
		} `json:"blueBerryCount"`
		GoldBerryCount []interface{} `json:"goldBerryCount"`
		GateControls   []struct {
			ID         int     `json:"id"`
			TimeAsBlue float64 `json:"timeAsBlue"`
			TimeAsRed  float64 `json:"timeAsRed"`
		} `json:"gateControls"`
		BerriesNeeded    int     `json:"berriesNeeded"`
		TotalGates       int     `json:"totalGates"`
		BlueSnailGatePos float64 `json:"blueSnailGatePos"`
		GoldSnailGatePos float64 `json:"goldSnailGatePos"`
		SnailDisengaged  bool    `json:"snailDisengaged"`
		HeadShotVictory  bool    `json:"headShotVictory"`
		StartTime        float64 `json:"startTime"`
		EndTime          float64 `json:"endTime"`
		WinCondition     int     `json:"winCondition"`
		Duration         float64 `json:"duration"`
	} `json:"games"`
}

func (statJson *StatsJSON) Players() []string {
	names := make([]string, len(statJson.PlayerMatchStats))
	for idx, profile := range statJson.PlayerMatchStats {
		names[idx] = profile.Nickname
	}
	return names
}

func (statJson *StatsJSON) MapsWon() map[string]int {
	blueMaps := 0
	goldMaps := 0
	for _, val := range statJson.GameWinners {
		if val == 1 {
			goldMaps++
		} else if val == 2 {
			blueMaps++
		}
	}
	output := map[string]int{
		"blue": blueMaps,
		"gold": goldMaps,
	}
	return output
}

func (statJson *StatsJSON) AdvancedStats() []map[string]map[string]int {
	goldStats := make(map[string]map[string]int)
	blueStats := make(map[string]map[string]int)

	for _, player := range statJson.PlayerMatchStats {
		playerMap := make(map[string]int)
		if player.Team == 1 {
			goldStats[player.Nickname] = playerMap
		} else {
			blueStats[player.Nickname] = playerMap
		}
	}
	for _, game := range statJson.Games {
		for _, player := range game.PlayerStats {
			if player.Team == 1 {
				goldStats[player.Nickname]["QueenKills"] += player.TotalQueenKillCount
				goldStats[player.Nickname]["WarriorKills"] += player.TotalWarriorKillCount
				goldStats[player.Nickname]["WorkerKills"] += player.TotalWorkerKillCount
				goldStats[player.Nickname]["WarriorDeaths"] += player.WarriorAndQueenDeathCount
				goldStats[player.Nickname]["WorkerDeaths"] += player.WorkerDeathCount
				goldStats[player.Nickname]["Team"] += player.Team
				goldStats[player.Nickname]["WarriorUptime"] += int(player.TimeSpentAsWarriorSeconds)
				goldStats[player.Nickname]["BerryDunks"] += (player.TotalBerryDeposits - player.TotalBerryThrowIns)
				goldStats[player.Nickname]["BerryThrows"] += player.TotalBerryThrowIns
				goldStats[player.Nickname]["EntityType"] += player.EntityType
				goldStats[player.Nickname]["Snail"] += int(player.TotalSnailDistance)
			} else {
				blueStats[player.Nickname]["QueenKills"] += player.TotalQueenKillCount
				blueStats[player.Nickname]["WarriorKills"] += player.TotalWarriorKillCount
				blueStats[player.Nickname]["WorkerKills"] += player.TotalWorkerKillCount
				blueStats[player.Nickname]["WarriorDeaths"] += player.WarriorAndQueenDeathCount
				blueStats[player.Nickname]["WorkerDeaths"] += player.WorkerDeathCount
				blueStats[player.Nickname]["Team"] += player.Team
				blueStats[player.Nickname]["WarriorUptime"] += int(player.TimeSpentAsWarriorSeconds)
				blueStats[player.Nickname]["BerryDunks"] += (player.TotalBerryDeposits - player.TotalBerryThrowIns)
				blueStats[player.Nickname]["BerryThrows"] += player.TotalBerryThrowIns
				blueStats[player.Nickname]["EntityType"] += player.EntityType
				blueStats[player.Nickname]["Snail"] += int(player.TotalSnailDistance)
			}
		}
	}
	output := make([]map[string]map[string]int, 2)
	output[0] = goldStats
	output[1] = blueStats
	return output
}
