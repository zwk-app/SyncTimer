package ui

import (
	"SyncTimer/timer"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zwk-app/zwk-tools/logs"
)

var targetTimeWindowInitialized = false
var targetTimeWindowContainer *fyne.Container
var targetTimeInput *NumbersEntry

func TargetTimeWindowOnClose() {
	logs.Debug("TargetTimeWindow", "OnClose", nil)
	ShowMainWindow()
}

func TargetTimeInputValidator(s string) error {
	if !timer.CheckTimeString(s) {
		return fmt.Errorf("invalid")
	}
	return nil
}

func TargetTimeInputOnSubmitted(s string) {
	logs.Debug("TargetTimeWindow", fmt.Sprintf("InputOnSubmitted: '%s'", s), nil)
	if TargetTimeInputValidator(s) == nil {
		TargetTimeWindowConfirmButtonOnClick()
	}
}

func TargetTimeWindowConfirmButtonOnClick() {
	logs.Debug("TargetTimeWindow", "ConfirmButtonOnClick", nil)
	s := targetTimeInput.Text
	if TargetTimeInputValidator(s) == nil {
		timer.SetTargetTimeString(s)
		TargetTimeWindowOnClose()
	}
}

func TargetTimeWindowContent() *fyne.Container {
	logs.Debug("TargetTimeWindow", "Content", nil)

	if !targetTimeWindowInitialized {

		/* Middle content */
		targetTimeInput = NewNumbersEntry()
		targetTimeInput.Validator = TargetTimeInputValidator
		targetTimeInput.OnSubmitted = TargetTimeInputOnSubmitted
		targetTimeInput.SetPlaceHolder("[[hh]mm]ss")
		targetTimeInput.TextStyle.Bold = true
		targetTimeInput.TextStyle.Monospace = true
		targetTimeContainer := container.NewCenter(targetTimeInput)

		/* Bottom buttons */
		targetTimeCancelButton := widget.NewButton("Cancel", TargetTimeWindowOnClose)
		targetTimeCancelButton.SetIcon(theme.CancelIcon())
		targetTimeConfirmButton := widget.NewButton("Confirm", TargetTimeWindowConfirmButtonOnClick)
		targetTimeConfirmButton.SetIcon(theme.ConfirmIcon())
		bottomContainer := container.NewGridWithRows(1, targetTimeCancelButton, targetTimeConfirmButton)

		targetTimeWindowContainer = container.NewBorder(nil, bottomContainer, nil, nil, targetTimeContainer)

		targetTimeWindowInitialized = true
	}
	return targetTimeWindowContainer
}
func ShowTargetTimeWindow() {
	logs.Debug("TargetTimeWindow", "Show", nil)
	FyneWindow.SetContent(TargetTimeWindowContent())
	targetTimeInput.SetText("")
	FyneWindow.Show()
	FyneWindow.Canvas().Focus(targetTimeInput)
}
