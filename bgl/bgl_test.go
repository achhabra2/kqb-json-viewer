package bgl

import "testing"

func TestGetMatches(t *testing.T) {
	bglData := BGLData{}
	err := bglData.LoadCurrentMatches()
	if err != nil {
		t.Log("Could not load matches from API. ", err)
		t.Fail()
	}
}

func TestGetMatchesLocal(t *testing.T) {
	bglData := BGLData{}
	err := bglData.LoadCurrentMatchesLocal()
	if err != nil {
		t.Log("Could not load matches from local json test file. ", err)
		t.Fail()
	}
}
