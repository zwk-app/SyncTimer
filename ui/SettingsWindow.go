package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
)

var settingsWindowInitialized = false
var settingsWindowContainer *fyne.Container
var settingsToolbarMenu *widget.Toolbar
var settingsToolbarMenuButton *widget.ToolbarAction
var settingsToolbarHelpButton *widget.ToolbarAction
var settingsVoiceAlertsCheck *widget.Check
var settingsNotificationsCheck *widget.Check
var settingsAlertSoundSelect *widget.Select

func SettingsOnClose() {
	log.Printf("SettingsOnClose")
	ShowMainWindow()
}

func SettingsToolbarMenuButtonOnClick() {
	log.Printf("SettingsToolbarMenuButtonOnClick")
	SettingsOnClose()
}

func SettingsToolbarHelpButtonOnClick() {
	log.Printf("SettingsToolbarHelpButtonOnClick")

}

func SettingsAlertSoundSelectOnChange(alertTitle string) {
	log.Printf("SettingsAlertSoundSelectOnChange: %s", alertTitle)
	//goland:noinspection GoUnhandledErrorResult
	go appEngine.Audio.Object.Play(appEngine.AlertName(alertTitle))
}

func SettingsSaveButtonOnClick() {
	log.Printf("SettingsSaveButtonOnClick")
	appEngine.Alerts.TextToSpeech = settingsVoiceAlertsCheck.Checked
	appEngine.Alerts.Notifications = settingsNotificationsCheck.Checked
	appEngine.Alerts.AlertSound = appEngine.AlertName(settingsAlertSoundSelect.Selected)
	_ = appEngine.SaveFyneSettings()
	SettingsOnClose()
}

func SettingsWindowContent() *fyne.Container {
	log.Printf("SettingsWindowContent")

	if !settingsWindowInitialized {

		/* Top toolbar */
		settingsToolbarMenuButton = widget.NewToolbarAction(theme.MenuIcon(), SettingsToolbarMenuButtonOnClick)
		settingsToolbarHelpButton = widget.NewToolbarAction(theme.HelpIcon(), SettingsToolbarHelpButtonOnClick)
		settingsToolbarMenu = widget.NewToolbar(
			settingsToolbarMenuButton,
			widget.NewToolbarSeparator(),
			widget.NewToolbarSpacer(),
			settingsToolbarHelpButton,
		)

		/* Middle content */
		settingsVoiceAlertsCheck = widget.NewCheck("", func(b bool) {})
		settingsVoiceAlertsForm := widget.NewFormItem("Voice Alerts", settingsVoiceAlertsCheck)
		settingsNotificationsCheck = widget.NewCheck("", func(b bool) {})
		settingsNotificationsForm := widget.NewFormItem("Notifications", settingsNotificationsCheck)
		settingsAlertSoundSelect = widget.NewSelect(appEngine.Alerts.AlertSoundTitles, SettingsAlertSoundSelectOnChange)
		settingsAlertSoundForm := widget.NewFormItem("Alert Sound", settingsAlertSoundSelect)
		settingsForm := widget.NewForm(settingsVoiceAlertsForm, settingsNotificationsForm, settingsAlertSoundForm)

		/* Bottom buttons */
		saveSettingsButton := widget.NewButton("Save", SettingsSaveButtonOnClick)
		saveSettingsButton.SetIcon(theme.DocumentSaveIcon())
		bottomContainer := container.NewCenter(saveSettingsButton)

		settingsWindowContainer = container.NewBorder(settingsToolbarMenu, bottomContainer, nil, nil, settingsForm)

		settingsWindowInitialized = true
	}
	return settingsWindowContainer
}
func ShowSettingsWindow() {
	log.Printf("ShowSettingsWindow")
	appEngine.Fyne.MainWindow.SetContent(SettingsWindowContent())
	settingsVoiceAlertsCheck.SetChecked(appEngine.Alerts.TextToSpeech)
	settingsNotificationsCheck.SetChecked(appEngine.Alerts.Notifications)
	settingsAlertSoundSelect.SetSelected(appEngine.Alerts.AlertSound)
	appEngine.Fyne.MainWindow.Show()
}
