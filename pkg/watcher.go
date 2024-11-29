package pkg

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	themes2 "shoggothforever/beefine/pkg/gui/themes"
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

var WatcherIndex = map[string][]string{
	"": {"welcome", "bpf", "docker"},
}

func CreateWatcher() fyne.CanvasObject {
	a := fyne.CurrentApp()
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	w := a.NewWindow("beefine")
	setWatcher := func(t Watcher) {
		title.SetText(t.Title)
		intro.SetText(t.Intro)
		title.Hide()
		intro.Hide()
		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}
	setWatcher(Watchers["welcome"])
	watcherBoarder := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)

	var preferenceCurrentWatcher = "currentWatcher"
	t := &widget.Tree{
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
	treeBoarder := container.NewBorder(nil, themes2.CreateThemes(a), nil, nil, t)

	split := container.NewHSplit(treeBoarder, watcherBoarder)
	split.SetOffset(0.2)
	return split
}

func WatcherStart() {
	a := app.NewWithID("io.watcher.beefine")
	w := a.NewWindow("beefine")
	// 初始化窗口内容
	w.SetContent(CreateWatcher())
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
