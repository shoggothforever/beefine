package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// HomeUI  creates the Home page with a list of available features.
func HomeUI(w fyne.Window) fyne.CanvasObject {
	// 按钮：List Probes
	listProbesButton := NewUITabButton(ProbesUIName)
	// 其他功能按钮
	placeholderButton := widget.NewButton("Placeholder Feature", func() {
		widget.ShowPopUp(widget.NewLabel("Feature coming soon!"), w.Canvas())
	})

	// 布局
	return container.NewVBox(
		widget.NewLabel("Select a feature:"),
		listProbesButton,
		placeholderButton,
	)
}
