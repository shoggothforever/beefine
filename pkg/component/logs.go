package component

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"shoggothforever/beefine/internal/cli"
	"strings"
	"sync"
)

const MaxLogRow = 200

// LogBoard 采用TextGrid性能不高，暂时采用entry替代
type LogBoard struct {
	widget.BaseWidget
	lines          int
	name           string
	entry          *widget.Entry
	m              sync.Mutex
	scroll         *container.Scroll
	cleanButton    *widget.Button
	analysisButton *widget.Button

	defaultText string
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (l *LogBoard) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	lg := container.NewVBox(l.scroll, l.cleanButton, l.analysisButton)
	return widget.NewSimpleRenderer(lg)
}
func (l *LogBoard) AppendLogf(format string, args ...any) {
	text := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	l.entry.Append(text)
	l.m.Lock()
	l.lines++
	if l.lines%MaxLogRow == 0 {
		l.Clear()
		//str := l.entry.Text
		//strs := strings.SplitAfterN(str, "\n", MaxLogRow/2)
		//if len(strs) == 0 {
		//	l.SetText("")
		//} else {
		//	l.SetText(strs[len(strs)-1])
		//}
		l.lines = 0
	}
	l.m.Unlock()
	//l.m.Lock()
	//defer l.m.Unlock()
	//l.logs.SetRow(len(l.logs.Rows), widget.NewTextGridFromString(text).Row(0))
	//if len(l.logs.Rows) > MaxLogRow {
	//	l.Clear()
	//}
	l.entry.DragEnd()
	l.Refresh()
}
func (l *LogBoard) SetText(text string) {
	l.entry.SetText(text)
}
func (l *LogBoard) Clear() {
	l.entry.SetText("")
	l.lines = 0
	l.Refresh()
}
func (l *LogBoard) Content() string {
	return l.entry.Text
}
func (l *LogBoard) WriteAnalysisLog(filePath, str string) {
	// 以追加模式打开文件（如果不存在则创建）
	// os.O_APPEND - 追加模式
	// os.O_CREATE - 如果文件不存在则创建
	// os.O_WRONLY - 只写模式
	// 0644 - 文件权限（rw-r--r--）
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	defer file.Close() // 确保文件最终被关闭
	// 写入文本
	_, err = file.WriteString(str + "\n") // 添加换行符
	if err != nil {
		log.Fatalf("写入文件失败: %v", err)
	}
	// 可选：同步写入磁盘（确保数据持久化）
	err = file.Sync()
	if err != nil {
		log.Printf("同步文件失败: %v", err)
	}
	log.Println("文本已成功追加到文件")
}
func NewLogBoard(name, text string, weight, height float32) *LogBoard {
	if len(text) == 0 {
		text = "\n"
	}
	entry := widget.NewMultiLineEntry()
	entry.Append(text)
	boardScroll := container.NewScroll(entry)
	boardScroll.SetMinSize(fyne.NewSize(weight, height))
	l := &LogBoard{
		name:        name,
		lines:       0,
		defaultText: text,
		entry:       entry,
		m:           sync.Mutex{},
		scroll:      boardScroll,
	}
	l.cleanButton = widget.NewButton("clear", func() { l.Clear() })
	l.analysisButton = widget.NewButton("analysis", func() {
		go func() {
			l.AppendLogf("start analysising logs\n")
			fmt.Println("analysising ", l.Content())
			rsp := cli.Analysis(l.Content())
			fmt.Println("analysis success")
			//l.AppendLogf("%s\n", rsp)
			l.WriteAnalysisLog("./tmp/"+l.name+"analysis.log", rsp)
		}()
	})
	l.ExtendBaseWidget(l)
	return l
}
