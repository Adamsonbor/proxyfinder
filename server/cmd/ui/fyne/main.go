package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/proxy"
)

func NewProxyTable() *widget.Table {
	length := func() (rows int, cols int) {
		return 10, 10
	}
	create := func () fyne.CanvasObject {
		return widget.NewLabel("test")
	}
	update := func (id widget.TableCellID, obj fyne.CanvasObject) {
		switch id.Col {
		case 0:
		}
	}

	return widget.NewTable(length, create, update)
}

func main() {
	app := app.New()
	window := app.NewWindow("Proxyfinder")
	window.Resize(fyne.NewSize(800, 600))

	// navbar := container.NewStack(canvas.NewRectangle(color.RGBA{255, 33, 33, 255}))
	navbar := widget.NewLabel("test")

	// proxyTable := container.NewStack(canvas.NewRectangle(color.RGBA{33, 255, 33, 255}))
	proxyTable := NewProxyTable()

	content := container.NewHSplit(
		navbar,
		proxyTable,
	)
	content.SetOffset(0.3)


	window.SetContent(content)
	window.ShowAndRun()
}
