package ui

import (
	"log"

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
			log.Println("Display help")
		}),
	)
	return toolbar
}

// should "show matches from"
func MakeUI(c *corpus.Corpus, w fyne.Window) fyne.CanvasObject {

	tabs := container.NewAppTabs(
		container.NewTabItem("Identify", MakeWordsTabUI(c)),
		container.NewTabItem("Explore", MakeNumbersTabUI(c)),
		container.NewTabItem("Discover", MakeCorpusTabUI(c, w)),
	)

	tabs.SetTabLocation(container.TabLocationBottom)
	content := container.NewBorder(makeToolbar(c, w), nil, nil, nil, tabs)
	return content
}
