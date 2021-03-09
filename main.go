package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func main() {
	names := stats.ListStatFiles()
	data := stats.ReadJson(names[0])

	a := app.New()
	w := a.NewWindow("KQB JSON Viewer")

	kqbApp := KQBApp{
		files:        names,
		selectedData: data,
		a:            a,
		w:            w,
	}

	kqbApp.ShowMainWindow()

}
