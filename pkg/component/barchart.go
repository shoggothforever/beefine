package component

import (
	"fmt"
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
		rect.SetMinSize(fyne.NewSize(40, height))
		rects[i] = rect
	}

	// 使用 HBox 布局
	container := container.NewHBox(rects...)
	return widget.NewSimpleRenderer(container)
}

// LiveChart 定义一个实时折线图组件
type LiveChart struct {
	widget.BaseWidget
	bar   *canvas.Rectangle
	value float32
	mu    sync.Mutex
}

// NewLiveChart 创建一个新的实时折线图组件
func NewLiveChart(maxPoints int, chartColor color.Color) *LiveChart {
	bar := canvas.NewRectangle(color.White)
	lc := &LiveChart{
		bar: bar,
	}
	lc.bar.Resize(fyne.NewSize(0, 50))
	lc.ExtendBaseWidget(lc)
	return lc
}

// AppendData 追加数据点并刷新图表
func (lc *LiveChart) AppendData(value float64) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.value = float32(value)
	lc.bar.Resize(fyne.NewSize(min(lc.value, 500), 50))

	lc.Refresh()
}

type LiveChartRenderer struct {
	l *LiveChart
}

func (l LiveChartRenderer) Destroy() {
}

func (l LiveChartRenderer) Layout(size fyne.Size) {
}

func (l LiveChartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(20, 0)
}

func (l LiveChartRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{l.l.bar, canvas.NewText(fmt.Sprintf("%d", int(l.l.value)), color.Black)}
}

func (l LiveChartRenderer) Refresh() {
	canvas.Refresh(l.l)
}

// CreateRenderer 实现 Fyne 的渲染逻辑
func (lc *LiveChart) CreateRenderer() fyne.WidgetRenderer {
	return &LiveChartRenderer{l: lc}
}
