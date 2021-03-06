package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func ShowAdvancedStats(app *fyne.App, data stats.StatsJSON) {
	a := *app
	w := a.NewWindow("Advanced Stats")

	plot := stats.PlotStats(data)

	image := canvas.NewImageFromImage(plot)

	image.FillMode = canvas.ImageFillOriginal

	w.SetContent(image)
	w.Show()
}
