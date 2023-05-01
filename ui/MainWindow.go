package ui

import (
	"SyncTimer/config"
	"SyncTimer/resources"
	"SyncTimer/timer"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/zwk-app/zwk-tools/logs"
	"image/color"
	"time"
)

var mainWindowInitialized = false
var mainWindowContainer *fyne.Container
var mainToolbarMenu *widget.Toolbar
var mainToolbarMenuButton *widget.ToolbarAction
var mainToolbarTimezoneButton *widget.ToolbarAction
var mainToolbarHelpButton *widget.ToolbarAction
var currentColor color.NRGBA
var currentLabel *canvas.Text
var currentText *canvas.Text
var targetColor color.NRGBA
var targetLabel *canvas.Text
var targetText *canvas.Text
var remainingColor color.NRGBA
var remainingLabel *canvas.Text
var remainingText *canvas.Text
var locationColor color.NRGBA
var locationText *canvas.Text

func MainToolbarMenuButtonOnClick() {
	logs.Debug("MainWindow", fmt.Sprintf("MenuButtonOnClick"), nil)
	ShowSettingsWindow()
}

func MainToolbarHelpButtonOnClick() {
	logs.Debug("MainWindow", fmt.Sprintf("HelpButtonOnClick"), nil)
	TextToSpeechAlert("about")
}

func RefreshDisplayedTimezone() {
	if timer.LocationName() == config.LocalLocationName {
		mainToolbarTimezoneButton.SetIcon(theme.HomeIcon())
	} else {
		mainToolbarTimezoneButton.SetIcon(theme.MediaRecordIcon())
	}
	locationText.Text = timer.LocationName()
	locationText.Refresh()
	mainToolbarMenu.Refresh()
}

func MainToolbarTimezoneButtonOnClick() {
	logs.Debug("MainWindow", fmt.Sprintf("TimezoneButtonOnClick"), nil)
	if locationText.Text == "UTC" {
		timer.SetLocationName(config.LocalLocationName)
	} else {
		timer.SetLocationName("UTC")
	}
	config.Location().Name = timer.LocationName()
	config.SaveFyneSettings(FyneApp)
	RefreshDisplayedTimezone()
}

func MainWindowLoop() {
	logs.Debug("MainWindow", fmt.Sprintf("Loop"), nil)
	go func() {
		for {
			/* default width: 320 */
			cw := FyneWindow.Content().Size().Width
			//goland:noinspection GoRedundantConversion
			r := float32(cw) / float32(320)

			targetLabel.TextSize = 16 * r
			currentLabel.TextSize = 16 * r
			remainingLabel.TextSize = 16 * r

			targetLabel.Refresh()
			currentLabel.Refresh()
			remainingLabel.Refresh()

			targetText.TextSize = 36 * r
			currentText.TextSize = 36 * r
			remainingText.TextSize = 36 * r

			targetLabel.Text = timer.AlarmName()
			targetText.Text = timer.TargetTimeText()
			currentText.Text = timer.CurrentTimeText()
			remainingText.Text = timer.RemainingTimeText()

			if timer.RemainingSeconds() < -30 {
				timer.NextTarget()
			}

			targetText.Refresh()
			currentText.Refresh()
			remainingText.Refresh()

			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func MainWindowContent() *fyne.Container {
	logs.Debug("MainWindow", fmt.Sprintf("Content"), nil)

	if !mainWindowInitialized {

		currentColor = color.NRGBA{R: 100, G: 100, B: 150, A: 255}
		targetColor = color.NRGBA{R: 100, G: 150, B: 100, A: 255}
		remainingColor = color.NRGBA{R: 150, G: 100, B: 100, A: 255}
		locationColor = color.NRGBA{R: 100, G: 100, B: 100, A: 255}

		/* Top toolbar */
		mainToolbarMenuButton = widget.NewToolbarAction(theme.SettingsIcon(), MainToolbarMenuButtonOnClick)
		mainToolbarTimezoneButton = widget.NewToolbarAction(theme.HomeIcon(), MainToolbarTimezoneButtonOnClick)
		mainToolbarHelpButton = widget.NewToolbarAction(theme.HelpIcon(), MainToolbarHelpButtonOnClick)
		mainToolbarMenu = widget.NewToolbar(
			mainToolbarMenuButton,
			widget.NewToolbarSeparator(),
			mainToolbarTimezoneButton,
			widget.NewToolbarSpacer(),
			mainToolbarHelpButton,
		)

		/* Current Time */
		currentLabel = canvas.NewText("Current", currentColor)
		currentLabel.Alignment = fyne.TextAlignCenter
		currentLabel.TextSize = 16
		currentLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		currentText = canvas.NewText(timer.CurrentTimeText(), currentColor)
		currentText.Alignment = fyne.TextAlignCenter
		currentText.TextSize = 36
		currentText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		currentGrid := container.NewGridWithColumns(1, currentLabel, currentText)

		/* Target Time */
		targetLabel = canvas.NewText(timer.AlarmName(), targetColor)
		targetLabel.Alignment = fyne.TextAlignCenter
		targetLabel.TextSize = 16
		targetLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		targetText = canvas.NewText(timer.TargetTimeText(), targetColor)
		targetText.Alignment = fyne.TextAlignCenter
		targetText.TextSize = 36
		targetText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		targetGrid := container.NewGridWithColumns(1, targetLabel, targetText)

		/* Remaining Time */
		remainingLabel = canvas.NewText("Remaining", remainingColor)
		remainingLabel.Alignment = fyne.TextAlignCenter
		remainingLabel.TextSize = 16
		remainingLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		remainingText = canvas.NewText(timer.RemainingTimeText(), remainingColor)
		remainingText.Alignment = fyne.TextAlignCenter
		remainingText.TextSize = 36
		remainingText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		remainingGrid := container.NewGridWithColumns(1, remainingLabel, remainingText)

		/* Current Location (UTC/Local) */
		locationText = canvas.NewText(timer.LocationName(), locationColor)
		locationText.Alignment = fyne.TextAlignCenter
		locationText.TextSize = 16
		locationText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		locationGrid := container.NewGridWithColumns(1, locationText)

		/* Middle content */
		middleGrid := container.NewGridWithColumns(1, currentGrid, targetGrid, remainingGrid, locationGrid)
		middleContainer := container.NewCenter(middleGrid)

		/* Bottom buttons */
		setTargetDelayIcon := fyne.NewStaticResource("delay.svg", resources.ReadImage("delay.svg"))
		setTargetDelayButton := widget.NewButtonWithIcon("Set Delay", setTargetDelayIcon, ShowTargetDelayWindow)
		setTargetTimeIcon := fyne.NewStaticResource("target.svg", resources.ReadImage("target.svg"))
		setTargetTimeButton := widget.NewButtonWithIcon("Set Target", setTargetTimeIcon, ShowTargetTimeWindow)
		bottomContainer := container.NewGridWithColumns(2, setTargetDelayButton, setTargetTimeButton)

		RefreshDisplayedTimezone()

		MainWindowLoop()

		mainWindowContainer = container.NewBorder(mainToolbarMenu, bottomContainer, nil, nil, middleContainer)

		mainWindowInitialized = true
	}
	return mainWindowContainer
}

func ShowMainWindow() {
	logs.Debug("MainWindow", fmt.Sprintf("Show"), nil)
	FyneWindow.SetContent(MainWindowContent())
	FyneWindow.Show()
}
