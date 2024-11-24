package pkg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"shoggothforever/beefine/pkg/gui"
)

// MainUI creates the main application window with navigation.
func MainUI(w fyne.Window) {
	// 创建选项卡容器
	tabM := gui.NewTabManager(&w)

	// 创建 Home 页面
	homeTab := container.NewTabItem("Home", gui.HomeUI(w))
	tabM.Append(homeTab, true)

	// 初始化窗口内容
	w.SetContent(tabM.Tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
