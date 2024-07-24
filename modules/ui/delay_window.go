package ui

//go:generate fyne bundle --package resources --name ResourceDivineIco -o resources/icon.go divine.ico

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DelayWindow struct {
	fyne.Window

	delayEntry    *widget.Entry
	confirmButton *widget.Button
	cancelButton  *widget.Button
	linkID        int
}

func NewDelayWindow(app fyne.App, delay int64, linkId int) *DelayWindow {

	dw := &DelayWindow{linkID: linkId}

	DelayWindow := app.NewWindow("Delay change")
	DelayWindow.SetFixedSize(true)
	DelayWindow.Resize(fyne.NewSize(405, 150))
	DelayWindow.CenterOnScreen()
	dw.Window = DelayWindow

	delayEntry := widget.NewEntry()
	dw.delayEntry = delayEntry
	delayEntry.Move(fyne.NewPos(10, 20))
	delayEntry.Resize(fyne.NewSize(380, 40))
	delayEntry.SetPlaceHolder("Enter delay in milliseconds")
	delayEntry.SetText(textFromDelay(delay))
	delayEntry.Validator = validateDelay
	delayEntry.Refresh()

	confirmButton := widget.NewButton("OK", nil)
	dw.confirmButton = confirmButton
	confirmButton.Move(fyne.NewPos(210, 80))
	confirmButton.Resize(fyne.NewSize(180, 40))

	cancelButton := widget.NewButton("Cancel", func() { dw.Close() })
	dw.cancelButton = cancelButton
	cancelButton.Move(fyne.NewPos(10, 80))
	cancelButton.Resize(fyne.NewSize(180, 40))

	dw.SetContent(container.NewWithoutLayout(
		delayEntry,
		confirmButton,
		cancelButton,
	))

	return dw
}

func (w *DelayWindow) OnConfirmDelay(f func()) {
	w.confirmButton.OnTapped = f
}

func textFromDelay(delay int64) string {
	if delay == 0 {
		return ""
	} else {
		return fmt.Sprint(delay)
	}
}

func validateDelay(s string) error {
	if s == "" {
		return nil
	}
	_, err := strconv.Atoi(s)
	return err
}
