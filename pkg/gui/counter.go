package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
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
	// 动态更新柱状图的函数
	// 柱状图容器
	chartContainer := container.NewVBox()
	updateChart := func(data int) {
		// 模拟柱状图高度变化
		bar := canvas.NewRectangle(color.Black)
		bar.SetMinSize(fyne.NewSize(10, float32(data)))
		chartContainer.Add(bar)

		// 保留最近 10 个柱状图
		if len(chartContainer.Objects) > 10 {
			chartContainer.Objects = chartContainer.Objects[1:]
		}
		canvas.Refresh(chartContainer)
	}

	var cancelFunc func()
	selectIface.OnChanged = func(s string) {
		req := counter.CounterReq{IfName: s}
		out, cancel := counter.Start(req)
		cancelFunc = cancel
		statusLabel.SetText(fmt.Sprintf("Status: Monitoring %s", s))
		stopButton.Enable()
		go func() {
			for v := range out {
				updateChart(int(v.Count))
			}
			statusLabel.SetText("Status: Idle")
			stopButton.Disable()
		}()
	}
	// 停止按钮事件
	stopButton.OnTapped = func() {
		if cancelFunc != nil {
			cancelFunc() // 调用关闭函数
			statusLabel.SetText("Status: Stopped")
			stopButton.Disable()
		}
		chartContainer.Objects = []fyne.CanvasObject{}
		chartContainer.Refresh()
	}
	// 模拟动态数据（调试用）

	return NewUIVBox(
		CounterUIName,
		widget.NewLabel("click to check how many net package have been received"),
		selectIface,
		statusLabel,
		stopButton,
		chartContainer,
	)
}
