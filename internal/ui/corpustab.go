package ui

import (
	"bufio"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/calculator"
	"github.com/Nubicola/alwcipher/pkg/corpus"
)

// changes to the input value will update all these
type boundCorpusOutput struct {
	// label/str that shows data about what's been imported
	importedBoundValue    binding.Int
	importedBoundValueStr binding.String
	importedMethodStr     binding.String
}

func makeNewBoundCorpusOutputs() *boundCorpusOutput {
	bov := new(boundCorpusOutput)
	bov.importedBoundValue = binding.NewInt()
	bov.importedBoundValueStr = binding.IntToString(bov.importedBoundValue)
	bov.importedMethodStr = binding.NewString()
	return bov
}

type outputCorpusUI struct {
	importedValLabel, importedLabel, importedMethodLabel *widget.Label
}

func makeNewCorpusOutputUI(b *boundCorpusOutput) fyne.CanvasObject {
	oui := new(outputCorpusUI)
	oui.importedValLabel = widget.NewLabelWithData(b.importedBoundValueStr)
	oui.importedLabel = widget.NewLabel("Entries imported:")
	oui.importedMethodLabel = widget.NewLabelWithData(b.importedMethodStr)

	box := container.NewHBox(oui.importedLabel, oui.importedValLabel, oui.importedMethodLabel)
	return box
}

func MakeCorpusTabUI(c *corpus.Corpus, w fyne.Window) fyne.CanvasObject {

	calcs := new(calculators)
	calcs.base = new(calculator.EQBaseCalculator)
	calcs.first = new(calculator.EQFirstCalculator)
	calcs.last = new(calculator.EQLastCalculator)

	bo := makeNewBoundCorpusOutputs()
	ou := makeNewCorpusOutputUI(bo)

	radio := widget.NewRadioGroup([]string{"Words", "Lines" /*, "Sentences", "Phrases"*/}, func(value string) {
		bo.importedMethodStr.Set(value)
	})
	radio.SetSelected("Words")

	button := widget.NewButton("Import", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				scanner := bufio.NewScanner(reader)
				mode, er := bo.importedMethodStr.Get()
				if er == nil {
					switch m := mode; m {
					case "Words":
						scanner.Split(bufio.ScanWords)
					case "Lines":
						scanner.Split(bufio.ScanLines)
					}
					i := 0
					i, err = c.ReadFromScanner(scanner)
					bo.importedBoundValue.Set(i)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}, w)
	})

	box := container.NewVBox(radio, button, ou)

	return box
}
