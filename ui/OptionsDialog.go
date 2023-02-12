package ui

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
)

var optionsDialogInitialized = false
var optionsDialogForm []*widget.FormItem
var voiceAlertsCheck *widget.Check
var notificationsCheck *widget.Check

func OptionsDialogInit() {
	log.Println("OptionsDialogInit")
	if !optionsDialogInitialized {
		voiceAlertsCheck = widget.NewCheck("", func(b bool) {})
		notificationsCheck = widget.NewCheck("", func(b bool) {})
		optionsDialogForm = append(optionsDialogForm, widget.NewFormItem("Voice alerts", voiceAlertsCheck))
		optionsDialogForm = append(optionsDialogForm, widget.NewFormItem("Notifications", notificationsCheck))
		optionsDialogInitialized = true
	}
}

func OptionsDialogOnExit(b bool) {
	log.Println("OptionsDialogOnExit")
	if b {
		voiceAlertsEnabled = voiceAlertsCheck.Checked
		notificationsEnabled = notificationsCheck.Checked
		mainApp.Preferences().SetBool("voiceAlertsEnabled", voiceAlertsEnabled)
		mainApp.Preferences().SetBool("notificationsEnabled", notificationsEnabled)
	}
}

func OptionsDialogShow() {
	log.Println("OptionsDialogShow")
	OptionsDialogInit()
	voiceAlertsCheck.SetChecked(mainApp.Preferences().BoolWithFallback("voiceAlertsEnabled", voiceAlertsEnabled))
	notificationsCheck.SetChecked(mainApp.Preferences().BoolWithFallback("notificationsEnabled", notificationsEnabled))
	dialog.ShowForm("Options", "OK", "Cancel",
		optionsDialogForm, OptionsDialogOnExit, mainWindow)
}
