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
	appName    string
	appVersion struct {
		major int
		minor int
		build int
	}
	appPath        string
	ConfigFileName string
	EmbeddedFS     *embed.FS
	Logs           struct {
		StdOut   bool   `json:"stdout"`
		FileName string `json:"log"`
	}
	Audio struct {
		Engine       *audio.AudioEngine
		Embedded     bool
		EmbeddedPath string
		LocalPath    string `json:"audioPath"`
		GenerateTTS  bool   `json:"generate-audio"`
	}
	Timer struct {
		Engine       *timer.TargetTime
		List         *timer.TargetList
		TargetsJson  string `json:"targets-json"`
		TargetTime   string `json:"target-time"`
		TargetDelay  string `json:"target-delay"`
		LocationName string
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
	last        struct {
		error error
	}
}

func NewAppEngine(appName string, major int, minor int, build int, appEmbeddedFS *embed.FS) *AppEngine {
	log.Printf("--- %s appVersion %d.%d build %02d ---", appName, major, minor, build)
	c := AppEngine{}
	c.appName = appName
	c.appVersion.major = major
	c.appVersion.minor = minor
	c.appVersion.build = build
	ex, e := os.Executable()
	if e != nil {
		ErrorExit(e.Error())
	}
	c.appPath = path.Dir(ex)
	c.ConfigFileName = path.Clean(c.appPath + string(os.PathSeparator) + c.appName + ".json")
	c.EmbeddedFS = appEmbeddedFS
	c.Logs.StdOut = false
	c.Logs.FileName = ""
	c.Audio.Engine = nil
	c.Audio.Embedded = true
	c.Audio.EmbeddedPath = "res/audio/"
	c.Audio.LocalPath = path.Clean(c.appPath + string(os.PathSeparator) + "res" + string(os.PathSeparator) + "audio" + string(os.PathSeparator))
	c.Audio.GenerateTTS = false
	c.Timer.Engine = nil
	c.Timer.List = nil
	c.Timer.TargetTime = ""
	c.Timer.TargetDelay = ""
	c.Timer.LocationName = timer.LocalLocationName
	c.Alerts.Notifications = false
	c.Alerts.TextToSpeech = true
	c.Alerts.AlertSound = "navy-12-lunch-time"
	c.Alerts.AlertSoundNames = append(c.Alerts.AlertSoundNames, "navy-01-wanking-of-combat", "navy-02-breakfast", "navy-12-lunch-time", "navy-14-wake-up", "navy-22-wanking-of-combat")
	c.Alerts.AlertSoundTitles = append(c.Alerts.AlertSoundTitles, "Wanking of Combat (Morning)", "Breakfast Time", "Lunch Time", "Wake Up", "Wanking of Combat (Evening)")
	c.Fyne.App = nil
	c.ProxyString = ""
	return &c
}

func (c *AppEngine) Error() error {
	lastError := c.last.error
	c.last.error = nil
	return lastError
}

func (c *AppEngine) AppName() string {
	return c.appName
}

func (c *AppEngine) AppVersion() string {
	return fmt.Sprintf("%d.%d.%d", c.appVersion.major, c.appVersion.minor, c.appVersion.build)
}

func (c *AppEngine) AppTitle() string {
	return fmt.Sprintf("%s v%s", c.AppName(), c.AppVersion())
}

func (c *AppEngine) AppPath() string {
	return c.appPath
}

func (c *AppEngine) AlertName(alertTitle string) string {
	for i, t := range c.Alerts.AlertSoundTitles {
		if t == alertTitle {
			return c.Alerts.AlertSoundNames[i]
		}
	}
	return ""
}

func (c *AppEngine) LoadEnvSettings() *AppEngine {
	log.Printf("AppEngine->LoadEnvSettings")
	c.ProxyString = os.Getenv("HTTP_PROXY")
	return c
}

func (c *AppEngine) jsonLoad(jsonFileName string) error {
	jsonFile, e := os.Open(jsonFileName)
	if e != nil {
		return e
	}
	jsonDecoder := json.NewDecoder(jsonFile)
	e = jsonDecoder.Decode(c)
	if e != nil {
		return e
	}
	log.Printf("loadJson: Success")
	return nil
}

func (c *AppEngine) LoadFileSettings(jsonFileName string) *AppEngine {
	if jsonFileName != "" {
		c.ConfigFileName = jsonFileName
	}
	if c.jsonLoad(c.ConfigFileName) == nil {
		log.Printf("AppEngine->LoadFileSettings: '%s'", c.ConfigFileName)
	}
	return c
}

func (c *AppEngine) LoadArgSettings() *AppEngine {
	log.Printf("AppEngine->LoadArgSettings")
	flag.BoolVar(&c.Logs.StdOut, "stdout", c.Logs.StdOut, "Display logs in Stdout")
	flag.StringVar(&c.Logs.FileName, "log", c.Logs.FileName, "Save logs in file")
	flag.StringVar(&c.Audio.LocalPath, "audioPath", c.Audio.LocalPath, "enforce audio local path")
	flag.BoolVar(&c.Audio.GenerateTTS, "generate-audio", c.Audio.GenerateTTS, "generate all TTS audio files")
	flag.StringVar(&c.Timer.TargetsJson, "targets-json", c.Timer.TargetsJson, "set targets list Json URL or filename")
	flag.StringVar(&c.Timer.TargetTime, "target-time", c.Timer.TargetTime, "set target time to <hh[mm[ss]]>")
	flag.StringVar(&c.Timer.TargetDelay, "target-delay", c.Timer.TargetDelay, "set target delay in <[[hh]mm]ss>")
	flag.Parse()
	return c
}

func (c *AppEngine) LoadFyneSettings() *AppEngine {
	log.Printf("AppEngine->LoadFyneSettings")
	if c.Fyne.App == nil {
		log.Printf("AppEngine->LoadFyneSettings error: Fyne.App not set")
		c.last.error = fmt.Errorf("Fyne.App not set")
		return c
	}
	c.Timer.LocationName = c.Fyne.App.Preferences().StringWithFallback("currentLocationName", c.Timer.LocationName)
	c.Alerts.TextToSpeech = c.Fyne.App.Preferences().BoolWithFallback("voiceAlertsEnabled", c.Alerts.TextToSpeech)
	c.Alerts.Notifications = c.Fyne.App.Preferences().BoolWithFallback("notificationsEnabled", c.Alerts.Notifications)
	c.Alerts.AlertSound = c.Fyne.App.Preferences().StringWithFallback("alertSound", c.Alerts.AlertSound)
	c.SetTargetJson(c.Fyne.App.Preferences().StringWithFallback("targetsJson", c.Timer.TargetsJson))
	return c
}

func (c *AppEngine) SaveFyneSettings() *AppEngine {
	log.Printf("AppEngine->SaveFyneSettings")
	if c.Fyne.App == nil {
		log.Printf("AppEngine->SaveFyneSettings error: Fyne.App not set")
		c.last.error = fmt.Errorf("Fyne.App not set")
		return c
	}
	c.Fyne.App.Preferences().SetString("currentLocationName", c.Timer.LocationName)
	c.Fyne.App.Preferences().SetBool("voiceAlertsEnabled", c.Alerts.TextToSpeech)
	c.Fyne.App.Preferences().SetBool("notificationsEnabled", c.Alerts.Notifications)
	c.Fyne.App.Preferences().SetString("alertSound", c.Alerts.AlertSound)
	c.Fyne.App.Preferences().SetString("targetsJson", c.Timer.TargetsJson)
	c.SetTargetJson(c.Timer.TargetsJson)
	return c
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

func (c *AppEngine) Play(shortName string) *AppEngine {
	log.Printf("AppEngine->Play '%s'", shortName)
	if c.Audio.Engine == nil {
		c.Audio.Engine = audio.NewAudioEngine(c.EmbeddedFS, c.Audio.EmbeddedPath, c.Audio.LocalPath, "en")
	}
	c.Audio.Engine.Play(shortName)
	return c
}

func (c *AppEngine) SetTargetTime(targetString string) *AppEngine {
	log.Printf("AppEngine->SetTargetTime '%s'", targetString)
	c.Timer.Engine.SetTargetString(targetString)
	e := c.Timer.Engine.Error()
	if e != nil {
		c.last.error = e
		log.Printf("AppEngine->SetTargetTime error: %s", e.Error())
	} else {
		c.Timer.Engine.SetTextLabel(timer.DefaultTextLabel)
		c.Timer.Engine.SetAlertSound(c.Alerts.AlertSound)
	}
	return c
}

func (c *AppEngine) SetTargetJson(s string) *AppEngine {
	log.Printf("AppEngine->SetTargetJson '%s'", s)
	if c.Timer.List == nil {
		c.Timer.List = timer.NewTargetList()
	}
	if c.Timer.TargetsJson != s {
		c.Timer.TargetsJson = s
		c.Timer.List.LoadJson(c.Timer.TargetsJson)
		c.NextTarget()
	}
	return c
}

func (c *AppEngine) NextTarget() *AppEngine {
	log.Printf("AppEngine->NextTarget")
	if c.Timer.Engine == nil {
		c.Timer.Engine = timer.NewTargetTimer()
	}
	if c.Timer.TargetDelay != "" {
		c.Timer.Engine.SetDelayString(c.Timer.TargetDelay)
	} else if c.Timer.TargetTime != "" {
		c.Timer.Engine.SetTargetString(c.Timer.TargetTime)
	} else {
		next := c.Timer.List.NextTargetListItem()
		c.Timer.Engine.SetTargetTime(next.Time())
		c.Timer.Engine.SetTextLabel(next.TextLabel())
		if len(next.AlertSound()) > 0 {
			c.Timer.Engine.SetAlertSound(next.AlertSound())
		} else {
			c.Timer.Engine.SetAlertSound(c.Alerts.AlertSound)
		}
	}
	return c
}
