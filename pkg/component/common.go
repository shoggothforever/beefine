package component

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"sync"
)

type TabManager struct {
	TabItemsMap map[string]*container.TabItem
	Tabs        *container.AppTabs
	w           *fyne.Window
	m           sync.Mutex
}

var tabManagerMap = make(map[string]*TabManager)

func NewTabManager(pkg string, w *fyne.Window) *TabManager {
	if _, ok := tabManagerMap[pkg]; !ok {
		tabManagerMap[pkg] = &TabManager{
			TabItemsMap: make(map[string]*container.TabItem),
			Tabs:        container.NewAppTabs(),
			w:           w,
			m:           sync.Mutex{},
		}
	}
	return tabManagerMap[pkg]
}
func GetTabManager(pkg string) *TabManager {
	return tabManagerMap[pkg]
}

func (t *TabManager) Remove(item *container.TabItem) {
	t.m.Lock()
	defer t.m.Unlock()
	if t.Tabs == nil || item == nil {
		return
	}
	t.Tabs.Remove(item)
	fmt.Println("remove ", item.Text)
	if t.TabItemsMap == nil {
		return
	}
	delete(t.TabItemsMap, item.Text)
}
func (t *TabManager) Append(item *container.TabItem, newTab bool) {
	t.m.Lock()
	defer t.m.Unlock()
	if t.Tabs == nil || item == nil {
		return
	}
	if t.TabItemsMap == nil {
		return
	}
	if _, ok := t.TabItemsMap[item.Text]; ok {
		t.Tabs.Select(item)
		return
	}
	t.TabItemsMap[item.Text] = item
	t.Tabs.Append(item)
	if newTab {
		t.Tabs.Select(item)
	}

}
func (t *TabManager) Select(item *container.TabItem) {
	if t.Tabs == nil {
		return
	}
	t.Tabs.Select(item)
}
func (t *TabManager) Get(title string) *container.TabItem {
	t.m.Lock()
	defer t.m.Unlock()
	if t.TabItemsMap == nil {
		return nil
	}
	var item *container.TabItem
	var ok bool
	if item, ok = t.TabItemsMap[title]; ok {
		t.Tabs.Select(item)
		return item
	}
	return item
}
func NewUIVBox(pkg string, title string, cancel func(), objs ...fyne.CanvasObject) *fyne.Container {
	objs = append([]fyne.CanvasObject{NewCloseButton(pkg, title, cancel)}, objs...)
	return container.NewVBox(objs...)
}
