package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/docker/docker/api/types"
	"shoggothforever/beefine/pkg/component"
	"strings"
)

const (
	PKGName         = "container"
	GetNsPeerScript = "scripts/nswatch.sh"
)

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	// 动态bpf日志区域
	bpfLogs := component.NewLogBoard("Real-Time BpfLogs", 600, 200)
	// image 日志
	containerLogs := component.NewLogBoard("Container-Creating Logs", 600, 200)
	toolbar := NewContainerToolBar(containerLogs, bpfLogs)
	content := container.NewHBox(
		toolbar,
		containerLogs,
		bpfLogs,
	)
	return content
}

func buildTag(name string, c *types.Container) string {
	return name + " " + c.ID[:min(len(c.ID), 16)] + " " + c.Image[:min(len(c.Image), 16)]
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
