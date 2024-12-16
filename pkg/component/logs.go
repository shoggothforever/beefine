package component

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
	"sync"
)

const MaxLogRow = 500

type LogBoard struct {
	widget.BaseWidget
	entry       *widget.Entry
	logs        *widget.TextGrid
	scroll      *container.Scroll
	defaultText string
	m           sync.Mutex
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (l *LogBoard) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(l.scroll)
}
func (l *LogBoard) AppendLogf(format string, args ...any) {
	l.m.Lock()
	defer l.m.Unlock()
	text := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	l.entry.Append(text)
	//l.logs.SetRow(len(l.logs.Rows), widget.NewTextGridFromString(text).Row(0))
	if len(l.logs.Rows) > MaxLogRow {
		l.logs.SetText("")
	}
	l.Refresh()
}
func (l *LogBoard) SetText(text string) {
	l.m.Lock()
	defer l.m.Unlock()
	l.entry.SetText(text)
}
func (l *LogBoard) Clear() {
	l.m.Lock()
	defer l.m.Unlock()
	l.logs.SetText(l.defaultText)
	l.Refresh()
}
func NewLogBoard(text string, weight, height float32) *LogBoard {
	boardLog := widget.NewTextGrid()
	boardLog.ShowLineNumbers = true
	if len(text) == 0 {
		text = "\n"
	}
	boardLog.SetText(text)
	entry := widget.NewEntry()
	entry.SetText(text)
	entry.Append(text)
	boardScroll := container.NewScroll(entry)
	boardScroll.SetMinSize(fyne.NewSize(weight, height))
	l := &LogBoard{
		defaultText: text,
		entry:       entry,
		logs:        boardLog,
		scroll:      boardScroll,
		m:           sync.Mutex{},
	}
	l.ExtendBaseWidget(l)
	return l
}
