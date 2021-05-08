package main

import "testing"

func TestTimeParsing(t *testing.T) {
	filename := "Custom-2021-05-06-20-22-26.json"
	matchTime := FileNameToTime("", filename)
	t.Log(matchTime)
	t.Log(matchTime.Format("3:04PM"))

	fname2 := "sample.json"
	bpath2 := "/Users/aman/Documents/code/kqb-json-viewer/fixtures"
	matchTime2 := FileNameToTime(bpath2, fname2)
	t.Log(matchTime2)
	t.Log(matchTime2.Format("3:04PM"))
}
