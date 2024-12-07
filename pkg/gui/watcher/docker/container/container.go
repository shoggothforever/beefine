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
	// 动态bpf日志区域
	bpfLogs := widget.NewTextGrid()
	bpfLogs.ShowLineNumbers = true
	bpfLogs.SetText("Real-Time BpfLogs")
	//bpfLogs.
	bpfLogsScroll := container.NewScroll(bpfLogs)
	bpfLogsScroll.SetMinSize(fyne.NewSize(400, 200)) // 限制宽度为 400，高度为 200
	// image 日志
	containerLogs := widget.NewTextGrid()
	containerLogs.SetText("Container-Creating Logs")
	ImageLogsScroll := container.NewScroll(containerLogs)
	ImageLogsScroll.SetMinSize(fyne.NewSize(400, 200)) // 限制宽度为 400，高度为 200

	//var bpfChoices
	toolbar := NewContainerToolBar(containerLogs, bpfLogs)
	content := container.NewHBox(
		toolbar,
		ImageLogsScroll,
		bpfLogsScroll,
	)
	return content
}
