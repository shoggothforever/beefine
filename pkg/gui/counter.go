package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"net"
	"shoggothforever/beefine/bpf/counter"
)

const CounterUIName = "countNetPackage"

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

	// 创建 Select 控件
	selectIface := widget.NewSelect(options, nil)
	selectIface.PlaceHolder = "Select a network interface"

	// 停止按钮（初始状态不可用）
	stopButton := widget.NewButton("Stop", nil)
	stopButton.Disable()
	// 运行逻辑
	// 计数标签
	cntLabel := widget.NewLabel("Counter")
	cntLabel.SetText("waiting to count")

	var cancelFunc func()
	selectIface.OnChanged = func(s string) {
		req := counter.CounterReq{IfName: s}
		out, cancel := counter.Start(req)
		cancelFunc = cancel
		statusLabel.SetText(fmt.Sprintf("Status: Monitoring %s", s))
		stopButton.Enable()
		go func() {
			for v := range out {
				cntLabel.SetText(fmt.Sprintf("Received %d packets", v.Count))
			}
			statusLabel.SetText("Status: Idle")
			stopButton.Disable()
		}()
	}
	stop := func() {
		if cancelFunc != nil {
			cancelFunc() // 调用关闭函数
			statusLabel.SetText("Status: Stopped")
			stopButton.Disable()
			cntLabel.SetText("waiting to count")
		}
	}
	defer stop()
	// 停止按钮事件
	stopButton.OnTapped = stop
	// 模拟动态数据（调试用）

	return NewUIVBox(
		CounterUIName,
		stop,
		widget.NewLabel("click to check how many net package have been received"),
		selectIface,
		statusLabel,
		stopButton,
		widget.NewLabel("Real-Time Packet Counter:"),
		cntLabel,
	)
}
