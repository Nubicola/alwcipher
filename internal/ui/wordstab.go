package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/calculator"
	"github.com/Nubicola/alwcipher/pkg/corpus"
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

type calculators struct {
	base, first, last calculator.ALWCalculator
}

// requires the input and the output as this is where the connection happens between the two!
func makeInputArea(i *boundInput, o *boundOutput, c *corpus.Corpus, calcs *calculators) fyne.CanvasObject {
	ui := new(userInputUI)
	// closures are great!
	calcfunc := func() {
		s, _ := i.userInput.Get()
		// calculate the value of that thing
		var val = c.Calculate(s)

		o.baseBoundValue.Set(calcs.base.StringValue(s))
		o.firstBoundValue.Set(calcs.first.Calculate(strings.Split(s, " ")))
		o.lastBoundValue.Set(calcs.last.Calculate(strings.Split(s, " ")))

		// build the string to show in the output area. It's either the corpus after adding this string,
		// or this string appended to the corpus' string (thus not affecting the corpus itself)
		var ss = make([]string, 1)
		if fyne.CurrentApp().Preferences().Bool("WriteToCorpus") {
			c.Add(s)
			ss = c.Get(val)
		} else {
			ss = append(c.Get(val), s)
		}
		o.outputFieldBoundValue.Set(strings.Join(ss, "\n"))
	}

	ui.calculateButton = widget.NewButton("Calculate", calcfunc)
	ui.userInput = makeNewWordKeyedEntry()
	ui.userInput.Bind(i.userInput)
	ui.userInput.OnTypedKey = calcfunc
	box := container.NewGridWithColumns(2, ui.userInput, ui.calculateButton)
	return box
}

func MakeWordsTabUI(c *corpus.Corpus) fyne.CanvasObject {

	calcs := new(calculators)
	calcs.base = new(calculator.EQBaseCalculator)
	calcs.first = new(calculator.EQFirstCalculator)
	calcs.last = new(calculator.EQLastCalculator)

	bo := makeNewBoundOutputs()
	ou := makeNewOutputUI(bo)

	bi := makeNewBoundInput()
	ia := makeInputArea(bi, bo, c, calcs)

	// write to corpus?
	writeToCorpusCheckbox := widget.NewCheck("Save to corpus?", func(value bool) {
		fyne.CurrentApp().Preferences().SetBool("WriteToCorpus", value)
	})
	writeToCorpusCheckbox.Checked = fyne.CurrentApp().Preferences().BoolWithFallback("WriteToCorpus", true)
	//writeToCorpusCheckbox.Disable() // not ready yet

	box := container.NewVBox(ia, writeToCorpusCheckbox, ou)

	/*var jbutton = widget.NewButton("Marshal", func() {
		var b bytes.Buffer
		err := c.Save(&b)
		if err != nil {
			bo.outputFieldBoundValue.Set(string(err.Error()))
		}
		bo.outputFieldBoundValue.Set(b.String())
	})*/

	return box
}
