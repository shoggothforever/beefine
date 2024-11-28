package pkg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/pkg/component"
	"shoggothforever/beefine/pkg/gui/watcher/bpf"
	"shoggothforever/beefine/pkg/gui/watcher/welcome"
)

// Watcher defines the data structure for a tutorial
type Watcher struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var Watchers = map[string]Watcher{
	"welcome": {"welcome", "", welcome.Screen},
	"bpf":     {"bpf", "", bpf.Screen},
}

func CreateWatcherTree() *widget.Tree {
	//a := fyne.CurrentApp()
	t := &widget.Tree{
		BaseWidget:     widget.BaseWidget{},
		Root:           "",
		HideSeparators: false,
		ChildUIDs:      nil,
		CreateNode:     nil,
		IsBranch:       nil,
		OnBranchClosed: nil,
		OnBranchOpened: nil,
		OnSelected:     nil,
		OnUnselected:   nil,
		UpdateNode:     nil,
	}

	return t
}

var WatcherIndex = map[string][]string{
	"": {"welcome", "bpf", "docker"},
}

func WatcherStart() {
	a := app.NewWithID("io.watch.ebpf")
	w := a.NewWindow("eBPF Hook Manager")
	// 创建选项卡容器
	tabM := component.NewTabManager("bpf", &w)

	// 创建 Home 页面
	homeTab := container.NewTabItem("Home", bpf.Screen(w))
	tabM.Append(homeTab, true)

	// 初始化窗口内容
	w.SetContent(tabM.Tabs)

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
