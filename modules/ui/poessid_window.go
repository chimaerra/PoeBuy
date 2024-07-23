package ui

//go:generate fyne bundle --package resources --name ResourceDivineIco -o resources/icon.go divine.ico

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PoessidWindow struct {
	fyne.Window

	poesessidEntry *widget.Entry
	confirmButton  *widget.Button
	cancelButton   *widget.Button
}

func NewPoessidWindow(app fyne.App) *PoessidWindow {

	pw := &PoessidWindow{}

	PoessidWindow := app.NewWindow("PoeBuy")
	PoessidWindow.SetFixedSize(true)
	PoessidWindow.Resize(fyne.NewSize(405, 150))
	pw.Window = PoessidWindow

	poesessidEntry := widget.NewEntry()
	pw.poesessidEntry = poesessidEntry
	poesessidEntry.Move(fyne.NewPos(10, 20))
	poesessidEntry.Resize(fyne.NewSize(380, 40))
	poesessidEntry.SetPlaceHolder("Enter your POESSID here")
	poesessidEntry.Password = true
	poesessidEntry.TextStyle = fyne.TextStyle{Bold: true, Italic: false, Monospace: false}
	poesessidEntry.Refresh()

	confirmButton := widget.NewButton("OK", nil)
	pw.confirmButton = confirmButton
	confirmButton.Move(fyne.NewPos(210, 80))
	confirmButton.Resize(fyne.NewSize(180, 40))

	cancelButton := widget.NewButton("Cancel", nil)
	pw.cancelButton = cancelButton
	cancelButton.Move(fyne.NewPos(10, 80))
	cancelButton.Resize(fyne.NewSize(180, 40))

	pw.SetContent(container.NewWithoutLayout(
		poesessidEntry,
		confirmButton,
		cancelButton,
	))

	return pw
}

func (w *PoessidWindow) OnConfirmPoessid(f func()) {
	w.confirmButton.OnTapped = f
}

func (w *PoessidWindow) OnClose(f func()) {
	w.cancelButton.OnTapped = f
}
