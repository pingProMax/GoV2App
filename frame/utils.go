package frame

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// 间距
func Spacing(size float32) *fyne.Container {
	return container.NewGridWrap(
		fyne.NewSize(size, size),
	)
}
