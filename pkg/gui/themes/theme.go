package themes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type ForcedVariant struct {
	fyne.Theme

	Variant fyne.ThemeVariant
}

func (f *ForcedVariant) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return f.Theme.Color(name, f.Variant)
}
func CreateThemes(a fyne.App) *fyne.Container {
	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(&ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(&ForcedVariant{Theme: theme.DefaultTheme(), Variant: theme.VariantLight})
		}),
	)
	return themes
}
