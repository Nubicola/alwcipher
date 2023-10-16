package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/corpus"
)

func makeToolbar(c *corpus.Corpus) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
	return toolbar
}

// should "show matches from"
func MakeUI(c *corpus.Corpus) fyne.CanvasObject {

	tabs := container.NewAppTabs(
		container.NewTabItem("Words", MakeWordsTabUI(c)),
		container.NewTabItem("Numbers", MakeNumbersTabUI(c)),
		container.NewTabItem("Corpus", widget.NewLabel("The corpus")),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationBottom)

	content := container.NewBorder(makeToolbar(c), nil, nil, nil, tabs)
	return content
}
