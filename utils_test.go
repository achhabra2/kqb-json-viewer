package main

import "testing"

func TestTimeParsing(t *testing.T) {
	filename := "Custom-2021-05-06-20-22-26.json"
	matchTime, _ := FileNameToTime(filename)
	t.Log(matchTime)
	t.Log(matchTime.Format("3:04PM"))
}
