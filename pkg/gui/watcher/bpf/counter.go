package bpf

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net"
	"shoggothforever/beefine/bpf/counter"
	"shoggothforever/beefine/pkg/component"
)

const CounterUIName = "InspectNetwork"

func CounterUI() fyne.CanvasObject {
	// 显示状态
	statusLabel := widget.NewLabel("Status: Idle")

	// 获取系统网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		statusLabel.SetText(fmt.Sprintf("Error: %v", err))
		return nil
	}

	// 填充 Select 选项
	options := []string{}
	for _, iface := range interfaces {
		options = append(options, iface.Name)
	}
	log := component.NewLogBoard("catch network info through xdp", 200, 400)
	// 创建 Select 控件
	selectIface := widget.NewSelect(options, nil)
	selectIface.PlaceHolder = "Select a network interface"

	// 停止按钮（初始状态不可用）
	stopButton := component.NewStopButton()
	// 运行逻辑
	// 计数标签
	cntLabel := widget.NewLabel("Counter")
	cntLabel.SetText("waiting to count")
	var cancelFunc func()
	var lastSelect string
	selectIface.OnChanged = func(s string) {
		if s == lastSelect {
			return
		}
		if cancelFunc != nil {
			cancelFunc()
		}
		lastSelect = s
		req := &counter.CounterReq{IfName: s}
		out, cancel := counter.Start(req)
		cancelFunc = cancel
		statusLabel.SetText(fmt.Sprintf("Status: Monitoring %s", s))
		stopButton.Enable()
		go func() {
			for v := range out {
				cntLabel.SetText(fmt.Sprintf("Received %d packets", v.Count))
				log.AppendLogf("Src IP: %s, Dst IP: %s, Src Port: %d, Dst Port: %d, Protocol: %s\n",
					v.SrcIp, v.DstIp, v.SrcPort, v.DstPort, v.Protocol)
			}
			stopButton.Disable()
		}()
	}
	stop := func() {
		if cancelFunc != nil {
			cancelFunc() // 调用关闭函数
			statusLabel.SetText("Status: Idle")
			stopButton.Disable()
			cntLabel.SetText("waiting to count")
			log.Clear()
		}
	}
	// 停止按钮事件
	stopButton.OnTapped = stop
	// 模拟动态数据（调试用）

	return component.NewUIVBox(
		PKGName,
		CounterUIName,
		stop,
		widget.NewLabel("using xdp to catch network package information"),
		widget.NewLabel("select network iface to attach xdp program"),
		selectIface,
		statusLabel,
		container.NewHBox(stopButton),
		widget.NewLabel("Real-Time Packet Counter:"),
		cntLabel,
		log,
	)
}
