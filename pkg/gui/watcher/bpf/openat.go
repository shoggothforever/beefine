package bpf

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/bpf/image_prep"
	"shoggothforever/beefine/internal/helper"
	"shoggothforever/beefine/pkg/component"
)

const OpenAtUIName = "openat"

func OpenAtUI() fyne.CanvasObject {
	// 显示状态
	statusLabel := widget.NewLabel("Status: Idle")
	var cancelFunc func()
	req := &image_prep.ImagePrepReq{}
	out, cancel := image_prep.Start(req)
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
	log := component.NewLogBoard(OpenAtUIName, " ", 200, 400)
	go func() {
		mp := make(map[string]int)
		for event := range out {
			comm := helper.Bytes2String(event.Comm[:])
			filename := helper.Bytes2String(event.Filename[:])
			str := fmt.Sprintf("ppid:%d comm:%s,operation:%s,filename:%s", event.Ppid, comm, helper.Bytes2String(event.Operation[:]), filename)
			if _, ok := mp[str]; ok {
				mp[str]++
				continue
			}
			mp[str] = 1
			log.AppendLogf(str)
		}
	}()
	return component.NewUIVBox(
		PKGName,
		OpenAtUIName,
		stopButton.OnTapped,
		stopButton,
		statusLabel,
		log,
	)
}
