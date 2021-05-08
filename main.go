package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/achhabra2/kqb-json-viewer/icons"
	"github.com/achhabra2/kqb-json-viewer/stats"
)

func main() {
	setupLogs()
	os.Setenv("FYNE_SCALE", "0.9")
	mainIcon := fyne.NewStaticResource("logo.png", icons.Logo)

	a := app.NewWithID("com.kqb-json-viewer.app")
	a.SetIcon(mainIcon)
	appTheme := myTheme{}
	a.Settings().SetTheme(&appTheme)
	w := a.NewWindow("KQB JSON Viewer")

	names := stats.ListStatFiles()
	if len(names) == 0 {
		err := errors.New("No Stat Files Available to Read")
		log.Println("No Stat Files Available to Read")
		dialog := dialog.NewError(err, w)
		dialog.SetOnClosed(func() { os.Exit(1) })
		dialog.Show()
		w.Resize(fyne.NewSize(500, 850))
		w.CenterOnScreen()
		w.ShowAndRun()
	}
	data, err := stats.ReadJson(names[0])
	if err != nil {
		dialog := dialog.NewError(err, w)
		dialog.SetOnClosed(func() { os.Exit(1) })
		dialog.Show()
		w.Resize(fyne.NewSize(500, 850))
		w.CenterOnScreen()
		w.ShowAndRun()
	} else {
		basePath, _ := filepath.Split(names[0])
		kqbApp := KQBApp{
			files:        names,
			selectedData: data,
			a:            a,
			w:            w,
			basePath:     basePath,
		}

		kqbApp.ShowMainWindow()
	}

}

func setupLogs() {
	f, err := os.OpenFile("./kqb-json-viewer-output.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	log.SetOutput(f)
}
