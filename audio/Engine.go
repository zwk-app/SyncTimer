package audio

import (
	"SyncTimer/resources"
	"fmt"
	"github.com/zwk-app/go-tools/logs"
	"path"
	"strings"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type AudioEngine struct {
	local struct {
		path string
	}
	last struct {
		error error
	}
}

var engine *AudioEngine

func Engine() *AudioEngine {
	if engine == nil {
		engine = &AudioEngine{}
		engine.last.error = nil
	}
	return engine
}

func SetError(e error) {
	Engine().last.error = e
}

func HasError() bool {
	return Engine().last.error != nil
}

func GetError() error {
	lastError := Engine().last.error
	SetError(nil)
	return lastError
}

func FileName(shortName string) string {
	if strings.HasPrefix(shortName, "target") || strings.HasPrefix(shortName, "about") {
		return path.Clean(Engine().local.path + shortName + "." + TextToSpeech().language + ".mp3")
	}
	return path.Clean(Engine().local.path + shortName + ".mp3")
}

func Create(shortName string, message string) {
	logs.Debug("AudioEngine", fmt.Sprintf("Create '%s' for '%s'", shortName, message), nil)
	CreateFile(FileName(shortName), message)
	if HasError() {
		logs.Error("AudioEngine", "", Engine().last.error)
	}
}

func PlayEmbedded(shortName string) {
	logs.Debug("AudioEngine", fmt.Sprintf("PlayEmbedded '%s'", shortName), nil)
	PlayMp3Content(resources.ReadAudio(FileName(shortName)))
	if HasError() {
		logs.Error("AudioEngine", "", Engine().last.error)
	}
}

func PlayLocal(shortName string) {
	logs.Debug("AudioEngine", fmt.Sprintf("PlayLocal '%s'", shortName), nil)
	PlayMp3File(FileName(shortName))
	if HasError() {
		logs.Error("AudioEngine", "", Engine().last.error)
	}
}

func Play(shortName string) {
	logs.Debug("AudioEngine", fmt.Sprintf("Play '%s'", shortName), nil)
	PlayEmbedded(shortName)
	if HasError() {
		SetError(nil)
		PlayLocal(shortName)
	}
}

func GenerateLang(appName string, language string) {
	logs.Debug("AudioEngine", fmt.Sprintf("GenerateLang '%s' for '%s", language, appName), nil)
	name := "about"
	mesg := fmt.Sprintf("Hello dear! My name is %s, I am so pleased to meet you!", appName)
	Create(name, mesg)
	for i := 1; i <= 3; i++ {
		name = fmt.Sprintf("target-%02d-hours", i)
		mesg = fmt.Sprintf("%d hours left", i)
		Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 1; i < 10; i += 1 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("%d minutes left", i)
		Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 10; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("%d minutes left", i)
		Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 1; i < 10; i += 1 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("%d", i)
		Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 10; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("%d seconds left", i)
		Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
}

func GenerateAll(appName string) {
	logs.Debug("AudioEngine", fmt.Sprintf("GenerateAll"), nil)
	GenerateLang(appName, "en")
}
