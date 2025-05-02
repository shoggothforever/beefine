package pod

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
)

const (
	PKGName = "pod"
)

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {
	// 动态bpf日志区域
	bpfLogs := component.NewLogBoard(PKGName+"BPF", "Real-Time BpfLogs", 600, 500)
	// image 日志
	containerLogs := component.NewLogBoard(PKGName+"runtime", "pod-communication Logs", 600, 500)
	toolbar := NewToolBar(containerLogs, bpfLogs)
	content := container.NewHBox(
		toolbar,
		containerLogs,
		bpfLogs,
	)
	return content
}

func NewToolBar(podLogs *component.LogBoard, bpfLogs *component.LogBoard) *fyne.Container {
	podSelector := NewPodSelect()
	podSelector.podLogs = podLogs
	podSelector.bpfLogs = bpfLogs
	return container.NewVBox(
		podSelector,
		widget.NewSeparator(),
		widget.NewCheck("network", podSelector.chooseNetwork),
	)
}
