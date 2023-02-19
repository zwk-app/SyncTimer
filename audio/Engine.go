package audio

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type AudioEngine struct {
	player         *AudioPlayer
	textToSpeech   *TextToSpeech
	localAudioPath string
}

func NewAudioEngine(embeddedFS *embed.FS, embeddedAudioPath string, localAudioPath string, language string) *AudioEngine {
	t := AudioEngine{}
	t.player = NewAudioPlayer(embeddedFS, embeddedAudioPath)
	t.textToSpeech = NewTextToSpeech(language)
	t.localAudioPath = localAudioPath
	return &t
}

func (t *AudioEngine) GetFileName(shortName string) string {
	if strings.HasPrefix(shortName, "target") || strings.HasPrefix(shortName, "about") {
		return path.Clean(t.localAudioPath + string(os.PathSeparator) + shortName + "." + t.textToSpeech.language + ".mp3")
	}
	return path.Clean(t.localAudioPath + string(os.PathSeparator) + shortName + ".mp3")
}

func (t *AudioEngine) Create(shortName string, message string) error {
	log.Printf("AudioEngine->Create '%s' for '%s'", shortName, message)
	e := t.textToSpeech.CreateFile(t.GetFileName(shortName), message)
	if e != nil {
		log.Printf("AudioEngine->Create error: '%s'", e.Error())
		return e
	}
	return nil
}

func (t *AudioEngine) Play(shortName string) error {
	log.Printf("AudioEngine->Play '%s'", shortName)
	e := t.player.Play(t.GetFileName(shortName))
	if e != nil {
		log.Printf("AudioEngine->Play error: '%s'", e.Error())
		return e
	}
	return nil
}

func (t *AudioEngine) generateLanguageFiles(appName string, language string) {
	log.Printf("Audio->GenerateAudioFiles for '%s' (%s)", appName, language)
	name := "about"
	mesg := fmt.Sprintf("Hello dear! My name is %s, I am so pleased to meet you!", appName)
	_ = t.Create(name, mesg)
	for i := 1; i <= 3; i++ {
		name = fmt.Sprintf("target-%02d-hours", i)
		mesg = fmt.Sprintf("%d hours left", i)
		_ = t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 1; i < 10; i += 1 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("%d minutes left", i)
		_ = t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 10; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("%d minutes left", i)
		_ = t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 1; i < 10; i += 1 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("%d", i)
		_ = t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 10; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("%d seconds left", i)
		_ = t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
}

func (t *AudioEngine) GenerateAllAudioFiles(appName string) {
	log.Printf("GenerateAllAudioFiles for '%s' in '%s'", appName, t.localAudioPath)
	t.generateLanguageFiles(appName, "en")
}
