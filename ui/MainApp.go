package ui

import (
	"SyncTimer/ttm"
	"SyncTimer/tts"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"log"
	"os"
	"time"
)

var mainAppName string
var mainAppVersion string

var mainApp fyne.App
var voiceAlertsEnabled = true
var notificationsEnabled = false

var currentLocation *time.Location
var currentLocationName string
var remainingSeconds int

func TextToSpeechAlert(name string) {
	if voiceAlertsEnabled {
		log.Printf("TextToSpeechAlert '%s'", name)
		tts.TextToSpeechEngine.Play(name)
	}
}

func NotificationAlert(message string) {
	if notificationsEnabled {
		log.Printf("NotificationAlert '%s'", message)
		mainApp.SendNotification(fyne.NewNotification(mainAppName+" v"+mainAppVersion, message))
	}
}

func SetCurrentLocation(locationName string) {
	if ttm.CurrentTargetTime.GetLocationName() != locationName {
		_ = ttm.SetCurrentTargetLocation(locationName)
	}
	switch locationName {
	case "UTC":
		currentLocation = time.UTC
		currentLocationName = "UTC"
		toolbarTimezoneButtonIcon = theme.MediaRecordIcon()
	default:
		currentLocation = time.Local
		currentLocationName = "Local Time"
		toolbarTimezoneButtonIcon = theme.HomeIcon()
	}
}

func TimeLoop() {
	log.Println("TimeLoop : Start")
	go func() {
		for {
			SetCurrentLocation(currentLocationName)
			currentTime := time.Now().In(currentLocation)
			currentString = ttm.GetTimeString(currentTime.Hour(), currentTime.Minute(), currentTime.Second())
			remainingDuration := ttm.CurrentTargetTime.GetTime().In(currentLocation).Sub(currentTime)
			d := remainingDuration.Round(time.Second)
			h := d / time.Hour
			d -= h * time.Hour
			m := d / time.Minute
			d -= m * time.Minute
			s := d / time.Second
			remainingSeconds = int(s) + (int(m) * 60) + (int(h) * 3600)
			if h > 0 {
				remainingString = fmt.Sprintf("%d:%02d:%02d", h, m, s)
			} else if m > 0 {
				remainingString = fmt.Sprintf("%02d:%02d", m, s)
			} else if s > -1 {
				remainingString = fmt.Sprintf("%02d", s)
			} else {
				if remainingString == "" {
					remainingString = fmt.Sprintf("%02d", remainingSeconds)
				} else {
					remainingString = ""
				}
				if remainingSeconds < -30 {
					ttm.SetNextTargetTime()
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func AlertLoop() {
	log.Println("AlertLoop : Start")
	time.Sleep(1500 * time.Millisecond)
	go func() {
		currentCheck := 0
		lastCheck := 0
		lastCheckDiff := 0
		for {
			currentCheck = remainingSeconds
			lastCheckDiff = lastCheck - currentCheck
			if lastCheck < currentCheck {
				log.Printf("AlertLoop : %08d %08d %02d", lastCheck, currentCheck, lastCheckDiff)
				lastCheck = currentCheck + 1
			}
			if (lastCheckDiff > 1) && (lastCheckDiff < 5) {
				log.Printf("AlertLoop : %08d %08d %02d", lastCheck, currentCheck, lastCheckDiff)
				currentCheck = lastCheck - 1
			}
			if currentCheck < lastCheck {
				r := currentCheck
				h := r / 3600
				s := r - (h * 3600)
				m := s / 60
				s -= m * 60
				if r > 0 {
					if r < 61 {
						if r%10 == 0 {
							// every 10 sec if T <= 1m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-seconds", s))
						}
					} else if r < 310 { // 5m = 300s
						if r%60 == 0 {
							// every min if T <= 5m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if r%300 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if r < 910 { // 15m = 900s
						if r%300 == 0 {
							// every 5 min if T <= 15m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if r%900 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if r < 1810 { // 30m = 1800s
						if r%600 == 0 {
							// every 10 min if T <= 30m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if r%1800 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if r < 10810 { // 3h = 10800s
						if r%3600 == 0 {
							// every 1 hour if T <= 3h
							TextToSpeechAlert(fmt.Sprintf("target-%02d-hours", h))
						}
					}
				}
				lastCheck = currentCheck
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func MainApp(appName string, appVersion string, enforceTarget bool) {
	mainAppName = appName
	mainAppVersion = appVersion

	e := os.Setenv("FYNE_THEME", "dark")
	if e != nil {
		log.Println("Setenv error?")
	}

	mainApp = app.NewWithID(mainAppName)

	/* Current Location */
	SetCurrentLocation(mainApp.Preferences().StringWithFallback("currentLocationName", "Use App Default"))

	/* Next Target */
	if !enforceTarget {
		ttm.SetNextTargetTime()
	}

	/* Time calculation */
	TimeLoop()

	/* Alerts */
	AlertLoop()

	MainWindowShow()
	mainApp.Run()
}
