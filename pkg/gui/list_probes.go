package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/logic"
)

const ProbesUIName = "ListProbes"

// ProbesUI creates the UI for the List Probes feature.
func ProbesUI() fyne.CanvasObject {
	// 获取 hooks 数据
	progTypeToHooks := logic.ListProbes()

	// 下拉菜单：progtype
	progTypeSelect := widget.NewSelect([]string{}, nil)
	progTypeSelect.PlaceHolder = "Select ProgType"

	// 下拉菜单：hook
	hookSelect := widget.NewSelect([]string{}, nil)
	hookSelect.PlaceHolder = "Select Hook"

	// 初始化 ProgType 选项
	for progType := range progTypeToHooks {
		progTypeSelect.Options = append(progTypeSelect.Options, progType)
	}
	progTypeSelect.Refresh()

	// 监听 ProgType 的选择事件
	progTypeSelect.OnChanged = func(progType string) {
		// 更新 hookSelect 的内容
		hooks, exists := progTypeToHooks[progType]
		if !exists {
			hookSelect.Options = []string{}
		} else {
			hookSelect.Options = hooks
		}
		hookSelect.Refresh()
		hookSelect.SetSelected("") // 重置选择
	}

	// 布局
	return NewUIVBox(
		ProbesUIName,
		widget.NewLabel("Select a ProgType and Hook:"),
		progTypeSelect,
		hookSelect,
	)
}
