package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sync"
)

type LogBoard struct {
	widget.BaseWidget
	logs   *widget.TextGrid
	scroll *container.Scroll
	m      sync.Mutex
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (l *LogBoard) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(l.scroll)
}
func (l *LogBoard) AppendLog(text string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.logs.SetRow(len(l.logs.Rows), widget.NewTextGridFromString(text).Row(0))
	if len(l.logs.Rows) > 200 {
		l.logs.SetText("")
	}
}
func NewLogBoard(text string, weight, height float32) *LogBoard {
	boardLog := widget.NewTextGrid()
	boardLog.ShowLineNumbers = true
	boardLog.SetText(text)
	boardScroll := container.NewScroll(boardLog)
	boardScroll.SetMinSize(fyne.NewSize(weight, height)) // 限制宽度为 600，高度为 200
	l := &LogBoard{
		logs:   boardLog,
		scroll: boardScroll,
		m:      sync.Mutex{},
	}
	l.ExtendBaseWidget(l)
	return l
}
