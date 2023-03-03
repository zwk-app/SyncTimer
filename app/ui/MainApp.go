package ui

import (
	app2 "SyncTimer/app"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
	"os"
	"time"
)

var appEngine *app2.AppEngine

func TextToSpeechAlert(name string) {
	if appEngine.Config.Alerts.TextToSpeech {
		log.Printf("TextToSpeechAlert '%s'", name)
		//goland:noinspection GoUnhandledErrorResult
		go appEngine.Audio.Play(name)
	}
}

func NotificationAlert(message string) {
	if appEngine.Config.Alerts.Notifications {
		log.Printf("NotificationAlert '%s'", message)
		go appEngine.FyneApp.SendNotification(fyne.NewNotification(appEngine.AppTitle(), message))
	}
}

func AlertLoop() {
	log.Println("AlertLoop : Start")
	time.Sleep(1500 * time.Millisecond)
	go func() {
		currentCheck := 0
		lastCheck := 0
		lastCheckDiff := 0
		for {
			currentCheck = appEngine.Timer.RemainingSeconds()
			lastCheckDiff = lastCheck - currentCheck
			if lastCheck < currentCheck {
				log.Printf("AlertLoop : %08d << %08d (%02d)", lastCheck, currentCheck, lastCheckDiff)
				lastCheck = currentCheck + 1
			}
			if (lastCheckDiff > 1) && (lastCheckDiff < 5) {
				log.Printf("AlertLoop : %08d <> %08d (%02d)", lastCheck, currentCheck, lastCheckDiff)
				currentCheck = lastCheck - 1
			}
			if currentCheck < lastCheck {
				h, m, s := appEngine.Timer.RemainingTime()
				if currentCheck >= 0 {
					if currentCheck < 11 {
						if currentCheck == 0 {
							TextToSpeechAlert(appEngine.Config.Alerts.AlarmSound)
						} else if currentCheck%2 == 0 {
							// every 2 sec if T <= 10s
							TextToSpeechAlert(fmt.Sprintf("target-%02d-seconds", s))
						}
					} else if currentCheck < 61 {
						if currentCheck%10 == 0 {
							// every 10 sec if T <= 1m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-seconds", s))
						}
					} else if currentCheck < 310 { // 5m = 300s
						if currentCheck%60 == 0 {
							// every min if T <= 5m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if currentCheck%300 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if currentCheck < 910 { // 15m = 900s
						if currentCheck%300 == 0 {
							// every 5 min if T <= 15m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if currentCheck%900 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if currentCheck < 1810 { // 30m = 1800s
						if currentCheck%600 == 0 {
							// every 10 min if T <= 30m
							TextToSpeechAlert(fmt.Sprintf("target-%02d-minutes", m))
						}
						if currentCheck%1800 == 0 {
							NotificationAlert(fmt.Sprintf("Target in %d minutes", m))
						}
					} else if currentCheck < 10810 { // 3h = 10800s
						if currentCheck%3600 == 0 {
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

func MainApp(a *app2.AppEngine) {
	appEngine = a
	e := os.Setenv("FYNE_THEME", "dark")
	if e != nil {
		log.Println("Setenv error?")
	}
	appEngine.FyneApp = app.NewWithID(appEngine.AppName())
	_ = appEngine.Config.LoadFyneSettings()
	AlertLoop()
	appEngine.FyneWindow = appEngine.FyneApp.NewWindow(appEngine.AppTitle())
	appEngine.FyneWindow.Resize(fyne.NewSize(320, 540))
	ShowMainWindow()
	appEngine.FyneApp.Run()
}
