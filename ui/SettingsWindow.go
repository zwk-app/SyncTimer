package ui

import (
	"SyncTimer/audio"
	"SyncTimer/config"
	"SyncTimer/logs"
	"SyncTimer/resources"
	"SyncTimer/timer"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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
	logs.Debug("SettingsWindow", fmt.Sprintf("OnClose"), nil)
	ShowMainWindow()
}

func SettingsToolbarMenuButtonOnClick() {
	logs.Debug("SettingsWindow", fmt.Sprintf("MenuButtonOnClick"), nil)
	SettingsWindowOnClose()
}

func SettingsToolbarHelpButtonOnClick() {
	logs.Debug("SettingsWindow", fmt.Sprintf("HelpButtonOnClick"), nil)
}

func SettingsAlertSoundSelectOnChange(alertTitle string) {
	logs.Debug("SettingsWindow", fmt.Sprintf("AlertSoundSelectOnChange: '%s'", alertTitle), nil)
	go audio.Play(resources.AlarmSoundName(alertTitle))
}

func SettingsSaveButtonOnClick() {
	logs.Debug("SettingsWindow", fmt.Sprintf("SaveButtonOnClick"), nil)
	config.Alerts().TextToSpeech = settingsVoiceAlertsCheck.Checked
	config.Alerts().Notifications = settingsNotificationsCheck.Checked
	config.Alerts().AlarmSound = resources.AlarmSoundName(settingsAlertSoundSelect.Selected)
	timer.SetTargetJson(settingsTargetsJsonEntry.Text)
	config.SaveFyneSettings(FyneApp)
	SettingsWindowOnClose()
}

func SettingsWindowContent() *fyne.Container {
	logs.Debug("SettingsWindow", fmt.Sprintf("Content"), nil)

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
		settingsAlertSoundSelect = widget.NewSelect(resources.AlarmSoundTitles(), SettingsAlertSoundSelectOnChange)
		settingsAlertSoundForm := widget.NewFormItem("Alert Sound", settingsAlertSoundSelect)
		settingsTargetsJsonEntry = widget.NewEntry()
		settingsTargetsJsonEntry = widget.NewMultiLineEntry()
		settingsTargetsJsonEntry.SetMinRowsVisible(3)
		settingsTargetsJsonEntry.SetText(config.Target().JsonName)
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
	logs.Debug("SettingsWindow", fmt.Sprintf("Show"), nil)
	FyneWindow.SetContent(SettingsWindowContent())
	settingsVoiceAlertsCheck.SetChecked(config.Alerts().TextToSpeech)
	settingsNotificationsCheck.SetChecked(config.Alerts().Notifications)
	settingsAlertSoundSelect.SetSelected(config.Alerts().AlarmSound)
	FyneWindow.Show()
}
