package beefine

import (
	"fyne.io/fyne/v2/app"
	"shoggothforever/beefine/pkg"
)

func main() {
	// 初始化 Fyne 应用
	a := app.NewWithID("io.watch.ebpf")
	w := a.NewWindow("eBPF Hook Manager")
	// 加载主 UI
	pkg.MainUI(w)
}
