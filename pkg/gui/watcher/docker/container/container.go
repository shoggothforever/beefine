package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const PKGName = "container"

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("Docker Container Monitoring"),
		widget.NewLabel("Monitor container performance"),
		// 在这里加入性能指标图表或实时监控视图
	)
}
