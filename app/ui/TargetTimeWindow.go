package ui

import (
	"SyncTimer/app/timer"
	"fmt"
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

func TargetInputValidator(s string) error {
	if !timer.CheckTimeString(s) {
		return fmt.Errorf("invalid")
	}
	return nil
}

func TargetInputOnSubmitted(s string) {
	log.Printf("TargetInputOnSubmitted %s", s)
	if TargetInputValidator(s) == nil {
		TargetWindowConfirmButtonOnClick()
	}
}

func TargetWindowConfirmButtonOnClick() {
	log.Printf("TargetWindowConfirmButtonOnClick")
	s := targetInput.Text
	if TargetInputValidator(s) == nil {
		appEngine.SetTargetTime(s)
		TargetWindowOnClose()
	}
}

func TargetWindowContent() *fyne.Container {
	log.Printf("TargetWindowContent")

	if !targetWindowInitialized {

		/* Middle content */
		targetInput = NewNumbersEntry()
		targetInput.Validator = TargetInputValidator
		targetInput.OnSubmitted = TargetInputOnSubmitted
		targetInput.SetPlaceHolder("hh[mm[ss]]")
		targetInput.TextStyle.Bold = true
		targetInput.TextStyle.Monospace = true
		targetContainer := container.NewCenter(targetInput)

		/* Bottom buttons */
		targetCancelButton := widget.NewButton("Cancel", TargetWindowOnClose)
		targetCancelButton.SetIcon(theme.CancelIcon())
		targetConfirmButton := widget.NewButton("Confirm", TargetWindowConfirmButtonOnClick)
		targetConfirmButton.SetIcon(theme.ConfirmIcon())
		bottomContainer := container.NewGridWithRows(1, targetCancelButton, targetConfirmButton)

		targetWindowContainer = container.NewBorder(nil, bottomContainer, nil, nil, targetContainer)

		targetWindowInitialized = true
	}
	return targetWindowContainer
}
func ShowTargetWindow() {
	log.Printf("ShowTargetWindow")
	appEngine.FyneWindow.SetContent(TargetWindowContent())
	targetInput.SetText("")
	appEngine.FyneWindow.Show()
	appEngine.FyneWindow.Canvas().Focus(targetInput)
}
