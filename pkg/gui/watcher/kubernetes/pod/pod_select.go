package pod

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
	"sync"
)

// PodSelect 是自定义控件，包装了 Select 并添加了额外的字段
type PodSelect struct {
	widget.BaseWidget                // 嵌入 BaseWidget
	base              *widget.Select // 内嵌 Select
	podLogs           *component.LogBoard
	bpfLogs           *component.LogBoard
	cancelMap         map[string]func()
	watcherID         int //观测的docker 容器id
	m                 sync.Mutex
}

// NewImageSelect 创建自定义控件实例
func NewPodSelect() *PodSelect {
	// 初始化 Select
	selectWidget := widget.NewSelect(nil, nil)
	// 创建 MyCustomWidget 实例
	s := &PodSelect{
		base:      selectWidget,
		cancelMap: make(map[string]func()),
		m:         sync.Mutex{},
	}
	//s.base.OnChanged = s.OnChanged
	s.base.PlaceHolder = "select existed pod"
	s.ExtendBaseWidget(s) // 必须扩展 BaseWidget
	return s
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (w *PodSelect) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(w.base)
}
func (w *PodSelect) chooseNetwork(b bool) {
	if b {
		w.podLogs.AppendLogf("tracing pod communication networks")

	} else {

	}

}
