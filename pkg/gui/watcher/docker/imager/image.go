package imager

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"reflect"
	"shoggothforever/beefine/internal/cli"
	"strings"
)

const PKGName = "imager"

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

const (
	netTracePointScript       = "scripts/net.bt"
	isolationTracePointScript = "scripts/isolation.bt"
)

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
	ImageLogs.ShowLineNumbers = true
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

func NewToolBar(ImageLogs *widget.TextGrid, bpfLogs *widget.TextGrid) *fyne.Container {
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
	imagePullerButton := widget.NewButton("Pull                  				", func() {
		pullImageFunc(imagePuller.Text)
	})
	imagePuller.SetPlaceHolder("Enter image name to pull")
	imagePuller.OnSubmitted = pullImageFunc
	imageSelector := NewImageSelect()
	imageSelector.imageLogs = ImageLogs
	imageSelector.bpfLogs = bpfLogs
	jsonEditor := widget.NewMultiLineEntry()
	jsonEditor.SetMinRowsVisible(8)
	jsonEditor.SetPlaceHolder(`{
    "image": "nginx",
    "name": "my-container",
    "ports": ["80:80"],
    "volumes": [],
    "env": [],
    "detach": true,
    "rm": true
	}`)
	onClick := func() {
		input := jsonEditor.Text
		result, err := cli.ParseAndRunDockerRun(input, imageSelector.base.Selected)
		//result, err := cli.ExecDockerCmd(input)
		imageSelector.m.Lock()
		if err != nil {
			// 显示错误信息
			imageSelector.AppendLogInLock(ImageLogs, fmt.Sprintf("Error: %v\n%s", err, result))
		} else {
			// 显示成功信息
			imageSelector.AppendLogInLock(ImageLogs, fmt.Sprintf("Success:\n%s", result))
		}
		imageSelector.m.Unlock()
	}
	jsonEditor.OnSubmitted = func(s string) {
		onClick()
	}
	runButton := widget.NewButton("run", onClick)
	v := reflect.TypeOf(imageSelector)
	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		//fmt.Println(m.Name)
		if strings.HasPrefix(m.Name, "choose") {
			imageSelector.cancelMap[m.Name] = nil
		}
	}
	return container.NewVBox(
		imagePuller,
		imagePullerButton,
		imageSelector,
		widget.NewSeparator(),
		widget.NewCheck("unionFS", imageSelector.chooseUnionFS),
		widget.NewCheck("mount", imageSelector.chooseMount),
		widget.NewCheck("network", imageSelector.chooseNetwork),
		widget.NewCheck("isolation", imageSelector.chooseIsolation),
		//widget.NewCheck("process", imageSelector.chooseProcess),
		jsonEditor,
		runButton,
	)
}
