package ui

import (
	"SyncTimer/timer"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zwk-app/go-tools/logs"
)

var targetDelayWindowInitialized = false
var targetDelayWindowContainer *fyne.Container
var targetDelayInput *NumbersEntry

func TargetDelayWindowOnClose() {
	logs.Debug("TargetDelayWindow", "OnClose", nil)
	ShowMainWindow()
}

func TargetDelayInputValidator(s string) error {
	if !timer.CheckDelayString(s) {
		return fmt.Errorf("invalid")
	}
	return nil
}

func TargetDelayInputOnSubmitted(s string) {
	logs.Debug("TargetDelayWindow", fmt.Sprintf("InputOnSubmitted: '%s'", s), nil)
	if TargetDelayInputValidator(s) == nil {
		TargetDelayWindowConfirmButtonOnClick()
	}
}

func TargetDelayWindowConfirmButtonOnClick() {
	logs.Debug("TargetDelayWindow", "ConfirmButtonOnClick", nil)
	s := targetDelayInput.Text
	if TargetDelayInputValidator(s) == nil {
		timer.SetTargetDelayString(s)
		TargetDelayWindowOnClose()
	}
}

func TargetDelayWindowContent() *fyne.Container {
	logs.Debug("TargetDelayWindow", "Content", nil)

	if !targetDelayWindowInitialized {

		/* Middle content */
		targetDelayInput = NewNumbersEntry()
		targetDelayInput.Validator = TargetDelayInputValidator
		targetDelayInput.OnSubmitted = TargetDelayInputOnSubmitted
		targetDelayInput.SetPlaceHolder("[[hh]mm]ss")
		targetDelayInput.TextStyle.Bold = true
		targetDelayInput.TextStyle.Monospace = true
		targetDelayContainer := container.NewCenter(targetDelayInput)

		/* Bottom buttons */
		targetDelayCancelButton := widget.NewButton("Cancel", TargetDelayWindowOnClose)
		targetDelayCancelButton.SetIcon(theme.CancelIcon())
		targetDelayConfirmButton := widget.NewButton("Confirm", TargetDelayWindowConfirmButtonOnClick)
		targetDelayConfirmButton.SetIcon(theme.ConfirmIcon())
		bottomContainer := container.NewGridWithRows(1, targetDelayCancelButton, targetDelayConfirmButton)

		targetDelayWindowContainer = container.NewBorder(nil, bottomContainer, nil, nil, targetDelayContainer)

		targetDelayWindowInitialized = true
	}
	return targetDelayWindowContainer
}
func ShowTargetDelayWindow() {
	logs.Debug("TargetDelayWindow", "Show", nil)
	FyneWindow.SetContent(TargetDelayWindowContent())
	targetDelayInput.SetText("")
	FyneWindow.Show()
	FyneWindow.Canvas().Focus(targetDelayInput)
}
