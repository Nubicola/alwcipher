package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Nubicola/alwcipher/pkg/calculator"
	"github.com/Nubicola/alwcipher/pkg/corpus"
	"github.com/mozillazg/go-unidecode"
)

func Ints(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

type FyneCorpus struct {
	c                 *corpus.Corpus
	base, first, last calculator.ALWCalculator
}

func (c *FyneCorpus) CleanString(s string) (string, error) {
	sdeunicoded := unidecode.Unidecode(s)
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err == nil {
		processedString := reg.ReplaceAllString(sdeunicoded, "")
		return strings.ToUpper(processedString), nil
	}
	return "", err
}

func (c *FyneCorpus) Add(s string) []string {
	var ss = make([]string, 1)
	// calculate the value of that thing
	val, verr := c.c.Calculate(s)
	if verr != nil {
		dialog.ShowError(verr, nil)
		fmt.Println(verr)
	} else {
		if fyne.CurrentApp().Preferences().Bool("WriteToCorpus") {
			// add it to the corpus
			err := c.c.Add(s)
			if err != nil {
				dialog.ShowError(err, nil)
			} else {
				// add it to the preferences
				ss = c.c.Get(val)
				sval := strconv.Itoa(val)
				fyne.CurrentApp().Preferences().SetStringList(sval, ss)

				ints := fyne.CurrentApp().Preferences().IntList("CorpusKeys")
				ints = append(ints, val)
				// store the corpus-key-list in one setting, and the corpus-value-list directly
				fyne.CurrentApp().Preferences().SetIntList("CorpusKeys", Ints(ints))

			}
		} else {
			// just add it to the return value
			ss = append(c.c.Get(val), s)
		}
	}
	return ss
}

func (c *FyneCorpus) Clear() {
	c.c.Clear()
	ints := fyne.CurrentApp().Preferences().IntList("CorpusKeys")
	for _, key := range ints {
		skey := strconv.Itoa(key)
		fyne.CurrentApp().Preferences().RemoveValue(skey)
	}
	fyne.CurrentApp().Preferences().RemoveValue("CorpusKeys")
}

func makeToolbar(fc FyneCorpus, w fyne.Window) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		/*widget.NewToolbarAction(theme.UploadIcon(), func() {
			d := dialog.NewCustom("Corpus Management", "OK", MakeCorpusTabUI(&fc, w), w)
			d.Show()
		}),
		widget.NewToolbarSpacer(),*/
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			richtext := widget.NewRichTextFromMarkdown(
				"Enjoying this application? \n\n I'd love it if you would [buy me a coffee!](https://www.buymeacoffee.com/nubicola) " +
					"\n\n The purpose of the utility is to help you discover relationships between words using English Gematria. In particular, you can build" +
					" a 'corpus', or a body of word-to-number (or phrase-to-number) associations." +
					"\n\n **NOTE** don't load huge text files. Later releases will support larger amounts of words. \n\n" +
					"\n\n Try it. Type in 'hello' and press enter. Type in 'helol' and see that it has the same value. Change to the other tab and type" +
					" '40' and see that both words are there. Import entire text files by line or by word. Take care: there is not a lot of error handling." +
					"\n\n# Words\n\n Type in a word or words. The EQ value will be shown on the right. Matching values from your corpus will be shown.\n\n" +
					"The selection box 'save to corpus' will do just that.\n\n" +
					"# Numbers\n\n Type in a number. All matching words in the corpus will be shown\n\n" +
					"# Corpus\n\n Operate on the corpus itself; you can import whole text files, export, and clear it out completely.\n\n")
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

// should "show matches from"
func MakeUI(w fyne.Window) fyne.CanvasObject {

	fyneCorpus := new(FyneCorpus)

	fyneCorpus.base = new(calculator.EQBaseCalculator)
	fyneCorpus.first = new(calculator.EQFirstCalculator)
	fyneCorpus.last = new(calculator.EQLastCalculator)
	fyneCorpus.c = corpus.NewCorpus(fyneCorpus.base)

	tabs := container.NewAppTabs(
		container.NewTabItem("Words", MakeWordsTabUI(fyneCorpus)),
		container.NewTabItem("Numbers", MakeNumbersTabUI(fyneCorpus)),
		container.NewTabItem("Corpus", MakeCorpusTabUI(fyneCorpus, w)),
	)

	// load the corpus stored in the preferences file into the corpus object
	ints := fyne.CurrentApp().Preferences().IntList("CorpusKeys")
	for _, key := range ints {
		skey := strconv.Itoa(key)
		corpusList := fyne.CurrentApp().Preferences().StringList(skey)
		for _, s := range corpusList {
			_ = fyneCorpus.c.Add(s)
		}
	}
	tabs.SetTabLocation(container.TabLocationBottom)

	content := container.NewBorder(makeToolbar(*fyneCorpus, w), nil, nil, nil, tabs)
	return content
}
