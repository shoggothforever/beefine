package container

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/docker/docker/api/types"
	"strings"
	"unsafe"
)

const PKGName = "container"

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	// 动态bpf日志区域
	bpfLogs := widget.NewTextGrid()
	bpfLogs.ShowLineNumbers = true
	bpfLogs.SetText("Real-Time BpfLogs")
	//bpfLogs.
	bpfLogsScroll := container.NewScroll(bpfLogs)
	bpfLogsScroll.SetMinSize(fyne.NewSize(600, 200)) // 限制宽度为 400，高度为 200
	// image 日志
	containerLogs := widget.NewTextGrid()
	containerLogs.ShowLineNumbers = true
	containerLogs.SetText("Container-Creating Logs")
	ImageLogsScroll := container.NewScroll(containerLogs)
	ImageLogsScroll.SetMinSize(fyne.NewSize(600, 200)) // 限制宽度为 400，高度为 200

	//var bpfChoices
	toolbar := NewContainerToolBar(containerLogs, bpfLogs)
	content := container.NewHBox(
		toolbar,
		ImageLogsScroll,
		bpfLogsScroll,
	)
	return content
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
