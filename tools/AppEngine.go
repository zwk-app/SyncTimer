package tools

import (
	"SyncTimer/timer"
	"SyncTimer/tts"
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
	TextToSpeech struct {
		Object       *tts.TextToSpeech
		Embedded     bool
		EmbeddedPath string
		LocalPath    string `json:"audioPath"`
		GenerateTTS  bool   `json:"generate-tts"`
	}
	Timer struct {
		Object        *timer.TargetTimer
		TargetTime    string `json:"time"`
		TargetDelay   string `json:"delay"`
		LocationName  string
		EnforceTarget bool
	}
	Alerts struct {
		TextToSpeech  bool
		Notifications bool
	}
	Fyne struct {
		App fyne.App
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
	c.TextToSpeech.Object = nil
	c.TextToSpeech.Embedded = true
	c.TextToSpeech.EmbeddedPath = "audio/"
	c.TextToSpeech.LocalPath = path.Clean(c.Path + string(os.PathSeparator) + "audio" + string(os.PathSeparator))
	c.TextToSpeech.GenerateTTS = false
	c.Timer.Object = nil
	c.Timer.TargetTime = ""
	c.Timer.TargetDelay = ""
	c.Timer.LocationName = timer.LocalLocationName
	c.Timer.EnforceTarget = false
	c.Alerts.Notifications = false
	c.Alerts.TextToSpeech = true
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
	flag.StringVar(&c.TextToSpeech.LocalPath, "audioPath", c.TextToSpeech.LocalPath, "enforce audio local path")
	flag.BoolVar(&c.TextToSpeech.GenerateTTS, "generate-tts", c.TextToSpeech.GenerateTTS, "generate all TTS audio files")
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
	return fmt.Sprintf("%d.%d.%02d", c.version.major, c.version.minor, c.version.build)
}

func (c *AppEngine) Title() string {
	return fmt.Sprintf("%s v%s", c.Name(), c.Version())
}
