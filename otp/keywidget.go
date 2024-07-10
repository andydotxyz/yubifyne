//go:generate fyne bundle -p otp key.svg -o bundled.go
package otp

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Input struct {
	widget.BaseWidget
	outline *canvas.Rectangle

	text        string
	onSubmitted func(string)
}

func NewInput(done func(string)) *Input {
	k := &Input{outline: canvas.NewRectangle(color.Transparent), onSubmitted: done}
	k.ExtendBaseWidget(k)

	k.outline.StrokeWidth = theme.InputBorderSize()
	k.outline.StrokeColor = theme.InputBorderColor()
	k.outline.CornerRadius = theme.InputRadiusSize()
	return k
}

func (k *Input) CreateRenderer() fyne.WidgetRenderer {
	labelHeight := widget.NewLabel("").MinSize().Height
	prop := canvas.NewRectangle(color.Transparent)
	prop.SetMinSize(fyne.NewSquareSize(labelHeight))

	r := theme.NewThemedResource(resourceKeySvg)
	return widget.NewSimpleRenderer(
		container.NewHBox(container.NewStack(k.outline,
			prop, container.NewPadded(canvas.NewImageFromResource(r)))))
}

func (k *Input) FocusGained() {
	k.text = ""
	k.outline.StrokeColor = theme.PrimaryColor()
	k.outline.Refresh()
}

func (k *Input) FocusLost() {
	k.text = ""
	k.outline.StrokeColor = theme.InputBorderColor()
	k.outline.Refresh()
}

func (k *Input) TypedRune(r rune) {
	k.text = k.text + string([]rune{r})
}

func (k *Input) TypedKey(ev *fyne.KeyEvent) {
	if ev.Name != fyne.KeyEnter && ev.Name != fyne.KeyReturn {
		return
	}

	k.onSubmitted(k.text)
}

func (k *Input) Tapped(_ *fyne.PointEvent) {
	fyne.CurrentApp().Driver().CanvasForObject(k).Focus(k)
}
