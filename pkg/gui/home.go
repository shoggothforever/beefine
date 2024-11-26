package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// HomeUI  creates the Home page with a list of available features.
func HomeUI(w fyne.Window) fyne.CanvasObject {
	// 按钮：List Probes
	listProbesButton := NewUITabButton(ProbesUIName, ProbesUI)
	// 按钮：counter Probes
	CounterButton := NewUITabButton(CounterUIName, CounterUI)
	// 按钮：tracePoint Exec
	ExecButton := NewUITabButton(ExecUIName, ExecUI)
	// 其他功能按钮
	placeholderButton := widget.NewButton("Placeholder Feature", func() {
		widget.ShowPopUp(widget.NewLabel("Feature coming soon!"), w.Canvas())
	})

	// 布局
	return container.NewGridWithColumns(4,
		widget.NewLabel("Select a feature:"),
		listProbesButton,
		CounterButton,
		ExecButton,
		placeholderButton,
	)
}
