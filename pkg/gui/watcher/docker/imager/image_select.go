package imager

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// MyCustomWidget 是自定义控件，包装了 Select 并添加了额外的字段
type ImageSelect struct {
	widget.BaseWidget                // 嵌入 BaseWidget
	base              *widget.Select // 内嵌 Select
	currentImage      string         // 额外的字段
}

// NewMyCustomWidget 创建自定义控件实例
func NewMyCustomWidget(options []string, onSelected func(string)) *ImageSelect {
	// 初始化 Select
	selectWidget := widget.NewSelect(options, onSelected)

	// 创建 MyCustomWidget 实例
	widget := &ImageSelect{
		base:         selectWidget,
		currentImage: "default data", // 设置默认值
	}

	widget.ExtendBaseWidget(widget) // 必须扩展 BaseWidget
	return widget
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (w *ImageSelect) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(w.base)
}
