package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	ui "github.com/Nubicola/alwcipher/internal/ui"
	calculator "github.com/Nubicola/alwcipher/pkg/calculator"
	corpus "github.com/Nubicola/alwcipher/pkg/corpus"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ALW Corpus Visualized")
	myWindow.Resize(fyne.Size{
		Width:  460,
		Height: 500,
	})

	var eqbc = new(calculator.EQBaseCalculator)
	/*	var eqfc = new(calculator.EQFirstCalculator)
		var eqlc = new(calculator.EQLastCalculator)*/

	var alwbasecorpus = corpus.NewCorpus(eqbc)

	myWindow.SetContent(ui.MakeUI(alwbasecorpus))

	// just for now...
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	myWindow.ShowAndRun()
}
