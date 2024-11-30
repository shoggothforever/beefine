package docker

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
)

const PKGName = "docker"

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	// 创建选项卡容器
	tabM := component.NewTabManager(PKGName, &w)
	if homeItem := tabM.Get("home"); homeItem != nil {
		// 布局
		tabM.Select(homeItem)
		return tabM.Tabs
	}
	objects := []fyne.CanvasObject{widget.NewLabel("Select a feature:")}
	for name, fn := range tabUIButtonFuncMap {
		objects = append(objects, component.NewUITabButton(PKGName, name, fn))
	}
	// 其他功能按钮
	placeholderButton := widget.NewButton("Placeholder Feature", func() {
		widget.ShowPopUp(widget.NewLabel("Feature coming soon!"), w.Canvas())
	})
	objects = append(objects, placeholderButton)
	home := container.NewGridWithColumns(4,
		objects...,
	)
	homeItem := container.NewTabItem("home", home)
	tabM.Append(homeItem, true)
	// 布局
	tabM.Select(homeItem)
	return tabM.Tabs
}
