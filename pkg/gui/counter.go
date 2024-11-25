package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const CounterUIName = "countNetPackage"

func CounterUI() *fyne.Container {
	return NewUIVBox(
		CounterUIName,
		widget.NewLabel("click to check how many net package have been received"),
	)
}
