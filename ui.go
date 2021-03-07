package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func ShowAdvancedStats(app *fyne.App, data stats.StatsJSON) {
	a := *app
	w := a.NewWindow("Advanced Stats")

	advStatsPlot := stats.PlotStats(data)
	objStatsPlot := stats.PlotObjectiveStats(data)
	advStatsCanvas := canvas.NewImageFromImage(advStatsPlot)
	advStatsCanvas.SetMinSize(fyne.NewSize(1280, 720))
	objStatsCanvas := canvas.NewImageFromImage(objStatsPlot)
	objStatsCanvas.SetMinSize(fyne.NewSize(1280, 720))

	advStatsCanvas.FillMode = canvas.ImageFillContain
	objStatsCanvas.FillMode = canvas.ImageFillContain
	cont := container.NewVBox()
	nextButton := widget.NewButtonWithIcon("Military", theme.MediaSkipNextIcon(), func() {
		cont.Objects[0] = advStatsCanvas
		cont.Refresh()
	})
	prevButton := widget.NewButtonWithIcon("Objective", theme.MediaSkipPreviousIcon(), func() {
		cont.Objects[0] = objStatsCanvas
		cont.Refresh()
	})
	controls := container.NewHBox(layout.NewSpacer(), prevButton, nextButton, layout.NewSpacer())
	cont.Add(advStatsCanvas)
	cont.Add(controls)
	w.CenterOnScreen()
	w.SetContent(cont)
	w.Show()
}
