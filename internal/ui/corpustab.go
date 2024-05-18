package ui

import (
	"bufio"
	"errors"

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
	importedBoundValue, corpusCountBV     binding.Int
	importedBoundValueStr, corpusCountBVS binding.String
	importedMethodStr, exportMethodStr    binding.String
}

func makeNewBoundCorpusOutputs() *boundCorpusOutput {
	bov := new(boundCorpusOutput)
	bov.importedBoundValue = binding.NewInt()
	bov.corpusCountBV = binding.NewInt()
	bov.importedBoundValueStr = binding.IntToString(bov.importedBoundValue)
	bov.corpusCountBVS = binding.IntToString(bov.corpusCountBV)
	bov.importedMethodStr = binding.NewString()
	bov.exportMethodStr = binding.NewString()
	return bov
}

/*func makeCorpusOutputArea() fyne.CanvasObject {
	w := widget.NewEntry()
	w.MultiLine = true
	w.Wrapping = fyne.TextWrapBreak
	w.Disable()
	return container.NewBorder(nil, nil, nil, nil, w)
	//return r
}*/

type outputCorpusUI struct {
	/*corpusCountLabel, corpusCountString,*/ importedValLabel, importedLabel, importedMethodLabel *widget.Label
}

func makeNewCorpusOutputUI(b *boundCorpusOutput) fyne.CanvasObject {
	oui := new(outputCorpusUI)
	oui.importedValLabel = widget.NewLabelWithData(b.importedBoundValueStr)
	oui.importedLabel = widget.NewLabel("Recent imports:")
	oui.importedMethodLabel = widget.NewLabelWithData(b.importedMethodStr)
	/*oui.corpusCountString = widget.NewLabel("Corpus size:")
	oui.corpusCountLabel = widget.NewLabelWithData(b.corpusCountBVS)*/

	ibox := container.NewHBox(oui.importedLabel, oui.importedValLabel, oui.importedMethodLabel)
	//cbox := container.NewHBox(oui.corpusCountString, oui.corpusCountLabel)
	return ibox //container.NewVBox(cbox, ibox)
}

func MakeCorpusTabUI(fc *FyneCorpus, w fyne.Window) fyne.CanvasObject {

	bo := makeNewBoundCorpusOutputs()
	ou := makeNewCorpusOutputUI(bo)
	bo.corpusCountBV.Set(fc.c.GetWordCount())

	radio := widget.NewRadioGroup([]string{"Words", "Lines" /*, "Sentences", "Phrases"*/}, func(value string) {
		bo.importedMethodStr.Set(value)
	})
	radio.SetSelected("Words")

	importFile := func(reader fyne.URIReadCloser) (int, error) {
		uri := reader.URI()
		if uri.MimeType() == "text/plain" {
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
				for scanner.Scan() {
					//fmt.Println("scanned", scanner.Text())
					s, _ := fc.CleanString(scanner.Text())
					_ = fc.Add(s)
					i += 1
				}
			}
			return i, er
		}
		return 0, errors.New("input file must be text/plain")

	}

	fileImportButton := widget.NewButton("Import File", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil {
				if reader != nil {
					i, er := importFile(reader)
					if er == nil {
						bo.importedBoundValue.Set(i)
						bo.corpusCountBV.Set(fc.c.GetWordCount())
					} else {
						dialog.ShowError(er, w)
					}
				}
			} else {
				dialog.ShowError(err, w)
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
				bo.corpusCountBV.Set(fc.c.GetWordCount())
			}
		}, w)
	})

	/*	annotateButton := widget.NewButton("Annotate file", func() {
			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil {
					if reader != nil {
						uri := reader.URI()
						if uri.MimeType() == "text/plain" {
							var rstring strings.Builder
							scanner := bufio.NewScanner(reader)
							scanner.Split(bufio.ScanWords)
							for scanner.Scan() {
								s, _ := fc.CleanString(scanner.Text())
								_ = fc.Add(s)
								v, _ := fc.base.StringValue(s)
								rstring.WriteString(s + "^" + strconv.Itoa(v) + "^ ")
							}
							md := widget.NewRichTextFromMarkdown(rstring.String())
							md.Wrapping = fyne.TextWrapWord
							c := container.NewBorder(nil, nil, nil, nil, md)
							fmt.Println("gonna show the md")
							d := dialog.NewCustom("Annotated Text", "OK", c, w)
							d.Show()
						}
					}
				} else {
					dialog.ShowError(err, w)
				}
			}, w)
		})
		annotateButton.Disable()*/

	exportRadio := widget.NewRadioGroup([]string{ /*"Alphabetically",*/ "Numerically" /*, "Sentences", "Phrases"*/}, func(value string) {
		bo.exportMethodStr.Set(value)
	})
	exportRadio.SetSelected("Numerically")
	exportRadio.Disable()

	exportButton := widget.NewButton("Export Corpus", func() {
		d := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
			if err == nil && closer != nil {
				fc.c.SaveNumerically(closer)
			}
		}, w)
		d.Show()
	})

	clearButton := widget.NewButton("Clear Corpus", func() {
		d := dialog.NewConfirm("Clear Corpus?", "Clear the corpus? This cannot be undone!", func(doIt bool) {
			if doIt {
				fc.Clear()
			}
		}, w)
		d.Show()
	})

	importArea := container.NewVBox(widget.NewLabel("Import"), widget.NewSeparator(),
		container.NewVBox(container.NewHBox(radio, folderImportButton, fileImportButton)), widget.NewSeparator(),
		/*container.NewHBox(annotateButton, widget.NewLabel("Import a text file. Add the items to the corpus,\n then show an annotated version of the text"))*/
	)
	//outputArea := makeCorpusOutputArea()
	exportArea := container.NewVBox(widget.NewLabel("Export"), widget.NewSeparator(), container.NewHBox(exportRadio, exportButton))
	clearArea := container.NewVBox(widget.NewLabel("Clear"), widget.NewSeparator(), container.NewHBox(clearButton))
	corpusArea := container.NewVBox(widget.NewLabel("Corpus Management"), widget.NewSeparator(), ou)
	border := container.NewBorder(corpusArea, nil, nil, nil, container.NewVBox(exportArea, importArea, clearArea))

	return border
}
