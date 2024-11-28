package bpf

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/bpf/exec"
	"shoggothforever/beefine/pkg/component"
)

const ExecUIName = "trace syscall/exec"

func ExecUI() fyne.CanvasObject {
	// 显示状态
	statusLabel := widget.NewLabel("Status: Idle")
	var cancelFunc func()

	//LogLabel := widget.NewLabel("tracing exec event")
	req := &exec.ExecReq{}
	out, cancel := exec.Start(req)
	cancelFunc = cancel
	stopButton := component.NewStopButton()
	stopButton.Enable()
	stopButton.OnTapped = func() {
		if cancelFunc != nil {
			cancelFunc()
			stopButton.Disable()
		}
	}
	go func() {
		for v := range out {
			fmt.Println(v.Pid)
		}
	}()
	return component.NewUIVBox(
		PKGName,
		ExecUIName,
		stopButton.OnTapped,
		stopButton,
		statusLabel,
	)
}
