package logic

import "fyne.io/fyne/v2/data/binding"

func Counter() {
	str := binding.NewString()
	str.Set("counter")
}
