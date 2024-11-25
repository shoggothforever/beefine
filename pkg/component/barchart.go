package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"sync"
)

type BarChart struct {
	widget.BaseWidget
	data   []int
	maxVal int
	mu     sync.Mutex
}

func NewBarChart() *BarChart {
	b := &BarChart{
		data:   make([]int, 0),
		maxVal: 1, // 防止除以零
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *BarChart) AppendData(value int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = append(b.data, value)
	if len(b.data) > 10 { // 限制为最近10个柱
		b.data = b.data[1:]
	}
	if value > b.maxVal {
		b.maxVal = value
	}
	b.Refresh()
}
func (b *BarChart) RemoveData() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = make([]int, 0)
	b.Refresh()
}
func (b *BarChart) CreateRenderer() fyne.WidgetRenderer {
	b.mu.Lock()
	defer b.mu.Unlock()

	rects := make([]fyne.CanvasObject, len(b.data))
	for i, v := range b.data {
		// 动态高度比例
		height := float32(v) / float32(b.maxVal) * 200
		rect := canvas.NewRectangle(color.RGBA{R: uint8(50 + i*20), G: 100, B: 200, A: 255})
		rect.SetMinSize(fyne.NewSize(20, height))
		rects[i] = rect
	}

	// 使用 HBox 布局
	container := container.NewHBox(rects...)
	return widget.NewSimpleRenderer(container)
}
