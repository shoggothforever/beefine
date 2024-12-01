package imager

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"math/rand"
	"shoggothforever/beefine/internal/cli"
	"strconv"
	"time"
)

const PKGName = "imager"

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	// 动态bpf日志区域
	bpfLogs := widget.NewTextGrid()
	bpfLogs.SetText("Real-Time BpfLogs")
	bpfLogsScroll := container.NewScroll(bpfLogs)
	bpfLogsScroll.SetMinSize(fyne.NewSize(400, 200)) // 限制宽度为 400，高度为 200
	// image 日志
	ImageLogs := widget.NewTextGrid()
	ImageLogs.SetText("Container-Creating Logs")
	ImageLogsScroll := container.NewScroll(ImageLogs)
	ImageLogsScroll.SetMinSize(fyne.NewSize(400, 200)) // 限制宽度为 400，高度为 200

	jsonEditor := widget.NewMultiLineEntry()
	jsonEditor.SetMinRowsVisible(10)
	jsonEditor.SetPlaceHolder(`Enter Docker JSON Config here...`)

	runButton := widget.NewButton("Run Docker", func() {
		jsonConfig := jsonEditor.Text
		result, err := cli.ParseAndRunDockerRun(jsonConfig)
		if err != nil {
			// 显示错误信息
			ImageLogs.SetText(fmt.Sprintf("Error: %v\n%s", err, result))
		} else {
			// 显示成功信息
			ImageLogs.SetText(fmt.Sprintf("Success:\n%s", result))
		}
	})

	// 筛选工具栏
	imagePuller := widget.NewEntry()
	// entry与button的事件触发函数
	pullImageFunc := func(s string) {
		imageName := s
		ImageLogs.SetText(fmt.Sprintf("pulling %s image\n", s))
		pullInfo, err := cli.PullDockerImage(imageName)
		if err != nil {
			return
		}
		ImageLogs.SetText(pullInfo)
	}
	imagePullerButton := widget.NewButton("Pull                  		", func() {
		pullImageFunc(imagePuller.Text)
	})
	imagePuller.SetPlaceHolder("Enter image name to pull")
	imagePuller.OnSubmitted = pullImageFunc

	ImageSelector := widget.NewSelect([]string{}, func(s string) {
		fmt.Printf("select existed image %s\n", s)
	})

	//var bpfChoices
	toolbar := container.NewVBox(
		imagePuller,
		imagePullerButton,
		ImageSelector,
		widget.NewSeparator(),
		widget.NewCheck("network", func(b bool) {}),
		widget.NewCheck("volume/io", func(b bool) {}),
		widget.NewCheck("cgroup", func(b bool) {}),
		widget.NewCheck("namespace", func(b bool) {}),
		widget.NewCheck("process", func(b bool) {}),
		widget.NewCheck("cpu", func(b bool) {}),
		jsonEditor,
		runButton,
	)
	toolbar.Resize(fyne.NewSize(200, 200))
	// 模拟实时更新日志和图表数据
	go simulateImageData(bpfLogs)
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
		logs.SetText(logs.Text() + "\n" + logEntry)
		// 模拟滚动：将光标移动到文本末尾
		//logs.CursorRow = len(logs.Text)
		logs.Refresh()
		time.Sleep(30 * time.Second)
	}
}
