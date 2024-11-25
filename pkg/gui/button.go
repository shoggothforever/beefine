package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type NewCanvasObjectFunc func() fyne.CanvasObject

func NewUITabButton(title string, fn NewCanvasObjectFunc) *widget.Button {
	return widget.NewButton(title, func() {
		// 打开一个新选项卡显示 List Probes UI
		tabItem := tabManager.Get(title)
		if tabItem == nil {
			tabItem = container.NewTabItem(title, fn())
			tabManager.Append(tabItem, true)
		} else {
			tabManager.Select(tabItem)
		}
	})
}

const CloseButtonTitle = "Close Tab"

func NewCloseButton(title string, cancel func()) fyne.CanvasObject {
	// 关闭按钮
	closeButton := widget.NewButtonWithIcon(CloseButtonTitle, theme.WindowCloseIcon(), func() {
		// 从 Tabs 中移除当前 Tab
		if item, ok := tabManager.TabItemsMap[title]; ok && item != nil {
			tabManager.Remove(item)
		}
		if cancel != nil {
			cancel()
		}
	})
	closeButton.Resize(fyne.NewSize(150, 40)) // 固定宽度和高度
	closeButton.Move(fyne.NewPos(300, 50))    // 设置位置
	content := container.NewHBox(
		closeButton,
	)
	return content
}

// NewClosableTab creates a new TabItem with a close button on the Tab title.
func NewClosableTab(title string, content fyne.CanvasObject) *container.TabItem {
	// 创建标题容器，包含标题文字和关闭按钮
	label := widget.NewLabel(title)
	closeButton := NewCloseButton(title, nil)
	// 使用 HBox 布局将标题和关闭按钮放在一起
	tabTitle := container.NewHBox(label, content, closeButton)
	// 创建透明的事件捕获区域
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(label.MinSize())

	// 创建 TabItem
	return container.NewTabItemWithIcon(title, nil, tabTitle)
}
