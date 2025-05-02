package bpf

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/bpf/mount"
	"shoggothforever/beefine/internal/helper"
	"shoggothforever/beefine/pkg/component"
)

const MountUIName = "mount"

func MountUI() fyne.CanvasObject {
	// 显示状态
	statusLabel := widget.NewLabel("Status: Idle")
	var cancelFunc func()

	req := &mount.MountReq{}
	out, cancel := mount.Start(req)
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
	log := component.NewLogBoard(MountUIName, " ", 200, 400)
	go func() {
		for e := range out {
			typ := helper.Bytes2String(e.Type[:])
			dir := helper.Bytes2String(e.DirName[:])
			dev := helper.Bytes2String(e.DevName[:])
			log.AppendLogf("pid: %d , type:%s , dir:%s , device:%s\n", e.Pid, typ, dir, dev)
		}
	}()
	return component.NewUIVBox(
		PKGName,
		MountUIName,
		stopButton.OnTapped,
		stopButton,
		statusLabel,
		log,
	)
}
