package docker

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log"
	"shoggothforever/beefine/internal/cli"
	"shoggothforever/beefine/pkg/component"
	"time"
)

const PKGName = "docker"

var tabUIButtonFuncMap = map[string]func() fyne.CanvasObject{}

// Screen
func Screen(w fyne.Window) fyne.CanvasObject {

	// 创建选项卡容器
	tabM := component.NewTabManager(PKGName, &w)
	if homeItem := tabM.Get("home"); homeItem != nil {
		// 布局
		tabM.Select(homeItem)
		return tabM.Tabs
	}
	objects := []fyne.CanvasObject{widget.NewLabel("Select a feature:")}
	for name, fn := range tabUIButtonFuncMap {
		objects = append(objects, component.NewUITabButton(PKGName, name, fn))
	}
	// 其他功能按钮
	placeholderButton := widget.NewButton("Placeholder Feature", func() {
		widget.ShowPopUp(widget.NewLabel("Feature coming soon!"), w.Canvas())
	})
	objects = append(objects, placeholderButton)
	home := container.NewGridWithColumns(4,
		objects...,
	)
	homeItem := container.NewTabItem("home", home)
	tabM.Append(homeItem, true)
	// 布局
	tabM.Select(homeItem)
	containerCountLabel := widget.NewLabel("Containers: 0")
	imageCountLabel := widget.NewLabel("Images: 0")
	cpuUsageBar := widget.NewProgressBar()
	memoryUsageBar := widget.NewProgressBar()
	go func(imageCountLabel *widget.Label, cpuUsageBar *widget.ProgressBar) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			// 设置周期性数据更新
			select {
			case <-ticker.C:
				usageData, err := cli.GetDockerDashBoardData()
				if err != nil {
					log.Println("Error fetching Docker data:", err)
					continue
				}
				// 更新界面上的容器和镜像数量
				containerCountLabel.SetText(fmt.Sprintf("Containers: %d", usageData.ContainerLen))
				imageCountLabel.SetText(fmt.Sprintf("Images: %d", usageData.ImagesLen))

				// 更新 CPU 和内存使用情况
				cpuUsageBar.SetValue(usageData.CpuUsage / float64(100))
				memoryUsageBar.SetValue(usageData.MemUsage / float64(4))
			}
		}
	}(imageCountLabel, cpuUsageBar)
	// 布局容器
	content := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Docker Dashboard"),
		containerCountLabel,
		imageCountLabel,
		widget.NewLabel("CPU Usage"),
		cpuUsageBar,
		widget.NewLabel("Memory Usage"),
		memoryUsageBar,
	)

	return content
}
