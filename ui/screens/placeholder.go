package screens

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// newPlaceholderScreen returns a centred label for screens not yet implemented.
func newPlaceholderScreen(title string) fyne.CanvasObject {
	label := widget.NewLabel(title + " — coming soon")
	return container.NewCenter(label)
}
