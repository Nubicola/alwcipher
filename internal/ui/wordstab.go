package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// changes to the input value will update all these
type boundOutput struct {
	// these are the integer values that receive updates
	baseBoundValue, firstBoundValue, lastBoundValue binding.Int
	// these are string representations bound to those values
	baseBoundValueStr, firstBoundValueStr, lastBoundValueStr binding.String
	// bound to the string output area
	outputFieldBoundValue binding.String
}

func makeNewBoundOutputs() *boundOutput {
	bov := new(boundOutput)
	bov.baseBoundValue = binding.NewInt()
	bov.baseBoundValueStr = binding.IntToString(bov.baseBoundValue)
	bov.firstBoundValue = binding.NewInt()
	bov.firstBoundValueStr = binding.IntToString(bov.firstBoundValue)
	bov.lastBoundValue = binding.NewInt()
	bov.lastBoundValueStr = binding.IntToString(bov.lastBoundValue)
	bov.outputFieldBoundValue = binding.NewString()

	return bov
}

type outputUI struct {
	baseLabel, firstLabel, lastLabel          *widget.Label
	baseLabelVal, firstLabelVal, lastLabelVal *widget.Label
	outputField                               *widget.Entry
}

func makeNewOutputUI(b *boundOutput) fyne.CanvasObject {
	oui := new(outputUI)
	oui.baseLabel = widget.NewLabel("ALW")
	oui.baseLabelVal = widget.NewLabelWithData(b.baseBoundValueStr)
	oui.firstLabel = widget.NewLabel("First")
	oui.firstLabelVal = widget.NewLabelWithData(b.firstBoundValueStr)
	oui.lastLabel = widget.NewLabel("Last")
	oui.lastLabelVal = widget.NewLabelWithData(b.lastBoundValueStr)

	oui.outputField = widget.NewEntryWithData(b.outputFieldBoundValue)
	oui.outputField.MultiLine = true
	oui.outputField.Wrapping = fyne.TextWrapBreak
	oui.outputField.Disable()
	outputValueGrid := container.New(layout.NewFormLayout(), oui.baseLabel, oui.baseLabelVal, oui.firstLabel, oui.firstLabelVal, oui.lastLabel, oui.lastLabelVal)
	box := container.NewBorder(nil, nil, nil, outputValueGrid, oui.outputField)
	//box := container.NewGridWithColumns(2, oui.outputField, outputValueGrid)
	return box
}

type boundInput struct {
	userInput binding.String
}

func makeNewBoundInput() *boundInput {
	i := new(boundInput)
	i.userInput = binding.NewString()
	return i
}

type wordsKeyedEntry struct {
	widget.Entry
	OnTypedKey func()
}

func makeNewWordKeyedEntry() *wordsKeyedEntry {
	entry := &wordsKeyedEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *wordsKeyedEntry) TypedKey(key *fyne.KeyEvent) {
	e.Entry.TypedKey(key)
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		e.OnTypedKey()
	}
}

type userInputUI struct {
	calculateButton *widget.Button
	userInput       *wordsKeyedEntry
}

// requires the input and the output as this is where the connection happens between the two!
func makeInputArea(i *boundInput, o *boundOutput, fc FyneCorpus) fyne.CanvasObject {
	ui := new(userInputUI)
	// closures are great!
	calcfunc := func() {
		s, err := i.userInput.Get()
		if err != nil || len(s) < 1 {
			return
		}

		cleanString, _ := fc.CleanString(s)

		ss := fc.Add(cleanString)

		// calculate the value of that thing
		//val, verr := c.Calculate(s)
		bval, berr := fc.base.StringValue(cleanString)
		fval, ferr := fc.first.Calculate(strings.Split(cleanString, " "))
		lval, lerr := fc.last.Calculate(strings.Split(cleanString, " "))
		if berr != nil || ferr != nil || lerr != nil {
			o.outputFieldBoundValue.Set("invalid characters; use english words without numerals only please")
		} else {
			o.baseBoundValue.Set(bval)
			o.firstBoundValue.Set(fval)
			o.lastBoundValue.Set(lval)
			o.outputFieldBoundValue.Set(strings.Join(ss, "\n"))
		}
	}

	ui.calculateButton = widget.NewButton("Calculate", calcfunc)
	ui.userInput = makeNewWordKeyedEntry()
	ui.userInput.Bind(i.userInput)
	ui.userInput.OnTypedKey = calcfunc
	box := container.NewGridWithColumns(2, ui.userInput, ui.calculateButton)
	return box
}

func MakeWordsTabUI(fc *FyneCorpus) fyne.CanvasObject {

	bo := makeNewBoundOutputs()
	ou := makeNewOutputUI(bo)

	bi := makeNewBoundInput()
	ia := makeInputArea(bi, bo, *fc)

	// write to corpus?
	writeToCorpusCheckbox := widget.NewCheck("Save to corpus?", func(value bool) {
		fyne.CurrentApp().Preferences().SetBool("WriteToCorpus", value)
	})
	writeToCorpusCheckbox.Checked = fyne.CurrentApp().Preferences().BoolWithFallback("WriteToCorpus", true)

	box := container.NewVBox(ia, writeToCorpusCheckbox)
	border := container.NewBorder(box, nil, nil, nil, ou)
	return border
}
