package ui

import (
	"SyncTimer/timer"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

var targetWindowInitialized = false
var targetWindowContainer *fyne.Container
var targetInput *NumbersEntry

func TargetWindowOnClose() {
	log.Printf("TargetWindowOnClose")
	ShowMainWindow()
}

func TargetConfirmButtonOnClick() {
	log.Printf("TargetConfirmButtonOnClick")
	appEngine.SetTargetTime(targetInput.Text)
	TargetWindowOnClose()
}

func TargetInputOnSubmitted(s string) {
	log.Printf("TargetInputOnSubmitted %s", s)
	if timer.CheckTimeString(s) {
		TargetConfirmButtonOnClick()
	}
}

func TargetWindowContent() *fyne.Container {
	log.Printf("TargetWindowContent")

	if !targetWindowInitialized {

		/* Middle content */
		targetInput = NewNumbersEntry()
		targetInput.OnSubmitted = TargetInputOnSubmitted
		targetInput.SetPlaceHolder("hh[mm[ss]]")
		targetContainer := container.NewCenter(targetInput)

		/* Bottom buttons */
		targetCancelButton := widget.NewButton("Cancel", TargetWindowOnClose)
		targetCancelButton.SetIcon(theme.CancelIcon())
		targetConfirmButton := widget.NewButton("Confirm", TargetConfirmButtonOnClick)
		targetConfirmButton.SetIcon(theme.ConfirmIcon())
		bottomContainer := container.NewGridWithRows(1, targetCancelButton, targetConfirmButton)

		targetWindowContainer = container.NewBorder(nil, bottomContainer, nil, nil, targetContainer)

		targetWindowInitialized = true
	}
	return targetWindowContainer
}
func ShowTargetWindow() {
	log.Printf("ShowTargetWindow")
	appEngine.Fyne.MainWindow.SetContent(TargetWindowContent())
	targetInput.SetText("")
	appEngine.Fyne.MainWindow.Show()
	appEngine.Fyne.MainWindow.Canvas().Focus(targetInput)
}
