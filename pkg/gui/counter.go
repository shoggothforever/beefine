package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net"
	"shoggothforever/beefine/bpf/counter"
	"shoggothforever/beefine/pkg/component"
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
	stopButton := NewStopButton()
	// 运行逻辑
	// 计数标签
	cntLabel := widget.NewLabel("Counter")
	cntLabel.SetText("waiting to count")
	// 创建实时折线图
	liveChart := component.NewLiveChart(50, color.RGBA{R: 255, G: 100, B: 100, A: 255})

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
				liveChart.AppendData(float64(v.Count))
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
		}
	}
	// 停止按钮事件
	stopButton.OnTapped = stop
	// 模拟动态数据（调试用）

	return NewUIVBox(
		CounterUIName,
		stop,
		widget.NewLabel("click to check how many net package have been received"),
		selectIface,
		statusLabel,
		container.NewHBox(stopButton),
		widget.NewLabel("Real-Time Packet Counter:"),
		cntLabel,
		liveChart,
	)
}
