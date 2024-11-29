package welcome

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net/url"
	"strings"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func Screen(_ fyne.Window) fyne.CanvasObject {
	logo := canvas.NewImageFromFile("internal/data/assets/logo.png")
	logo.FillMode = canvas.ImageFillContain
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	footer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewHyperlink("beefine.io", parseURL("")),
		widget.NewLabel("-"),
		widget.NewHyperlink("documentation", parseURL("")),
		widget.NewLabel("-"),
		widget.NewHyperlink("github", parseURL("https://github.com/shoggothforever/beefine")),
		widget.NewLabel("-"),
		widget.NewHyperlink("gitlab", parseURL("https://gitlab.eduxiji.net/T202410336994295/project2608128-273615")),
		layout.NewSpacer(),
	)
	author := fyne.StaticResource{}
	author.StaticName = "author"
	author.StaticContent = []byte("蔡龙祥 <1337231450@qq.com>\n谭文轩 <2@qq.com>\n李睿涵 <3@qq.com>\n\n")
	authors := widget.NewRichTextFromMarkdown(formatAuthors(string(author.Content())))
	content := container.NewVBox(
		widget.NewLabelWithStyle("\n\nWelcome to the beefine app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		logo,
		container.NewCenter(authors),
	)
	scroll := container.NewScroll(content)

	bgColor := withAlpha(theme.Color(theme.ColorNameBackground), 0xe0)
	shadowColor := withAlpha(theme.Color(theme.ColorNameBackground), 0x33)

	bg := canvas.NewRectangle(bgColor)

	footerBG := canvas.NewRectangle(shadowColor)

	listen := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(listen)
	go func() {
		for range listen {
			bgColor = withAlpha(theme.Color(theme.ColorNameBackground), 0xe0)
			bg.FillColor = bgColor
			bg.Refresh()

			shadowColor = withAlpha(theme.Color(theme.ColorNameBackground), 0x33)
			footerBG.FillColor = bgColor
			footer.Refresh()
		}
	}()

	return container.NewStack(container.New(unpad{top: true}, bg),
		container.NewBorder(nil,
			container.NewStack(footerBG, footer), nil, nil,
			container.New(unpad{top: true, bottom: true}, scroll)))
}

func withAlpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}

type underLayout struct {
	offset float32
}

func (u underLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	under := objs[0]
	left := size.Width/2 - under.Size().Width/2
	under.Move(fyne.NewPos(left, u.offset-50))
}

func (u underLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{}
}

type unpad struct {
	top, bottom bool
}

func (u unpad) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	pad := theme.Padding()
	var pos fyne.Position
	if u.top {
		pos = fyne.NewPos(0, -pad)
	}
	size := s
	if u.top {
		size = size.AddWidthHeight(0, pad)
	}
	if u.bottom {
		size = size.AddWidthHeight(0, pad)
	}
	for _, o := range objs {
		o.Move(pos)
		o.Resize(size)
	}
}

func (u unpad) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}

func formatAuthors(lines string) string {
	markdown := &strings.Builder{}
	markdown.WriteString("### Authors\n\n")

	for _, line := range strings.Split(lines, "\n") {
		if len(line) == 0 {
			continue
		}

		markdown.WriteString("* ")
		markdown.WriteString(line)
		markdown.WriteByte('\n')
	}

	return markdown.String()
}
