package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/corpus"
)

func makeToolbar(c *corpus.Corpus, w fyne.Window) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
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
					}
				}
			}, w)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			richtext := widget.NewRichTextFromMarkdown(
				"Thank you for using ALW Visualized\n\n" +
					"# Identify\n\n Type in a word or words. The EQ value will be shown on the right. Matching values from your corpus will be shown.\n\n" +
					"The selection box 'save to corpus' will do just that.\n\n" +
					"# Explore\n\n Type in a number. All matching words in the corpus will be shown\n\n" +
					"# Discover\n\n Import words or lines from an external file. Every word or line will undergo the EQ calculation and be added to the existing corpus\n\n" +
					"# Load and Save\n\n Save and load a corpus file in this app's native format\n\n")
			richtext.Wrapping = fyne.TextWrapWord
			d := dialog.NewCustom("ALW Help", "OK", richtext, w)
			d.Show()
			width, height := w.Canvas().Size().Components()
			d.Resize(w.Canvas().Size().SubtractWidthHeight(width*.2, height*.4))
		}),
	)
	return toolbar
}

// should "show matches from"
func MakeUI(c *corpus.Corpus, w fyne.Window) fyne.CanvasObject {

	tabs := container.NewAppTabs(
		container.NewTabItem("Words", MakeWordsTabUI(c)),
		container.NewTabItem("Numbers", MakeNumbersTabUI(c)),
		container.NewTabItem("Import", MakeCorpusTabUI(c, w)),
	)

	tabs.SetTabLocation(container.TabLocationBottom)
	content := container.NewBorder(makeToolbar(c, w), nil, nil, nil, tabs)
	return content
}
