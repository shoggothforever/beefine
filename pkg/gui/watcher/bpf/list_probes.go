package bpf

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
	"shoggothforever/beefine/pkg/gui/watcher/bpf/logic"
	"strings"
)

const HelpersUIName = "ListHelpFunction"

// HelpersUI creates the UI for the List Probes feature.
func HelpersUI() fyne.CanvasObject {
	// 获取 hooks 数据
	progTypeToHooks := logic.ListHelperFunc()
	// 下拉菜单：progtype
	progTypeSelect := widget.NewSelect([]string{}, nil)
	progTypeSelect.PlaceHolder = "Select ProgType"

	// 初始化 ProgType 选项
	for progType := range progTypeToHooks {
		progTypeSelect.Options = append(progTypeSelect.Options, progType)
	}
	progTypeSelect.Refresh()
	// 可用hook看板
	log := component.NewLogBoard("eBPF helpers supported for program type", 200, 400)
	// 监听 ProgType 的选择事件
	progTypeSelect.OnChanged = func(progType string) {
		// 更新 hookSelect 的内容
		hooks, exists := progTypeToHooks[progType]
		if !exists {
			log.SetText("\n")
		} else {
			log.SetText(strings.Join(hooks, "\n"))
		}
	}
	// 布局
	return component.NewUIVBox(
		PKGName,
		HelpersUIName,
		nil,
		widget.NewLabel("get eBPF helpers function"),
		widget.NewLabel("Select a ProgType and Hook:"),
		progTypeSelect,
		log,
	)
}
