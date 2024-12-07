package container

import (
	"bufio"
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"os"
	"shoggothforever/beefine/internal/cli"
	"strings"
	"sync"
	"unsafe"
)

// ContainerSelect 是自定义控件，包装了 Select 并添加了额外的字段
type ContainersSelect struct {
	widget.BaseWidget                // 嵌入 BaseWidget
	base              *widget.Select // 内嵌 Select
	containerLogs     *widget.TextGrid
	bpfLogs           *widget.TextGrid
	containerButton   *widget.Button
	cancelMap         map[string]func()
	containers        map[string]*types.Container
	currentContainer  *types.Container
	watcherPID        int //观测的docker 容器pid
	m                 sync.Mutex
}

// NewContainersSelect 创建自定义控件实例
func NewContainersSelect(containerLogs *widget.TextGrid, bpfLogs *widget.TextGrid) *ContainersSelect {
	// 初始化 Select
	containers, err := cli.ListContainer()
	if err != nil {
		return nil
	}
	var containertags []string
	containerMap := make(map[string]*types.Container)
	for k, v := range containers {
		for _, name := range v.Names {
			name = buildTag(name, &v)
			containertags = append(containertags, name)
			containerMap[name] = &containers[k]
			break
		}
	}
	selectWidget := widget.NewSelect(containertags, nil)
	// 创建 MyCustomWidget 实例
	s := &ContainersSelect{
		base:          selectWidget,
		cancelMap:     make(map[string]func()),
		m:             sync.Mutex{},
		containers:    containerMap,
		containerLogs: containerLogs,
		bpfLogs:       bpfLogs,
	}
	s.base.OnChanged = s.OnChanged
	s.base.PlaceHolder = "select existed containers"
	s.containerButton = widget.NewButton("waiting...", s.OnClick)
	s.ExtendBaseWidget(s) // 必须扩展 BaseWidget
	return s
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (w *ContainersSelect) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(w.base)
}

func (w *ContainersSelect) OnChanged(s string) {
	if _, ok := w.containers[s]; ok {
		w.currentContainer = w.containers[s]
		fmt.Println("update currentContainer ", w.currentContainer.ID)
		if cli.CheckContainerRunningState(w.currentContainer.Status) {
			w.containerButton.SetText("stop container")
			w.AppendLogInLock(w.containerLogs, fmt.Sprintf("container %s is stopped", w.currentContainer.ID[:16]))
		} else {
			w.containerButton.SetText("start container")
			w.AppendLogInLock(w.containerLogs, fmt.Sprintf("container %s is running", w.currentContainer.ID[:16]))
		}
	}

}
func (w *ContainersSelect) OnClick() {
	if w.currentContainer == nil {
		return
	}
	err := cli.ChangeContainerState(w.currentContainer.ID, cli.CheckContainerRunningState(w.currentContainer.Status))
	if err != nil {
		fmt.Println("change container state error", err)
	}
	status, err := cli.GetContainerStat(w.currentContainer.ID)
	w.currentContainer.Status = status.State.Status
	w.containers[buildTag(w.currentContainer.Names[0], w.currentContainer)] = w.currentContainer
	if cli.CheckContainerRunningState(w.currentContainer.Status) {
		w.containerButton.SetText("stop container")
		w.AppendLogInLock(w.containerLogs, fmt.Sprintf("start container %s ", w.currentContainer.ID[:16]))
	} else {
		w.containerButton.SetText("start container")
		w.AppendLogInLock(w.containerLogs, fmt.Sprintf("stop container %s", w.currentContainer.ID[:16]))
	}
}
func NewContainerToolBar(containerLogs *widget.TextGrid, bpfLogs *widget.TextGrid) *fyne.Container {
	containerSelector := NewContainersSelect(containerLogs, bpfLogs)
	return container.NewVBox(
		containerSelector,
		widget.NewSeparator(),
		containerSelector.containerButton,
		widget.NewCheck("volume/io", containerSelector.chooseVolume),
		widget.NewCheck("cgroup", containerSelector.chooseCgroup),
		widget.NewCheck("namespace", containerSelector.chooseNamespace),
		widget.NewCheck("process", containerSelector.chooseProcess),
		widget.NewCheck("cpu", containerSelector.chooseCpu),
		widget.NewCheck("memory", containerSelector.chooseCpu),
	)
}

func (w *ContainersSelect) chooseVolume(b bool) {
	if b {

	} else {

	}
}

func (w *ContainersSelect) chooseCgroup(b bool) {
	if b {

	} else {

	}
}

func (w *ContainersSelect) chooseNamespace(b bool) {
	if b {

	} else {

	}
}

func (w *ContainersSelect) chooseProcess(b bool) {
	if b {

	} else {

	}
}

func (w *ContainersSelect) chooseCpu(b bool) {
	if b {

	} else {

	}
}
func (w *ContainersSelect) chooseMemory(b bool) {
	if b {

	} else {

	}
}

func (w *ContainersSelect) AppendLogInLock(logs *widget.TextGrid, text string) {
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

func Bytes2String(b []byte) string {
	trimmedData := bytes.TrimRight(b, "\x00")
	return *(*string)(unsafe.Pointer(&trimmedData))
}

func buildTag(name string, c *types.Container) string {
	return name + " " + c.ID + " " + c.Image
}
func parseTagName(tag string) string {
	ss := strings.Split(tag, " ")
	if len(ss) > 0 {
		return ss[0]
	}
	return ""
}
func parseTagID(tag string) string {
	ss := strings.Split(tag, " ")
	if len(ss) > 1 {
		return ss[1]
	}
	return ""
}
