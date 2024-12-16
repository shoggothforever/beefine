package bpf

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
)

const PKGName = "bpf"

type UITab struct {
	name   string
	uiFunc func() fyne.CanvasObject
}

var tabUIButtonFuncGroup = []UITab{
	{HelpersUIName, HelpersUI},
	{CounterUIName, CounterUI},
	{ExecUIName, ExecUI},
}

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
	for _, v := range tabUIButtonFuncGroup {
		objects = append(objects, component.NewUITabButton(PKGName, v.name, v.uiFunc))
	}
	// 其他功能按钮
	placeholderButton := placeHolder(w)
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
func placeHolder(w fyne.Window) fyne.CanvasObject {
	return widget.NewButton("Placeholder Feature", func() {
		dialog.ShowInformation("placeHolder", "Feature coming soon!", w)
	})
}
