package ui

import (
	"bufio"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
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

func MakeCorpusTabUI(fc *FyneCorpus, w fyne.Window) fyne.CanvasObject {

	bo := makeNewBoundCorpusOutputs()
	ou := makeNewCorpusOutputUI(bo)

	radio := widget.NewRadioGroup([]string{"Words", "Lines" /*, "Sentences", "Phrases"*/}, func(value string) {
		bo.importedMethodStr.Set(value)
	})
	radio.SetSelected("Words")

	importFile := func(reader fyne.URIReadCloser) (int, error) {
		scanner := bufio.NewScanner(reader)
		mode, er := bo.importedMethodStr.Get()
		i := 0
		if er == nil {
			switch m := mode; m {
			case "Words":
				scanner.Split(bufio.ScanWords)
			case "Lines":
				scanner.Split(bufio.ScanLines)
			}
			i := 0
			for scanner.Scan() {
				_ = fc.Add(scanner.Text())
				i += 1
			}
		}
		return i, er
	}

	fileImportButton := widget.NewButton("Import File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				i, er := importFile(reader)
				if er == nil {
					bo.importedBoundValue.Set(i)
				}
			}
		}, w)
	})

	folderImportButton := widget.NewButton("Import Folder", func() {
		dialog.ShowFolderOpen(func(listable fyne.ListableURI, err error) {
			if err == nil && listable != nil {
				l, _ := listable.List()
				i := 0
				for _, f := range l {
					s, eee := storage.Reader(f)
					if eee == nil {
						j, er := importFile(s)
						if er == nil {
							i += j
						}
					}
				}
				bo.importedBoundValue.Set(i)
			}
		}, w)
	})

	box := container.NewVBox(radio, fileImportButton, folderImportButton, ou)
	return box
}
