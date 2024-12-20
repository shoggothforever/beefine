package pkg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/internal/data/assets"
	themes2 "shoggothforever/beefine/pkg/gui/themes"
	"shoggothforever/beefine/pkg/gui/watcher/bpf"
	"shoggothforever/beefine/pkg/gui/watcher/docker"
	container2 "shoggothforever/beefine/pkg/gui/watcher/docker/container"
	"shoggothforever/beefine/pkg/gui/watcher/docker/imager"
	"shoggothforever/beefine/pkg/gui/watcher/welcome"
)

// Watcher defines the data structure for a tutorial
type Watcher struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

// Watchers 应用程序UI目录树中的结点信息
var Watchers = map[string]Watcher{
	"welcome":   {"welcome", "Welcome to the beefine observer", welcome.Screen},
	"BPF":       {"Load eBPF", "Observe system-level activities", bpf.Screen},
	"Docker":    {"Docker", "Monitor Docker activities", docker.Screen},
	"imager":    {"Image Monitoring", "Monitor Docker imager creation process", imager.Screen},
	"container": {"Container Monitoring", "Monitor running container performance", container2.Screen},
}

// WatcherIndex 目录树UI中各个节点的连接关系
var WatcherIndex = map[string][]string{
	"":       {"welcome", "BPF", "Docker"},
	"Docker": {"imager", "container"},
}

// viewsSet 保存创建过的canvas信息，避免重复创建
var viewsSet = map[string][]fyne.CanvasObject{}

func CreateTree(setWatcher func(t Watcher)) *widget.Tree {
	a := fyne.CurrentApp()
	var preferenceCurrentWatcher = "currentWatcher"
	return &widget.Tree{
		BaseWidget:     widget.BaseWidget{},
		Root:           "",
		HideSeparators: false,
		ChildUIDs: func(uid widget.TreeNodeID) (c []widget.TreeNodeID) {
			return WatcherIndex[uid]
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Placeholder Feature")
		},
		IsBranch: func(uid widget.TreeNodeID) (ok bool) {
			node, ok := WatcherIndex[uid]
			return ok && len(node) > 0
		},
		OnBranchClosed: nil,
		OnBranchOpened: nil,
		OnSelected: func(uid widget.TreeNodeID) {
			if t, ok := Watchers[uid]; ok {
				a.Preferences().SetString(preferenceCurrentWatcher, uid)
				setWatcher(t)
			}
		},
		OnUnselected: nil,
		UpdateNode: func(uid widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
			watcher, ok := Watchers[uid]
			if !ok {
				return
			}
			node.(*widget.Label).SetText(watcher.Title)
		},
	}
}

// CreateWatcher 创建watcher看板的UI结构
func CreateWatcher() fyne.CanvasObject {
	a := fyne.CurrentApp()
	a.Settings().SetTheme(&themes2.ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	// 设置fyne的tree UI 选中时会执行的操作
	setWatcher := func(t Watcher) {
		title.SetText(t.Title)
		intro.SetText(t.Intro)
		title.Hide()
		intro.Hide()
		if t.View != nil {
			if _, ok := viewsSet[t.Title]; !ok {
				viewsSet[t.Title] = []fyne.CanvasObject{t.View(TopWindow)}
			}
			content.Objects = viewsSet[t.Title]
		}
		content.Refresh()
	}

	//设置初始的选项卡
	setWatcher(Watchers["welcome"])
	watcherBoarder := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	//创建目录树结构
	t := CreateTree(setWatcher)
	treeBoarder := container.NewBorder(nil, themes2.CreateThemes(a), nil, nil, t)

	split := container.NewHSplit(treeBoarder, watcherBoarder)
	split.SetOffset(0.2)
	return split
}

var TopWindow fyne.Window

func WatcherStart() {
	//配置应用ID
	a := app.NewWithID("io.watcher.beefine")
	//设置应用ICON
	a.SetIcon(fyne.NewStaticResource("icon", assets.LogoIcon))
	//设置程序主窗口
	TopWindow = a.NewWindow("beefine")
	// 初始化窗口内容
	TopWindow.SetContent(CreateWatcher())
	TopWindow.Resize(fyne.NewSize(1600, 800))
	TopWindow.ShowAndRun()
}
