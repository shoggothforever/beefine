package imager

import (
	"bufio"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"reflect"
	"shoggothforever/beefine/bpf/image_prep"
	"shoggothforever/beefine/internal/cli"
	"slices"
	"strings"
	"sync"
	"unsafe"
)

// MyCustomWidget 是自定义控件，包装了 Select 并添加了额外的字段
type ImageSelect struct {
	widget.BaseWidget                // 嵌入 BaseWidget
	base              *widget.Select // 内嵌 Select
	currentImage      string         // 额外的字段
	watcherID         int            //观测的docker 容器id
	imageLogs         *widget.TextGrid
	bpfLogs           *widget.TextGrid
	cancelMap         map[string]func()
	m                 sync.Mutex
}

// NewImageSelect 创建自定义控件实例
func NewImageSelect(options []string, onSelected func(string)) *ImageSelect {
	// 初始化 Select
	selectWidget := widget.NewSelect(options, onSelected)

	// 创建 MyCustomWidget 实例
	s := &ImageSelect{
		base:         selectWidget,
		currentImage: "default data", // 设置默认值
		cancelMap:    make(map[string]func()),
		m:            sync.Mutex{},
	}

	s.ExtendBaseWidget(s) // 必须扩展 BaseWidget
	return s
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (w *ImageSelect) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(w.base)
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
	imagePullerButton := widget.NewButton("Pull                  		", func() {
		pullImageFunc(imagePuller.Text)
	})
	imagePuller.SetPlaceHolder("Enter image name to pull")
	imagePuller.OnSubmitted = pullImageFunc
	images, err := cli.ListImage()
	if err != nil {
		return nil
	}
	var imagetags []string
	for _, v := range images {
		for _, tag := range v.RepoTags {
			imagetags = append(imagetags, tag)
			break
		}
	}
	imageSelector := NewImageSelect(imagetags, func(s string) {
		fmt.Printf("select existed image %s\n", s)
	})
	imageSelector.base.PlaceHolder = "select existed image"
	imageSelector.imageLogs = ImageLogs
	imageSelector.bpfLogs = bpfLogs
	jsonEditor := widget.NewMultiLineEntry()
	jsonEditor.SetMinRowsVisible(10)
	jsonEditor.SetPlaceHolder(`{
    "image": "nginx",
    "name": "my-container",
    "ports": [],
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
			ImageLogs.SetText(fmt.Sprintf("Error: %v\n%s", err, result))
		} else {
			// 显示成功信息
			ImageLogs.SetText(fmt.Sprintf("Success:\n%s", result))
		}
	}
	jsonEditor.OnSubmitted = func(s string) {
		onClick()
	}
	runButton := widget.NewButton("run", onClick)
	v := reflect.TypeOf(imageSelector)
	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		fmt.Println(m.Name)
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
		widget.NewCheck("volume/io", imageSelector.chooseVolume),
		widget.NewCheck("cgroup", imageSelector.chooseCgroup),
		widget.NewCheck("namespace", imageSelector.chooseNamespace),
		widget.NewCheck("process", imageSelector.chooseProcess),
		widget.NewCheck("cpu", imageSelector.chooseCpu),
		jsonEditor,
		runButton,
	)
}
func (w *ImageSelect) chooseUnionFS(b bool) {
	if b == true {
		req := image_prep.ImagePrepReq{}
		out, cancel := image_prep.Start(&req)
		w.cancelMap["chooseUnionFS"] = cancel
		go func() {
			strs := []string{"sudo", "beefine", "gmain", "gnome-terminal-", "gnome-shell", "systemd-oomd"}
			for v := range out {
				//time.Sleep(time.Second)
				comm := Bytes2String(v.Comm[:])
				if slices.Contains(strs, comm) {
					break
				}
				w.m.Lock()
				str := fmt.Sprintf("pid:%d,comm:%s,operation:%s", v.Pid, comm, Bytes2String(v.Operation[:]))
				fmt.Println(str)
				w.AppendLogInLock(w.bpfLogs, str)
				w.m.Unlock()
			}
		}()
	} else {
		w.cancelMap["chooseUnionFS"]()
	}
	fmt.Println("choose watch unionfs")
}
func (w *ImageSelect) chooseVolume(b bool) {

}
func (w *ImageSelect) chooseCgroup(b bool) {

}
func (w *ImageSelect) chooseNamespace(b bool) {

}
func (w *ImageSelect) chooseProcess(b bool) {

}
func (w *ImageSelect) chooseCpu(b bool) {

}
func (w *ImageSelect) AppendLogInLock(logs *widget.TextGrid, text string) {
	file, err := os.OpenFile("tmplog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	logs.SetRow(len(logs.Rows), widget.NewTextGridFromString(text).Row(0))
	if len(logs.Rows) > 200 {
		logs.SetText("")
	}
	//logs.SetText(text)
}
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	trimmedData := bytes.TrimRight(b, "\x00")
	return *(*string)(unsafe.Pointer(&trimmedData))
}
