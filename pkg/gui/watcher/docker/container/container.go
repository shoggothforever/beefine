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
	// 容器选择
	containerSelector := widget.NewSelect([]string{"Container 1", "Container 2"}, func(value string) {
		// TODO: 更新对应容器的数据
	})
	containerSelector.SetSelected("Select a container")

	// 性能指标选项卡
	cpuUsage := widget.NewLabel("CPU Usage Chart Placeholder")
	memoryUsage := widget.NewLabel("Memory Usage Chart Placeholder")
	networkActivity := widget.NewLabel("Network Activity Chart Placeholder")

	tabs := container.NewAppTabs(
		container.NewTabItem("CPU", cpuUsage),
		container.NewTabItem("Memory", memoryUsage),
		container.NewTabItem("Network", networkActivity),
	)

	// 动态日志区域
	logs := widget.NewMultiLineEntry()
	logs.SetPlaceHolder("Container runtime logs...")
	logs.Disable()

	content := container.NewStack(
		containerSelector,
		tabs,
		widget.NewSeparator(),
		logs,
	)

	return content
}
