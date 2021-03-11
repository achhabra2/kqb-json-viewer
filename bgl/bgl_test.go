package bgl

import (
	"testing"
)

func TestGetMe(t *testing.T) {
	want := "prosive"
	bglData := BGLData{
		Token: "e19f5efd2fb9f6abd3d4522d18d08b2db7754c24",
	}
	bglData.GetMe()
	have := bglData.User.FirstName

	if have != want {
		t.Errorf("Error with user retreival, got: %v, have: %v\n", have, want)
	}
}

func TestGetMatches(t *testing.T) {
	bglData := BGLData{
		Token: "e19f5efd2fb9f6abd3d4522d18d08b2db7754c24",
	}
	bglData.GetMe()
	err := bglData.LoadCurrentMatches()
	if err != nil {
		t.Error("Could not load matches from API. ", err)
	}
}

func TestGetMatchesLocal(t *testing.T) {
	bglData := BGLData{}
	err := bglData.LoadCurrentMatchesLocal()
	if err != nil {
		t.Error("Could not load matches from local json test file. ", err)
	}
}
