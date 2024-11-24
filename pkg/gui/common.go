package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"sync"
)

var once sync.Once

type TabManager struct {
	TabItemsMap map[string]*container.TabItem
	Tabs        *container.AppTabs
	w           *fyne.Window
	m           sync.Mutex
}

var tabManager *TabManager

func NewTabManager(w *fyne.Window) *TabManager {
	once.Do(func() {
		tabManager = &TabManager{
			TabItemsMap: make(map[string]*container.TabItem),
			Tabs:        container.NewAppTabs(),
			w:           w,
			m:           sync.Mutex{},
		}
	})
	return tabManager
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
