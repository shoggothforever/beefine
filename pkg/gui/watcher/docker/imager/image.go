package imager

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"strconv"
	"time"
)

const PKGName = "imager"

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
	ImageLogs := widget.NewTextGrid()
	ImageLogs.SetText("Container-Creating Logs")
	ImageLogsScroll := container.NewScroll(ImageLogs)
	ImageLogsScroll.SetMinSize(fyne.NewSize(400, 200)) // 限制宽度为 400，高度为 200

	//var bpfChoices
	toolbar := NewToolBar(ImageLogs, bpfLogs)
	// 模拟实时更新日志和图表数据
	//go simulateImageData(bpfLogs)
	content := container.NewHBox(
		toolbar,
		ImageLogsScroll,
		bpfLogsScroll,
	)
	return content
}

// simulateImageData 模拟实时数据更新
func simulateImageData(logs *widget.TextGrid) {
	for {
		// 模拟随机更新镜像层数据
		layer := rand.Intn(4) // 4 层镜像层
		value := rand.Intn(20) + 10
		laystr := strconv.Itoa(layer)
		valstr := strconv.Itoa(value)
		// 更新日志
		logEntry := time.Now().Format("15:04:05") + " - Layer " + laystr + ": Wrote " + valstr + "MB"
		logs.SetRow(len(logs.Rows), widget.NewTextGridFromString(logEntry).Row(0))
		// 模拟滚动：将光标移动到文本末尾
		//logs.CursorRow = len(logs.Text)
		logs.Refresh()
		time.Sleep(2 * time.Second)
	}
}
