package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"shoggothforever/beefine/bpf/counter"
)

const CounterUIName = "countNetPackage"

func CounterUI() fyne.CanvasObject {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	entry := widget.NewEntry()
	entry.OnSubmitted = func(s string) {
		req := counter.CounterReq{IfName: s}
		out, cancel := counter.Start(req)
		defer cancel()
		for v := range out {
			fmt.Println(v)
		}
	}

	return NewUIVBox(
		CounterUIName,
		widget.NewLabel("click to check how many net package have been received"),
		entry,
	)
}
