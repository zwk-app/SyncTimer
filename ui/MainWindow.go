package ui

import (
	"SyncTimer/timer"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"time"
)

var mainWindowInitialized = false
var mainWindow fyne.Window
var toolbarMenu *widget.Toolbar
var toolbarMenuButton *widget.ToolbarAction
var toolbarTimezoneButton *widget.ToolbarAction
var toolbarHelpButton *widget.ToolbarAction
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

func ToolbarMenuButtonOnClick() {
	log.Println("ToolbarMenuButtonOnClick")
	OptionsDialogShow()
}

func ToolbarHelpButtonOnClick() {
	log.Println("ToolbarHelpButtonOnClick")
	TextToSpeechAlert("about")
}

func RefreshDisplayedTimezone() {
	if appEngine.Timer.Object.GetLocationName() == timer.LocalLocationName {
		toolbarTimezoneButton.SetIcon(theme.HomeIcon())
	} else {
		toolbarTimezoneButton.SetIcon(theme.MediaRecordIcon())
	}
	locationText.Text = appEngine.Timer.Object.GetLocationName()
	locationText.Refresh()
	toolbarMenu.Refresh()
}

func ToolbarTimezoneButtonOnClick() {
	log.Println("ToolbarTimezoneButtonOnClick")
	if locationText.Text == "UTC" {
		appEngine.Timer.Object.SetLocationName(timer.LocalLocationName)
	} else {
		appEngine.Timer.Object.SetLocationName("UTC")
	}
	appEngine.Timer.LocationName = appEngine.Timer.Object.GetLocationName()
	_ = appEngine.SaveFyneSettings()
	RefreshDisplayedTimezone()
}

func MainWindowLoop() {
	log.Println("MainWindowLoop : Start")
	go func() {
		for {
			/* default width: 320 */
			cw := mainWindow.Content().Size().Width
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

			targetText.Text = appEngine.Timer.Object.GetTargetTimeString()
			currentText.Text = appEngine.Timer.Object.GetCurrentTimeString()
			remainingText.Text = appEngine.Timer.Object.GetRemainingString()

			if appEngine.Timer.Object.GetRemainingSeconds() < -30 {
				appEngine.Timer.Object.Next()
			}

			targetText.Refresh()
			currentText.Refresh()
			remainingText.Refresh()

			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func MainWindowInit() {

	log.Println("MainWindowInit")

	if !mainWindowInitialized {

		currentColor = color.NRGBA{R: 100, G: 100, B: 150, A: 255}
		targetColor = color.NRGBA{R: 100, G: 150, B: 100, A: 255}
		remainingColor = color.NRGBA{R: 150, G: 100, B: 100, A: 255}
		locationColor = color.NRGBA{R: 100, G: 100, B: 100, A: 255}

		/* Main Window creation  */
		mainWindow = appEngine.Fyne.App.NewWindow(appEngine.Title())
		mainWindow.Resize(fyne.NewSize(320, 540))

		/* Top toolbar */
		toolbarMenuButton = widget.NewToolbarAction(theme.SettingsIcon(), ToolbarMenuButtonOnClick)
		toolbarTimezoneButton = widget.NewToolbarAction(theme.HomeIcon(), ToolbarTimezoneButtonOnClick)
		toolbarHelpButton = widget.NewToolbarAction(theme.HelpIcon(), ToolbarHelpButtonOnClick)
		toolbarMenu = widget.NewToolbar(
			toolbarMenuButton,
			widget.NewToolbarSeparator(),
			toolbarTimezoneButton,
			widget.NewToolbarSpacer(),
			toolbarHelpButton,
		)

		/* Current Time */
		currentLabel = canvas.NewText("Current", currentColor)
		currentLabel.Alignment = fyne.TextAlignCenter
		currentLabel.TextSize = 16
		currentLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		currentText = canvas.NewText(appEngine.Timer.Object.GetCurrentTimeString(), currentColor)
		currentText.Alignment = fyne.TextAlignCenter
		currentText.TextSize = 36
		currentText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		currentGrid := container.New(layout.NewGridLayout(1), currentLabel, currentText)

		/* Target Time */
		targetLabel = canvas.NewText("Target", targetColor)
		targetLabel.Alignment = fyne.TextAlignCenter
		targetLabel.TextSize = 16
		targetLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		targetText = canvas.NewText(appEngine.Timer.Object.GetTargetTimeString(), targetColor)
		targetText.Alignment = fyne.TextAlignCenter
		targetText.TextSize = 36
		targetText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		targetGrid := container.New(layout.NewGridLayout(1), targetLabel, targetText)

		/* Remaining Time */
		remainingLabel = canvas.NewText("Remaining", remainingColor)
		remainingLabel.Alignment = fyne.TextAlignCenter
		remainingLabel.TextSize = 16
		remainingLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		remainingText = canvas.NewText(appEngine.Timer.Object.GetRemainingString(), remainingColor)
		remainingText.Alignment = fyne.TextAlignCenter
		remainingText.TextSize = 36
		remainingText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		remainingGrid := container.New(layout.NewGridLayout(1), remainingLabel, remainingText)

		/* Current Location (UTC/Local) */
		locationText = canvas.NewText(appEngine.Timer.Object.GetLocationName(), locationColor)
		locationText.Alignment = fyne.TextAlignCenter
		locationText.TextSize = 16
		locationText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		locationGrid := container.New(layout.NewGridLayout(1), locationText)

		/* Middle content */
		middleGrid := container.New(layout.NewGridLayout(1), currentGrid, targetGrid, remainingGrid, locationGrid)
		middleContainer := container.New(layout.NewCenterLayout(), middleGrid)

		/* Bottom buttons */
		setTargetTimeButton := widget.NewButton("Set Target", TargetTimeDialogShow)
		setTargetTimeButton.SetIcon(theme.MediaPlayIcon())
		bottomContainer := container.New(layout.NewCenterLayout(), setTargetTimeButton)

		mainWindow.SetContent(container.NewBorder(toolbarMenu, bottomContainer, nil, nil, middleContainer))

		RefreshDisplayedTimezone()

		MainWindowLoop()

		mainWindowInitialized = true
	}
}

func MainWindowShow() {
	MainWindowInit()
	mainWindow.Show()
}
