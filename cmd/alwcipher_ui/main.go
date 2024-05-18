package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	ui "github.com/Nubicola/alwcipher/internal/ui"
)

// foreground: BBBBBBFF
// disabled: BBBBBB42
type alwTheme struct{}

func (alwTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameDisabled:
		// changing the disabled color to the regular color because I find the disabled one hard to read
		return theme.DefaultTheme().Color(theme.ColorNameForeground, v)
	default:
		return theme.DefaultTheme().Color(c, v)
	}
}

func (alwTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (alwTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (alwTheme) Size(s fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(s)
}

func main() {
	myApp := app.NewWithID("eu.thesaturdays.alwcipher")
	myWindow := myApp.NewWindow("ALW Corpus Explorer")
	myWindow.Resize(fyne.Size{
		Width:  800,
		Height: 1200,
	})

	myApp.Settings().SetTheme(&alwTheme{})
	myWindow.SetContent(ui.MakeUI(myWindow))
	myWindow.ShowAndRun()
}
