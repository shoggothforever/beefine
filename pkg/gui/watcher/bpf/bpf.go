package bpf

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
	themes2 "shoggothforever/beefine/pkg/gui/themes"
)

const PKGName = "bpf"

var TabUIButtonFuncMap = map[string]func() fyne.CanvasObject{
	ProbesUIName:  ProbesUI,
	CounterUIName: CounterUI,
	ExecUIName:    ExecUI,
}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	a := fyne.CurrentApp()
	objects := []fyne.CanvasObject{widget.NewLabel("Select a feature:")}
	for name, fn := range TabUIButtonFuncMap {
		objects = append(objects, component.NewUITabButton(PKGName, name, fn))
	}
	// 其他功能按钮
	placeholderButton := widget.NewButton("Placeholder Feature", func() {
		widget.ShowPopUp(widget.NewLabel("Feature coming soon!"), w.Canvas())
	})
	objects = append(objects, placeholderButton)

	themes := themes2.CreateThemes(a)
	objects = append(objects, themes)

	// 布局
	return container.NewGridWithColumns(4,
		objects...,
	)
}
