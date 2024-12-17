package main

import (
	_ "embed"
	"shoggothforever/beefine/pkg"
)

func main() {
	// 初始化 Fyne 应用

	// 加载主 UI
	pkg.WatcherStart()
}
