package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/zwk-app/zwk-tools/tools"
	"os"
	"path"
	"strings"
)

const LocalLocationName = "Local Time"

//goland:noinspection ALL
const UtcLocationName = "UTC"

const DefaultTextLabel = "Target"
const DefaultAlarmSound = "navy-14-wake-up"

type AppConfig struct {
	name    string
	version struct {
		major int
		minor int
		build int
	}
	Logs     LogsConfig  `json:"logs"`
	Proxy    ProxyConfig `json:"proxy"`
	Audio    AudioConfig `json:"audio"`
	Location LocationConfig
	Target   TargetConfig `json:"target"`
	Alerts   AlertsConfig `json:"alerts"`
}

type LogsConfig struct {
	StdOut   bool   `json:"stdout"`
	FileName string `json:"filename"`
	Verbose  bool   `json:"verbose"`
}

type ProxyConfig struct {
	Server string `json:"server"`
	Port   int    `json:"port"`
}

type AudioConfig struct {
	LocalPath string `json:"path"`
	Make      bool   `json:"make"`
}

type LocationConfig struct {
	Name string
}

type TargetConfig struct {
	JsonName string `json:"json"`
	Time     string `json:"time"`
	Delay    string `json:"delay"`
}

type AlertsConfig struct {
	TextToSpeech  bool   `json:"tts"`
	Notifications bool   `json:"notif"`
	AlarmSound    string `json:"alarm"`
}

var config *AppConfig = nil

func Config() *AppConfig {
	if config == nil {
		config = new(AppConfig)
		config.Location.Name = LocalLocationName
		config.Alerts.TextToSpeech = true
		config.Alerts.Notifications = false
		config.Alerts.AlarmSound = DefaultAlarmSound
	}
	return config
}

func SetAppInfo(appName string, majorVersion int, minorVersion int, buildNumber int) {
	Config().name = appName
	Config().version.major = majorVersion
	Config().version.minor = minorVersion
	Config().version.build = buildNumber
}
func Name() string {
	return Config().name
}
func Version() string {
	return fmt.Sprintf("%d.%d", Config().version.major, Config().version.minor)
}
func Build() string {
	return fmt.Sprintf("%d", Config().version.build)
}

//goland:noinspection GoUnusedExportedFunction
func VersionWithBuild() string {
	return fmt.Sprintf("%s.%s", Version(), Build())

}
func Title() string {
	return fmt.Sprintf("%s v%s", Name(), Version())
}

func Logs() *LogsConfig {
	return &Config().Logs
}
func Proxy() *ProxyConfig {
	return &Config().Proxy
}
func Audio() *AudioConfig {
	return &Config().Audio
}
func Location() *LocationConfig {
	return &Config().Location
}
func Target() *TargetConfig {
	return &Config().Target
}
func Alerts() *AlertsConfig {
	return &Config().Alerts
}

func LoadEnvironment() {
	proxyUrl := os.Getenv("HTTP_PROXY")
	if len(proxyUrl) > 0 {
		proxyPat := `http[s]{0,1}://(?P<server>[a-zA-Z0-9\._-]+):(?P<port>[0-9]+)`
		Config().Proxy.Server = tools.FirstMatch(proxyUrl, proxyPat, "server")
		Config().Proxy.Port = tools.Int(tools.FirstMatch(proxyUrl, proxyPat, "port"))
	}
}

func ToJson() ([]byte, error) {
	data, e := json.Marshal(Config())
	if e == nil {
		return data, e
	}
	return nil, e
}

func FromJson(data []byte) error {
	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	e := jsonDecoder.Decode(Config())
	if e != nil {
		return e
	}
	return nil
}

func DefaultConfig() string {
	ex, e := os.Executable()
	if e == nil {
		return strings.Replace(ex, path.Base(ex), tools.AlphaNums(Name())+".json", -1)
	}
	return ""
}

func LoadFile(jsonFileName string) error {
	if len(jsonFileName) == 0 {
		jsonFileName = DefaultConfig()
	}
	jsonBytes, e := os.ReadFile(jsonFileName)
	if e != nil {
		return e
	}
	return FromJson(jsonBytes)
}

func LoadArguments() {
	r := Config()
	flag.BoolVar(&r.Logs.StdOut, "stdout", r.Logs.StdOut, "Display app.logs in Stdout")
	flag.StringVar(&r.Logs.FileName, "log", r.Logs.FileName, "Save app.logs in file")
	flag.BoolVar(&r.Logs.Verbose, "verbose", r.Logs.Verbose, "Verbose mode")
	flag.StringVar(&r.Audio.LocalPath, "audio-path", r.Audio.LocalPath, "enforce audio local path")
	flag.BoolVar(&r.Audio.Make, "audio-make", r.Audio.Make, "generate all TTS audio files")
	flag.StringVar(&r.Target.JsonName, "targets-json", r.Target.JsonName, "set targets Json list URL or filename")
	flag.StringVar(&r.Target.Time, "time", r.Target.Time, "set target time to <hh[mm[ss]]>")
	flag.StringVar(&r.Target.Delay, "delay", r.Target.Delay, "set target delay in <[[hh]mm]ss>")
	flag.Parse()
}

func LoadFyneSettings(fyneApp fyne.App) {
	r := Config()
	r.Location.Name = fyneApp.Preferences().StringWithFallback("currentLocationName", r.Location.Name)
	r.Alerts.TextToSpeech = fyneApp.Preferences().BoolWithFallback("voiceAlertsEnabled", r.Alerts.TextToSpeech)
	r.Alerts.Notifications = fyneApp.Preferences().BoolWithFallback("notificationsEnabled", r.Alerts.Notifications)
	r.Alerts.AlarmSound = fyneApp.Preferences().StringWithFallback("alarmSound", r.Alerts.AlarmSound)
	r.Target.JsonName = fyneApp.Preferences().StringWithFallback("targetsJson", r.Target.JsonName)
}

func SaveFyneSettings(fyneApp fyne.App) {
	r := Config()
	fyneApp.Preferences().SetString("currentLocationName", r.Location.Name)
	fyneApp.Preferences().SetBool("voiceAlertsEnabled", r.Alerts.TextToSpeech)
	fyneApp.Preferences().SetBool("notificationsEnabled", r.Alerts.Notifications)
	fyneApp.Preferences().SetString("alarmSound", r.Alerts.AlarmSound)
	fyneApp.Preferences().SetString("targetsJson", r.Target.JsonName)
}
