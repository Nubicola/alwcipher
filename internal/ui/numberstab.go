package ui

import (
	"log"
	"strconv"

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

func makeNumInputArea(i *boundNumInput, o *boundNumOutput, w fyne.CanvasObject, c *corpus.Corpus) fyne.CanvasObject {
	ui := new(userNumInputUI)
	calc := func() {
		var v, err = i.userNumInput.Get()
		if err == nil && v > 0 {
			// this is kludgy
			a := make(map[string][]string)
			b := make(map[string]string)
			o.tree.Set(a, b)
			o.tree.Append(binding.DataTreeRootID, strconv.Itoa(v), strconv.Itoa(v))
			for i, s := range c.Get(v) {
				log.Println(i, s)
				o.tree.Append(strconv.Itoa(v), s, s)
			}
			w.Refresh()
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
	//textOutput binding.String
	tree binding.StringTree
}

func makeNewBoundNumOutput() *boundNumOutput {
	o := new(boundNumOutput)
	//o.textOutput = binding.NewString()
	o.tree = binding.NewStringTree()
	//o.tree = binding.NewIntTree()
	return o
}

func makeNumOutputArea(bo *boundNumOutput, c *corpus.Corpus) fyne.CanvasObject {
	/*	w := widget.NewEntryWithData(bo.textOutput)
		w.MultiLine = true
		w.Wrapping = fyne.TextWrapBreak*/

	treewg := widget.NewTreeWithData(
		/* (data binding.DataTree, createItem func(bool) fyne.CanvasObject, updateItem func(binding.DataItem, bool, fyne.CanvasObject)) *widget.Tree*/
		bo.tree,
		func(_ bool) fyne.CanvasObject {
			return widget.NewLabel("hello")
		},
		func(data binding.DataItem, isParent bool, obj fyne.CanvasObject) {

			l := obj.(*widget.Label)
			l.Bind(data.(binding.String))
		},
	)

	return treewg
}

func MakeNumbersTabUI(c *corpus.Corpus) fyne.CanvasObject {

	//tree.Append()
	// iterate over the keys of the corpus without knowing how the corpus is implemented
	// GetKeys perhaps?
	// currently Corpus exports Eqs so just grab it directly :P. Not the best design.

	bo := makeNewBoundNumOutput()
	oa := makeNumOutputArea(bo, c)
	bi := makeNewBoundNumInput()
	ia := makeNumInputArea(bi, bo, oa, c)
	return container.NewBorder(ia, nil, nil, nil, oa)
	//return container.NewVBox(ia, oa)
	//box := container.NewGridWithRows(4, layout.NewSpacer(), ia, layout.NewSpacer(), oa)
}
