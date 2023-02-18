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
		appEngine.Alerts.TextToSpeech = voiceAlertsCheck.Checked
		appEngine.Alerts.Notifications = notificationsCheck.Checked
		_ = appEngine.SaveFyneSettings()
	}
}

func OptionsDialogShow() {
	log.Println("OptionsDialogShow")
	OptionsDialogInit()
	voiceAlertsCheck.SetChecked(appEngine.Alerts.TextToSpeech)
	notificationsCheck.SetChecked(appEngine.Alerts.Notifications)
	dialog.ShowForm("Options", "OK", "Cancel",
		optionsDialogForm, OptionsDialogOnExit, mainWindow)
}
