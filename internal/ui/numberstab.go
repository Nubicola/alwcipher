package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/corpus"
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

func makeNumInputArea(i *boundNumInput, o *boundNumOutput, c *corpus.Corpus) fyne.CanvasObject {
	ui := new(userNumInputUI)
	calc := func() {
		var v, err = i.userNumInput.Get()
		if err != nil {
			o.textOutput.Set("error")
		} else {
			o.textOutput.Set(strings.Join(c.Get(v), "\n"))
		}
	}

	ui.calculateButton = widget.NewButton("Show me", calc)
	ui.userInput = makeNewNumKeyedEntry()
	ui.userInput.Bind(i.userInput)
	ui.userInput.OnTypedKey = calc
	box := container.NewGridWithColumns(2, ui.userInput, ui.calculateButton)
	return box
}

type boundNumOutput struct {
	textOutput binding.String
}

func makeNewBoundNumOutput() *boundNumOutput {
	o := new(boundNumOutput)
	o.textOutput = binding.NewString()
	return o
}

func makeNumOutputArea(bo *boundNumOutput, c *corpus.Corpus) fyne.CanvasObject {
	w := widget.NewEntryWithData(bo.textOutput)
	w.MultiLine = true
	w.Wrapping = fyne.TextWrapBreak
	return w
}

func MakeNumbersTabUI(c *corpus.Corpus) fyne.CanvasObject {

	/*calcs := new(calculators)
	calcs.base = new(calculator.EQBaseCalculator)
	calcs.first = new(calculator.EQFirstCalculator)
	calcs.last = new(calculator.EQLastCalculator)*/

	/*	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			switch id {
			case "":
				return []widget.TreeNodeID{"a", "b", "c"}
			case "a":
				return []widget.TreeNodeID{"a1", "a2"}
			}
			return []string{}
		},
		func(id widget.TreeNodeID) bool {
			return id == "" || id == "a"
		},
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("Branch template")
			}
			return widget.NewLabel("Leaf template")
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			text := id
			if branch {
				text += " (branch)"
			}
			o.(*widget.Label).SetText(text)
		})*/

	bo := makeNewBoundNumOutput()
	oa := makeNumOutputArea(bo, c)
	bi := makeNewBoundNumInput()
	ia := makeNumInputArea(bi, bo, c)
	return container.NewBorder(ia, nil, nil, nil, oa)
	//return container.NewVBox(ia, oa)
	//box := container.NewGridWithRows(4, layout.NewSpacer(), ia, layout.NewSpacer(), oa)
}
