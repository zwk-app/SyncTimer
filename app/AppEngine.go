package app

import (
	"SyncTimer/app/audio"
	"SyncTimer/app/timer"
	"embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"os"
	"path"
	"strings"
)

//goland:noinspection GoNameStartsWithPackageName
type AppEngine struct {
	name    string
	version struct {
		major int
		minor int
		build int
	}
	local struct {
		path string
	}
	embedded struct {
		fs   *embed.FS
		path string
	}
	Config     *Config
	Logs       *Logs
	Audio      *audio.AudioEngine
	Timer      *timer.TargetTime
	List       *timer.TargetList
	FyneApp    fyne.App
	FyneWindow fyne.Window
	lastError  error
}

func NewAppEngine(appName string, majorVersion int, minorVersion int, buildNumber int, embeddedFS *embed.FS, embeddedPath string) *AppEngine {
	r := &AppEngine{}
	r.name = appName
	r.version.major = majorVersion
	r.version.minor = minorVersion
	r.version.build = buildNumber
	r.SetLocalPath()
	r.embedded.fs = embeddedFS
	r.embedded.path = embeddedPath
	r.Config = NewConfig(r, path.Clean(r.local.path+string(os.PathSeparator)+r.name+".json"))
	r.Logs = NewAppLogs(r, LogsDebug)
	configJson, e := r.Config.ToJson()
	r.Logs.Debug("AppEngine", fmt.Sprintf("Config: \n%s", configJson), e)
	r.Audio = audio.NewAudioEngine(r.embedded.fs, r.embedded.path+"audio", r.local.path+"res/audio", "en")
	r.Timer = nil
	r.List = nil
	r.FyneApp = nil
	r.FyneWindow = nil
	r.lastError = nil
	return r
}

func (r *AppEngine) GetPath(linuxStylePath string) string {
	p := path.Clean(strings.ReplaceAll(linuxStylePath, "/", string(os.PathSeparator)))
	if !strings.HasSuffix(linuxStylePath, string(os.PathSeparator)) {
		p += string(os.PathSeparator)
	}
	return p
}

func (r *AppEngine) SetLocalPath() {
	ex, e := os.Executable()
	if e != nil {
		LogsCriticalErrorExit("AppEngine", "SetLocalPath: os.Executable() error", e)
	}
	r.local.path = r.GetPath(path.Dir(ex))
}

func (r *AppEngine) Error(title string, message string, e error) {
	r.lastError = e
	if r.FyneWindow != nil {
		dialog.NewInformation(title, message, r.FyneWindow).Show()
	}
	r.Logs.Error(title, message, e)
}

func (r *AppEngine) HasError() bool {
	return r.lastError != nil
}

func (r *AppEngine) GetError() error {
	if r.HasError() {
		e := r.lastError
		r.lastError = nil
		return e
	}
	return nil
}

func (r *AppEngine) AppName() string {
	return r.name
}

func (r *AppEngine) AppVersion() string {
	return fmt.Sprintf("%d.%d", r.version.major, r.version.minor)
}

func (r *AppEngine) AppBuild() string {
	return fmt.Sprintf("%d", r.version.build)
}

func (r *AppEngine) AppVersionWithBuild() string {
	return fmt.Sprintf("%s.%s", r.AppVersion(), r.AppBuild())

}

func (r *AppEngine) AppTitle() string {
	return fmt.Sprintf("%s v%s", r.AppName(), r.AppVersion())
}

func (r *AppEngine) AlarmSoundNames() []string {
	return []string{
		"navy-01-wanking-of-combat",
		"navy-02-breakfast",
		"navy-12-lunch-time",
		"navy-14-wake-up",
		"navy-22-wanking-of-combat",
	}
}

func (r *AppEngine) AlarmSoundTitles() []string {
	return []string{
		"Wanking of Combat (Morning)",
		"Breakfast Time",
		"Lunch Time",
		"Wake Up",
		"Wanking of Combat (Evening)",
	}
}

func (r *AppEngine) AlarmSoundName(alarmSoundTitle string) string {
	for i, v := range r.AlarmSoundTitles() {
		if v == alarmSoundTitle {
			return r.AlarmSoundNames()[i]
		}
	}
	return ""
}

func (r *AppEngine) ReadFile(fileName string) []byte {
	r.Logs.Debug("AppEngine", fmt.Sprintf("ReadFile: '%s'", fileName), nil)
	if r.HasError() {
		LogsCriticalErrorExit("AppEngine", "ReadFile: previous error not taken into account", nil)
	}
	/* First try : embeded file */
	if r.embedded.fs != nil {
		embeddedFileName := r.embedded.path + string(os.PathSeparator) + fileName
		embedBytes, embedError := r.embedded.fs.ReadFile(embeddedFileName)
		if embedError == nil {
			return embedBytes
		}
	}
	/* Else : local file */
	localFileName := path.Clean(r.local.path + string(os.PathSeparator) + fileName)
	localBytes, localError := os.ReadFile(localFileName)
	if localError == nil {
		return localBytes
	}
	r.Logs.Error("AppConfig", fmt.Sprintf("ReadFile: cannot read file '%s'", fileName), localError)
	return nil
}

func (r *AppEngine) Play(shortName string) *AppEngine {
	r.Logs.Debug("AppEngine", fmt.Sprintf("Play: '%s'", shortName), nil)
	if r.HasError() {
		LogsCriticalErrorExit("AppEngine", "Play: previous error not taken into account", nil)
	}
	r.Audio.Play(shortName)
	return r
}

func (r *AppEngine) SetTargetTime(targetString string) *AppEngine {
	r.Logs.Debug("AppEngine", fmt.Sprintf("SetTargetTime: '%s'", targetString), nil)
	r.Timer.SetTargetString(targetString)
	e := r.Timer.Error()
	if e != nil {
		r.lastError = e
		r.Logs.Error("AppEngine", fmt.Sprintf("SetTargetTime: '%s'", targetString), nil)
	} else {
		r.Timer.SetTextLabel(timer.DefaultTextLabel)
		r.Timer.SetAlertSound(r.Config.Alerts.AlarmSound)
	}
	return r
}

func (r *AppEngine) SetTargetJson(s string) *AppEngine {
	r.Logs.Debug("AppEngine", fmt.Sprintf("SetTargetJson: '%s'", s), nil)
	if r.List == nil {
		r.List = timer.NewTargetList()
	}
	if r.Config.Target.JsonName != s {
		r.Config.Target.JsonName = s
		r.List.LoadJson(r.Config.Target.JsonName)
		r.NextTarget()
	}
	return r
}

func (r *AppEngine) NextTarget() *AppEngine {
	r.Logs.Debug("AppEngine", fmt.Sprintf("NextTarget"), nil)
	if r.Timer == nil {
		r.Timer = timer.NewTargetTimer()
	}
	if r.Config.Target.Delay != "" {
		r.Timer.SetDelayString(r.Config.Target.Delay)
	} else if r.Config.Target.Time != "" {
		r.Timer.SetTargetString(r.Config.Target.Time)
	} else {
		next := r.List.NextTargetListItem()
		r.Timer.SetTargetTime(next.Time())
		r.Timer.SetTextLabel(next.TextLabel())
		if len(next.AlertSound()) > 0 {
			r.Timer.SetAlertSound(next.AlertSound())
		} else {
			r.Timer.SetAlertSound(r.Config.Alerts.AlarmSound)
		}
	}
	return r
}
