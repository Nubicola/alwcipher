package ui

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/corpus"
)

func makeToolbar(c *corpus.Corpus, w fyne.Window) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			/*	fyne.CurrentApp().Preferences().SetBool("WriteToCorpus", value)
				writeToCorpusCheckbox.Checked = fyne.CurrentApp().Preferences().BoolWithFallback("WriteToCorpus", true)*/
			/* corpusDir := fyne.CurrentApp().Preferences().StringWithFallback("CorpusDir", os.UserHomeDir())
			corpusName := fyne.CurrentApp().Preferences().StringWithFallback("CorpusDir", "ALWCorpus.json")
			uriString, uriErr := storage.ParseURI()*/

			dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader != nil {
					err = c.LoadNative(reader)
					if err != nil {
						dialog.ShowError(err, w)
					}
				}
			}, w)
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err == nil && writer != nil {
					err = c.SaveNative(writer)
					if err != nil {
						dialog.ShowError(err, w)
					} else {
						// store save location in preferences so it will be the default next time
						fyne.CurrentApp().Preferences().SetString("CorpusURI", writer.URI().String())
					}

				}
			}, w)
		}),
		widget.NewToolbarAction(theme.UploadIcon(), func() {
			d := dialog.NewCustom("Import from file", "OK", MakeCorpusTabUI(c, w), w)
			d.Show()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			richtext := widget.NewRichTextFromMarkdown(
				"The purpose of the utility is to help you discover relationships between words using English Gematria. In particular, you can build" +
					" a 'corpus', or a body of word-to-number (or phrase-to-number) associations." +
					"\n\n Try it. Type in 'hello' and press enter. Type in 'helol' and see that it has the same value. Change to the other tab and type" +
					" '40' and see that both words are there. Import entire text files by line or by word. Take care: there is not a lot of error handling." +
					"\n\n# Words\n\n Type in a word or words. The EQ value will be shown on the right. Matching values from your corpus will be shown.\n\n" +
					"The selection box 'save to corpus' will do just that.\n\n" +
					"# Numbers\n\n Type in a number. All matching words in the corpus will be shown\n\n" +
					"# Load and Save\n\n Save and load a corpus file in this app's native format\n\n" +
					"# Import\n\n Import words or lines from an **external** file. Every word or line will undergo the EQ calculation and be added to the existing corpus\n\n")
			richtext.Wrapping = fyne.TextWrapWord
			d := dialog.NewCustom("ALW Corpus Explorer Help", "OK", richtext, w)
			d.Show()
			width, height := w.Canvas().Size().Components()
			d.Resize(w.Canvas().Size().SubtractWidthHeight(width*.2, height*.2))
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			d := dialog.NewCustom("Credits", "OK", CreditsContainer(), w)
			d.Resize(fyne.NewSize(800, 400))
			d.Show()
		}),
	)
	return toolbar
}

func loadLastSavedCorpus(c *corpus.Corpus) error {
	corpusName := fyne.CurrentApp().Preferences().StringWithFallback("CorpusURI", fyne.CurrentApp().Storage().RootURI().String()+"/corpus.json")
	// try to load the file
	uri, errURI := storage.ParseURI(corpusName)
	if errURI != nil {
		return errURI
	}
	exists, _ := storage.Exists(uri)
	if !exists {
		return nil
	}
	reader, err := storage.Reader(uri)
	if err != nil {
		return err
	}
	defer reader.Close()
	return c.LoadNative(reader)
}

// should "show matches from"
func MakeUI(c *corpus.Corpus, w fyne.Window) fyne.CanvasObject {

	err := loadLastSavedCorpus(c)
	if err != nil {
		err = errors.Join(errors.New("unable to load previous corpus. Removing saved preference"), err)
		dialog.ShowError(err, w)
		fyne.CurrentApp().Preferences().RemoveValue("CorpusURI")
	}
	tabs := container.NewAppTabs(
		container.NewTabItem("Words", MakeWordsTabUI(c)),
		container.NewTabItem("Numbers", MakeNumbersTabUI(c)),
		//container.NewTabItem("Import", MakeCorpusTabUI(c, w)),
	)

	tabs.SetTabLocation(container.TabLocationBottom)
	content := container.NewBorder(makeToolbar(c, w), nil, nil, nil, tabs)
	return content
}
