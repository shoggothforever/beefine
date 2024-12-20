package imager

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"reflect"
	"shoggothforever/beefine/internal/cli"
	"shoggothforever/beefine/pkg/component"
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
	bpfLogs := component.NewLogBoard("Real-Time BpfLogs", 600, 200)
	// image 日志
	ImageLogs := component.NewLogBoard("Container-Creating Logs", 600, 200)

	//var bpfChoices
	toolbar := NewToolBar(ImageLogs, bpfLogs)
	// 模拟实时更新日志和图表数据
	//go simulateImageData(bpfLogs)
	content := container.NewHBox(
		toolbar,
		ImageLogs,
		bpfLogs,
	)
	return content
}

func NewToolBar(ImageLogs *component.LogBoard, bpfLogs *component.LogBoard) *fyne.Container {
	// 筛选工具栏
	imagePuller := widget.NewEntry()
	// entry与button的事件触发函数
	pullImageFunc := func(s string) {
		imageName := s
		ImageLogs.AppendLogf("pulling %s image\n", s)
		pullInfo, err := cli.PullDockerImage(imageName)
		if err != nil {
			return
		}
		ImageLogs.AppendLogf(pullInfo)
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
	“cmd" : "/bin/sh"
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
		if err != nil {
			// 显示错误信息
			ImageLogs.AppendLogf("Error: %v %s \n", err, result)
		} else {
			// 显示成功信息
			ImageLogs.AppendLogf("Success:%s \n", result)
		}
	}
	jsonEditor.OnSubmitted = func(s string) {
		onClick()
	}
	runButton := widget.NewButton("run", onClick)
	v := reflect.TypeOf(imageSelector)
	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		if strings.HasPrefix(m.Name, "choose") {
			imageSelector.cancelMap[m.Name] = nil
		}
	}
	return container.NewVBox(
		imagePuller,
		imagePullerButton,
		imageSelector,
		widget.NewSeparator(),
		widget.NewCheck("VFS", imageSelector.chooseVFS),
		widget.NewCheck("mount", imageSelector.chooseMount),
		widget.NewCheck("network", imageSelector.chooseNetwork),
		widget.NewCheck("isolation", imageSelector.chooseIsolation),
		jsonEditor,
		runButton,
	)
}
