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
	labels []string
	maxVal int
	mu     sync.Mutex
}

func NewBarChart() *BarChart {
	b := &BarChart{
		data:   []int{0, 0, 0, 0}, // 默认4层镜像层
		labels: []string{"Layer 1", "Layer 2", "Layer 3", "Layer 4"},
		maxVal: 100, // 最大值
	}
	b.ExtendBaseWidget(b)
	return b
}

func (b *BarChart) AppendData(layer int, value int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if layer >= 0 && layer < len(b.data) {
		b.data[layer] += value
		if b.data[layer] > b.maxVal {
			b.maxVal = b.data[layer]
		}
		b.Refresh()
	}
}

//func (bc *BarChart) CreateRenderer() fyne.WidgetRenderer {
//	bc.mu.Lock()
//	defer bc.mu.Unlock()
//
//	rects := make([]fyne.CanvasObject, len(b.data))
//	for i, v := range b.data {
//		// 动态高度比例
//		height := float32(v) / float32(b.maxVal) * 200
//		rect := canvas.NewRectangle(color.RGBA{R: uint8(50 + i*20), G: 100, B: 200, A: 255})
//		rect.SetMinSize(fyne.NewSize(40, height))
//		rects[i] = rect
//	}
//
//	// 使用 HBox 布局
//	container := container.NewHBox(rects...)
//	return widget.NewSimpleRenderer(container)
//}

// CreateRenderer 实现 Fyne 的渲染逻辑
func (b *BarChart) CreateRenderer() fyne.WidgetRenderer {
	return &BarChartRenderer{b}
}

type BarChartRenderer struct {
	*BarChart
}

func (b BarChartRenderer) Destroy() {

}

func (b BarChartRenderer) Layout(size fyne.Size) {

}

func (b BarChartRenderer) MinSize() fyne.Size {
	return fyne.NewSize(20, 0)
}

func (b BarChartRenderer) Objects() []fyne.CanvasObject {
	b.mu.Lock()
	defer b.mu.Unlock()

	bars := make([]fyne.CanvasObject, len(b.data))
	for i, val := range b.data {
		height := float32(val) / float32(b.maxVal) * 200
		rect := canvas.NewRectangle(color.RGBA{R: uint8(50 + i*50), G: 150, B: 200, A: 255})
		rect.SetMinSize(fyne.NewSize(40, height))
		label := canvas.NewText(b.labels[i], color.Black)
		label.Alignment = fyne.TextAlignCenter

		bars[i] = container.NewVBox(rect, label)
	}

	barContainer := container.NewHBox(bars...)
	return []fyne.CanvasObject{barContainer}
}

func (b BarChartRenderer) Refresh() {
	canvas.Refresh(b)
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

// CreateRenderer 实现 Fyne 的渲染逻辑
func (lc *LiveChart) CreateRenderer() fyne.WidgetRenderer {
	return &LiveChartRenderer{l: lc}
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
