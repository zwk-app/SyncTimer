package tools

import (
	"SyncTimer/audio"
	"SyncTimer/timer"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"log"
	"os"
	"path"
)

type AppEngine struct {
	name    string
	version struct {
		major int
		minor int
		build int
	}
	Path           string
	ConfigFileName string
	EmbeddedFS     *embed.FS
	Logs           struct {
		StdOut   bool   `json:"stdout"`
		FileName string `json:"log"`
	}
	Audio struct {
		Object       *audio.AudioEngine
		Embedded     bool
		EmbeddedPath string
		LocalPath    string `json:"audioPath"`
		GenerateTTS  bool   `json:"generate-audio"`
	}
	Timer struct {
		Object        *timer.TargetTimer
		TargetTime    string `json:"time"`
		TargetDelay   string `json:"delay"`
		LocationName  string
		EnforceTarget bool
	}
	Alerts struct {
		TextToSpeech     bool
		Notifications    bool
		AlertSound       string
		AlertSoundNames  []string
		AlertSoundTitles []string
	}
	Fyne struct {
		App        fyne.App
		MainWindow fyne.Window
	}
	ProxyString string `json:"proxy"`
}

func NewAppEngine(appName string, major int, minor int, build int, appEmbeddedFS *embed.FS) *AppEngine {
	log.Printf("--- %s version %d.%d build %02d ---", appName, major, minor, build)
	c := AppEngine{}
	c.name = appName
	c.version.major = major
	c.version.minor = minor
	c.version.build = build
	ex, e := os.Executable()
	if e != nil {
		ErrorExit(e.Error())
	}
	c.Path = path.Dir(ex)
	c.ConfigFileName = path.Clean(c.Path + string(os.PathSeparator) + c.name + ".json")
	c.EmbeddedFS = appEmbeddedFS
	c.Logs.StdOut = false
	c.Logs.FileName = ""
	c.Audio.Object = nil
	c.Audio.Embedded = true
	c.Audio.EmbeddedPath = "res/audio/"
	c.Audio.LocalPath = path.Clean(c.Path + string(os.PathSeparator) + "res" + string(os.PathSeparator) + "audio" + string(os.PathSeparator))
	c.Audio.GenerateTTS = false
	c.Timer.Object = nil
	c.Timer.TargetTime = ""
	c.Timer.TargetDelay = ""
	c.Timer.LocationName = timer.LocalLocationName
	c.Timer.EnforceTarget = false
	c.Alerts.Notifications = false
	c.Alerts.TextToSpeech = true
	c.Alerts.AlertSound = "navy-12-lunch-time"
	c.Alerts.AlertSoundNames = append(c.Alerts.AlertSoundNames, "navy-01-wanking-of-combat", "navy-02-breakfast", "navy-12-lunch-time", "navy-14-wake-up", "navy-22-wanking-of-combat")
	c.Alerts.AlertSoundTitles = append(c.Alerts.AlertSoundTitles, "Wanking of Combat (Morning)", "Breakfast Time", "Lunch Time", "Wake Up", "Wanking of Combat (Evening)")
	c.Fyne.App = nil
	c.ProxyString = ""
	return &c
}

func (c *AppEngine) LoadEnvSettings() *AppEngine {
	log.Printf("AppEngine->LoadEnvSettings")
	c.ProxyString = os.Getenv("HTTP_PROXY")
	return c
}

func (c *AppEngine) jsonLoad(jsonFileName string) error {
	log.Printf("AppEngine->jsonLoad: '%s'", jsonFileName)
	jsonFile, e := os.Open(jsonFileName)
	if e != nil {
		return e
	}
	jsonDecoder := json.NewDecoder(jsonFile)
	e = jsonDecoder.Decode(c)
	if e != nil {
		return e
	}
	log.Printf("JsonLoad: Success")
	return nil
}

func (c *AppEngine) LoadFileSettings(jsonFileName string) *AppEngine {
	if len(jsonFileName) > 0 {
		c.ConfigFileName = jsonFileName
	}
	log.Printf("AppEngine->LoadFileSettings: '%s'", c.ConfigFileName)
	if c.jsonLoad(c.ConfigFileName) == nil {
		log.Printf("AppEngine->LoadFileSettings: Success")
	}
	return c
}

func (c *AppEngine) LoadArgSettings() *AppEngine {
	log.Printf("AppEngine->LoadArgSettings")
	flag.BoolVar(&c.Logs.StdOut, "stdout", c.Logs.StdOut, "Display logs in Stdout")
	flag.StringVar(&c.Logs.FileName, "log", c.Logs.FileName, "Save logs in file")
	flag.StringVar(&c.Audio.LocalPath, "audioPath", c.Audio.LocalPath, "enforce audio local path")
	flag.BoolVar(&c.Audio.GenerateTTS, "generate-audio", c.Audio.GenerateTTS, "generate all TTS audio files")
	flag.StringVar(&c.Timer.TargetTime, "time", c.Timer.TargetTime, "set target time to <hh[mm[ss]]>")
	flag.StringVar(&c.Timer.TargetDelay, "delay", c.Timer.TargetDelay, "set target delay in <[[hh]mm]ss>")
	flag.Parse()
	return c
}

func (c *AppEngine) LoadFyneSettings() error {
	log.Printf("AppEngine->LoadFyneSettings")
	if c.Fyne.App == nil {
		log.Printf("AppEngine->LoadFyneSettings error: Fyne.App not set")
		return fmt.Errorf("Fyne.App not set")
	}
	c.Timer.LocationName = c.Fyne.App.Preferences().StringWithFallback("currentLocationName", c.Timer.LocationName)
	c.Alerts.TextToSpeech = c.Fyne.App.Preferences().BoolWithFallback("voiceAlertsEnabled", c.Alerts.TextToSpeech)
	c.Alerts.Notifications = c.Fyne.App.Preferences().BoolWithFallback("notificationsEnabled", c.Alerts.Notifications)
	c.Alerts.AlertSound = c.Fyne.App.Preferences().StringWithFallback("alertSound", c.Alerts.AlertSound)
	return nil
}

func (c *AppEngine) SaveFyneSettings() error {
	log.Printf("AppEngine->SaveFyneSettings")
	if c.Fyne.App == nil {
		log.Printf("AppEngine->SaveFyneSettings error: Fyne.App not set")
		return fmt.Errorf("Fyne.App not set")
	}
	c.Fyne.App.Preferences().SetString("currentLocationName", c.Timer.LocationName)
	c.Fyne.App.Preferences().SetBool("voiceAlertsEnabled", c.Alerts.TextToSpeech)
	c.Fyne.App.Preferences().SetBool("notificationsEnabled", c.Alerts.Notifications)
	c.Fyne.App.Preferences().SetString("alertSound", c.Alerts.AlertSound)
	return nil
}

func (c *AppEngine) SetLogOptions() *AppEngine {
	log.Printf("AppEngine->SetLogOptions")
	if len(c.Logs.FileName) > 0 {
		logFile, e := os.OpenFile(c.Logs.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			ErrorExit(e.Error())
		}
		log.SetOutput(logFile)
	}
	if c.Logs.StdOut {
		log.SetOutput(os.Stdout)
	}
	return c
}

func (c *AppEngine) Name() string {
	return c.name
}

func (c *AppEngine) Version() string {
	return fmt.Sprintf("%d.%d.%d", c.version.major, c.version.minor, c.version.build)
}

func (c *AppEngine) Title() string {
	return fmt.Sprintf("%s v%s", c.Name(), c.Version())
}

func (c *AppEngine) AlertName(alertTitle string) string {
	for i, t := range c.Alerts.AlertSoundTitles {
		if t == alertTitle {
			return c.Alerts.AlertSoundNames[i]
		}
	}
	return ""
}
