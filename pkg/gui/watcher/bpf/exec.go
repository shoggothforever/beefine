package bpf

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/bpf/exec"
	"shoggothforever/beefine/internal/helper"
	"shoggothforever/beefine/pkg/component"
)

const ExecUIName = "Trace SyscallExec"

func ExecUI() fyne.CanvasObject {
	// 显示状态
	statusLabel := widget.NewLabel("Status: Idle")
	var cancelFunc func()

	req := &exec.ExecReq{}
	out, cancel := exec.Start(req)
	statusLabel.SetText("monitor process exec syscall")
	cancelFunc = cancel
	stopButton := component.NewStopButton()
	stopButton.Enable()
	stopButton.OnTapped = func() {
		if cancelFunc != nil {
			cancelFunc()
			stopButton.Disable()
		}
	}
	log := component.NewLogBoard(" ", 200, 400)
	go func() {
		mp := make(map[string]uint64)
		for e := range out {
			comm := helper.Bytes2String(e.Comm[:])
			if e.ExitEvent {
				fmt.Printf("exit duration_ns:%v,prio:%d, pid: %d, comm: %s\n", e.Ts-mp[comm], e.Prio, e.Pid, comm)
				log.AppendLogf("exit duration_ns:%v,prio:%d, pid: %d, comm: %s\n", e.Ts-mp[comm], e.Prio, e.Pid, comm)
				delete(mp, comm)
			} else {
				mp[comm] = e.Ts
				fmt.Printf("exec pid: %d, comm: %s\n", e.Pid, comm)
				log.AppendLogf("exec pid: %d, comm: %s\n", e.Pid, comm)
			}
		}
	}()
	return component.NewUIVBox(
		PKGName,
		ExecUIName,
		stopButton.OnTapped,
		stopButton,
		statusLabel,
		log,
	)
}
