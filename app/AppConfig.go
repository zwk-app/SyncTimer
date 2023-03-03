package app

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	app      *AppEngine
	FileName string
	Logs     struct {
		StdOut   bool   `json:"stdout"`
		FileName string `json:"filename"`
	} `json:"logs"`
	Proxy struct {
		String string `json:"url"`
		Server string `json:"server"`
		Port   int    `json:"port"`
	} `json:"proxy"`
	Audio struct {
		LocalPath string `json:"path"`
		Make      bool   `json:"make"`
	} `json:"audio"`
	Location struct {
		Name string
	}
	Target struct {
		JsonName string `json:"json"`
		Time     string `json:"time"`
		Delay    string `json:"delay"`
	} `json:"target"`
	Alerts struct {
		TextToSpeech  bool   `json:"tts"`
		Notifications bool   `json:"notif"`
		AlarmSound    string `json:"alarm"`
	} `json:"alerts"`
}

func NewConfig(appEngine *AppEngine, configFileName string) *Config {
	r := &Config{}
	r.app = appEngine
	if r.app == nil {
		LogsCriticalErrorExit("AppConfig", "NewAppConfig: AppEngine is not set", nil)
	}
	r.LoadEnvironment().LoadConfigFile(configFileName).LoadCommandLineArguments()
	return r
}

func (r *Config) LoadEnvironment() *Config {
	r.Proxy.String = os.Getenv("HTTP_PROXY")
	return r
}

func (r *Config) ToJson() ([]byte, error) {
	data, e := json.Marshal(r)
	if e == nil {
		return data, e
	}
	return nil, e
}

func (r *Config) FromJson(data []byte) error {
	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	e := jsonDecoder.Decode(r)
	if e != nil {
		return e
	}
	return nil
}

func (r *Config) LoadConfigFile(jsonFileName string) *Config {
	if len(jsonFileName) > 0 {
		if jsonBytes, e := os.ReadFile(jsonFileName); e == nil {
			if e = r.FromJson(jsonBytes); e == nil {
				r.FileName = jsonFileName
			}
		}
	}
	return r
}

func (r *Config) LoadCommandLineArguments() *Config {
	flag.BoolVar(&r.Logs.StdOut, "stdout", r.Logs.StdOut, "Display app.logs in Stdout")
	flag.StringVar(&r.Logs.FileName, "log", r.Logs.FileName, "Save app.logs in file")
	flag.StringVar(&r.Audio.LocalPath, "app.audio-path", r.Audio.LocalPath, "enforce app.audio local path")
	flag.BoolVar(&r.Audio.Make, "app.audio-make", r.Audio.Make, "generate all TTS app.audio files")
	flag.StringVar(&r.Target.JsonName, "targets-json", r.Target.JsonName, "set targets list Json URL or filename")
	flag.StringVar(&r.Target.Time, "target-time", r.Target.Time, "set target time to <hh[mm[ss]]>")
	flag.StringVar(&r.Target.Delay, "target-delay", r.Target.Delay, "set target delay in <[[hh]mm]ss>")
	flag.Parse()
	return r
}

func (r *Config) LoadFyneSettings() *Config {
	if r.app == nil {
		LogsCriticalErrorExit("AppConfig", "LoadFyneSettings: AppEngine is not set", nil)
	}
	if r.app.FyneApp == nil {
		LogsCriticalErrorExit("AppConfig", "SaveFyneSettings: FyneApp is not set", nil)
	}
	r.Location.Name = r.app.FyneApp.Preferences().StringWithFallback("currentLocationName", r.Location.Name)
	r.Alerts.TextToSpeech = r.app.FyneApp.Preferences().BoolWithFallback("voiceAlertsEnabled", r.Alerts.TextToSpeech)
	r.Alerts.Notifications = r.app.FyneApp.Preferences().BoolWithFallback("notificationsEnabled", r.Alerts.Notifications)
	r.Alerts.AlarmSound = r.app.FyneApp.Preferences().StringWithFallback("alarmSound", r.Alerts.AlarmSound)
	r.Target.JsonName = r.app.FyneApp.Preferences().StringWithFallback("targetsJson", r.Target.JsonName)
	return r
}

func (r *Config) SaveFyneSettings() *Config {
	if r.app == nil {
		LogsCriticalErrorExit("AppConfig", "SaveFyneSettings: AppEngine is not set", nil)
	}
	if r.app.FyneApp == nil {
		LogsCriticalErrorExit("AppConfig", "SaveFyneSettings: FyneApp is not set", nil)
	}
	r.app.FyneApp.Preferences().SetString("currentLocationName", r.Location.Name)
	r.app.FyneApp.Preferences().SetBool("voiceAlertsEnabled", r.Alerts.TextToSpeech)
	r.app.FyneApp.Preferences().SetBool("notificationsEnabled", r.Alerts.Notifications)
	r.app.FyneApp.Preferences().SetString("alarmSound", r.Alerts.AlarmSound)
	r.app.FyneApp.Preferences().SetString("targetsJson", r.Target.JsonName)
	return r
}
