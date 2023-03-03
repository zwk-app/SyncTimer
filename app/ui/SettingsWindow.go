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
var settingsTargetsJsonEntry *widget.Entry

func SettingsWindowOnClose() {
	log.Printf("SettingsWindowOnClose")
	ShowMainWindow()
}

func SettingsToolbarMenuButtonOnClick() {
	log.Printf("SettingsToolbarMenuButtonOnClick")
	SettingsWindowOnClose()
}

func SettingsToolbarHelpButtonOnClick() {
	log.Printf("SettingsToolbarHelpButtonOnClick")

}

func SettingsAlertSoundSelectOnChange(alertTitle string) {
	log.Printf("SettingsAlertSoundSelectOnChange: %s", alertTitle)
	//goland:noinspection GoUnhandledErrorResult
	go appEngine.Audio.Play(appEngine.AlarmSoundName(alertTitle))
}

func SettingsSaveButtonOnClick() {
	log.Printf("SettingsSaveButtonOnClick")
	appEngine.Config.Alerts.TextToSpeech = settingsVoiceAlertsCheck.Checked
	appEngine.Config.Alerts.Notifications = settingsNotificationsCheck.Checked
	appEngine.Config.Alerts.AlarmSound = appEngine.AlarmSoundName(settingsAlertSoundSelect.Selected)
	appEngine.SetTargetJson(settingsTargetsJsonEntry.Text)
	_ = appEngine.Config.SaveFyneSettings()
	SettingsWindowOnClose()
}

func SettingsWindowContent() *fyne.Container {
	log.Printf("SettingsWindowContent")

	if !settingsWindowInitialized {

		/* Top toolbar */
		settingsToolbarMenuButton = widget.NewToolbarAction(theme.NavigateBackIcon(), SettingsToolbarMenuButtonOnClick)
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
		settingsAlertSoundSelect = widget.NewSelect(appEngine.AlarmSoundTitles(), SettingsAlertSoundSelectOnChange)
		settingsAlertSoundForm := widget.NewFormItem("Alert Sound", settingsAlertSoundSelect)
		settingsTargetsJsonEntry = widget.NewEntry()
		settingsTargetsJsonEntry = widget.NewMultiLineEntry()
		settingsTargetsJsonEntry.SetMinRowsVisible(3)
		settingsTargetsJsonEntry.SetText(appEngine.Config.Target.JsonName)
		settingsTargetsJsonForm := widget.NewFormItem("Targets JSON", settingsTargetsJsonEntry)
		settingsForm := widget.NewForm(settingsVoiceAlertsForm, settingsNotificationsForm, settingsAlertSoundForm, settingsTargetsJsonForm)

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
	appEngine.FyneWindow.SetContent(SettingsWindowContent())
	settingsVoiceAlertsCheck.SetChecked(appEngine.Config.Alerts.TextToSpeech)
	settingsNotificationsCheck.SetChecked(appEngine.Config.Alerts.Notifications)
	settingsAlertSoundSelect.SetSelected(appEngine.Config.Alerts.AlarmSound)
	appEngine.FyneWindow.Show()
}
