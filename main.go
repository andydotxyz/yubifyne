package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	loadTheme(a)
	w := a.NewWindow("YubiFyne")

	g := newGUI(w)
	w.SetContent(g.makeUI())

	w.Resize(fyne.NewSize(240, 168))
	g.win.Canvas().Focus(g.user)

	w.ShowAndRun()
}
