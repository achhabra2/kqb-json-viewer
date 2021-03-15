package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func main() {
	names := stats.ListStatFiles()
	data := stats.ReadJson(names[0])

	a := app.NewWithID("com.kqb-json-viewer.app")

	appTheme := myTheme{}
	a.Settings().SetTheme(&appTheme)

	w := a.NewWindow("KQB JSON Viewer")

	kqbApp := KQBApp{
		files:        names,
		selectedData: data,
		a:            a,
		w:            w,
	}

	kqbApp.ShowMainWindow()

}
