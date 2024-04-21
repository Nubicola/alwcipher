package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type boundNumInput struct {
	userInput    binding.String
	userNumInput binding.Int
}

func makeNewBoundNumInput() *boundNumInput {
	i := new(boundNumInput)
	i.userInput = binding.NewString()
	i.userNumInput = binding.StringToInt(i.userInput)
	return i
}

// duplicate type from WordsTab, could be moved to common...later
type numKeyedEntry struct {
	widget.Entry
	OnTypedKey func()
}

// https://docs.fyne.io/extend/numerical-entry move to this to ensure only numbers are entered
func makeNewNumKeyedEntry() *numKeyedEntry {
	entry := &numKeyedEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numKeyedEntry) TypedKey(key *fyne.KeyEvent) {
	e.Entry.TypedKey(key)
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		e.OnTypedKey()
	}
}

type userNumInputUI struct {
	calculateButton *widget.Button
	userInput       *numKeyedEntry
}

func makeNumInputArea(i *boundNumInput, o *boundNumOutput, fc FyneCorpus) fyne.CanvasObject {
	ui := new(userNumInputUI)
	calc := func() {
		//var st, _ = i.userInput.Get()
		var v, err1 = i.userNumInput.Get()
		if err1 == nil {
			if v > 0 {
				o.matchingText.Set(strings.Join(fc.c.Get(v), "\n"))
			}
		}
	}

	ui.calculateButton = widget.NewButton("Show matches", calc)
	ui.userInput = makeNewNumKeyedEntry()
	ui.userInput.Bind(i.userInput)
	ui.userInput.OnTypedKey = calc
	box := container.NewGridWithColumns(2, ui.userInput, ui.calculateButton)
	return box
}

type boundNumOutput struct {
	matchingText binding.String
}

func makeNewBoundNumOutput() *boundNumOutput {
	o := new(boundNumOutput)
	o.matchingText = binding.NewString()
	return o
}

func makeNumOutputArea(bo *boundNumOutput, _ FyneCorpus) fyne.CanvasObject {
	w := widget.NewEntryWithData(bo.matchingText)
	w.MultiLine = true
	w.Wrapping = fyne.TextWrapBreak
	w.Disable()
	return container.NewBorder(nil, nil, nil, nil, w)
	//return r
}

func MakeNumbersTabUI(fc *FyneCorpus) fyne.CanvasObject {
	bo := makeNewBoundNumOutput()
	oa := makeNumOutputArea(bo, *fc)
	bi := makeNewBoundNumInput()
	ia := makeNumInputArea(bi, bo, *fc)
	return container.NewBorder(ia, nil, nil, nil, oa)
}
