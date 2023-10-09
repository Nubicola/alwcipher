package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	naeq_go "github.com/Nubicola/NAEQ_GO"
)

func getALWValuesFromCorpus(alw *naeq_go.NAEQ_Processor, s string) string {
	return "hello"
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ALW Visualized")
	myWindow.Resize(fyne.Size{
		Width:  460,
		Height: 500,
	})

	//var alwProcessor *naeq_go.NAEQ_Processor = naeq_go.MakeNewNQP("alw")
	var eqbc = new(naeq_go.EQBaseCalculator)
	var eqfc = new(naeq_go.EQFirstCalculator)
	var eqlc = new(naeq_go.EQLastCalculator)

	var alwProcessor = naeq_go.MakeNewNQP("alw")

	// set up the output values area
	// the BoundValues can be updated which will auto-update the widgets
	alwBoundValue := binding.NewInt()
	firstBoundValue := binding.NewInt()
	lastBoundValue := binding.NewInt()

	outputBoundValue := binding.NewString()
	userInputBoundValue := binding.NewString()

	// right side of the border layout. This is pretty much the key action:
	// user presses "calculate" and all the magic happens
	calculateButton := widget.NewButton("Calculate", func() {
		s, _ := userInputBoundValue.Get()

		var ns = []string{}
		ns = append(ns, s)
		alwBoundValue.Set(eqbc.Calculate(ns))
		// first and last need the string to be split by whitespace into a list
		var words = strings.Fields(s)
		firstBoundValue.Set(eqfc.Calculate(words))
		lastBoundValue.Set(eqlc.Calculate(words))

		// here we would load all matching items into the output field
		outputBoundValue.Set(getALWValuesFromCorpus(alwProcessor, s))
	})

	alwLabel := widget.NewLabel("ALW")
	firstLabel := widget.NewLabel("First")
	lastLabel := widget.NewLabel("Last")

	alwBoundValueStr := binding.IntToString(alwBoundValue)
	firstBoundValueStr := binding.IntToString(firstBoundValue)
	lastBoundValueStr := binding.IntToString(lastBoundValue)
	alwLabelVal := widget.NewLabelWithData(alwBoundValueStr)
	firstLabelVal := widget.NewLabelWithData(firstBoundValueStr)
	lastLabelVal := widget.NewLabelWithData(lastBoundValueStr)

	outputValueGrid := container.New(layout.NewFormLayout(), alwLabel, alwLabelVal, firstLabel, firstLabelVal, lastLabel, lastLabelVal)
	rightBox := container.NewVBox(calculateButton, outputValueGrid)

	// menu bar on top
	myWindow.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Sync...", func() {
				//
			}))))
	/*	toolbar := widget.NewToolbar(
			widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
			widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
			widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(theme.HelpIcon(), func() {
			}),
		)
	*/
	// container.NewBorder(top, bottom, left, right, center)

	// the middle; a vertical box
	userInput := widget.NewEntryWithData(userInputBoundValue)

	// write to corpus?
	writeToCorpusCheckbox := widget.NewCheck("Save to corpus?", func(value bool) {
		// nothing
	})
	writeToCorpusCheckbox.Checked = false
	writeToCorpusCheckbox.Disable() // not ready yet

	// the output area

	outputField := widget.NewEntryWithData(outputBoundValue)
	//outputField := widget.NewLabelWithData(outputBoundValue)
	outputField.MultiLine = true
	outputField.Wrapping = fyne.TextWrapBreak
	// fix the color by customizing theme: https://developer.fyne.io/extend/custom-theme
	//outputField.Disable()

	middleBox := container.NewVBox(userInput, writeToCorpusCheckbox, outputField)

	content := container.NewBorder(nil, nil, nil, rightBox, middleBox)
	myWindow.SetContent(content)

	// just for now...
	fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
	myWindow.ShowAndRun()
}
