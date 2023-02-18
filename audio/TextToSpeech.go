package audio

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const ttsDefaultLanguage = "en"
const ttsChannelCount = 2
const ttsBitDepthInBytes = 2

type TextToSpeech struct {
	name         string
	path         string
	lang         string
	embeddedFS   *embed.FS
	embeddedPath string
}

func NewTextToSpeech(mainAppName string, audioPath string, language string) *TextToSpeech {
	if !strings.HasSuffix(audioPath, string(os.PathSeparator)) {
		audioPath += string(os.PathSeparator)
	}
	switch language {
	case "en":
		return &TextToSpeech{name: mainAppName, path: audioPath, lang: "en"}
	}
	return &TextToSpeech{name: mainAppName, path: audioPath, lang: ttsDefaultLanguage}
}

func (t *TextToSpeech) SetEmbeddedAudioFS(embeddedAudioFS *embed.FS, embeddedAudioPath string) {
	t.embeddedFS = embeddedAudioFS
	t.embeddedPath = embeddedAudioPath
}

func (t *TextToSpeech) GetFileName(name string) string {
	if strings.HasPrefix(name, "target") || strings.HasPrefix(name, "about") {
		return t.path + name + "." + t.lang + ".mp3"
	}
	return t.path + name + ".mp3"
}

func (t *TextToSpeech) PlayFile(fileName string) error {
	embeddedFileName := t.embeddedPath + filepath.Base(fileName)
	fileBytes, e := t.embeddedFS.ReadFile(embeddedFileName)
	if e == nil {
		log.Printf("Audio->PlayFile embedded:%s", embeddedFileName)
	} else {
		fileBytes, e = os.ReadFile(fileName)
		if e != nil {
			log.Printf("Audio->PlayFile ERROR:%s", e.Error())
			return e
		}
		log.Printf("Audio->PlayFile local:%s", fileName)
	}
	fileDecoder, e := mp3.NewDecoder(bytes.NewReader(fileBytes))
	if e != nil {
		return e
	}
	otoContext, readyChan, e := oto.NewContext(fileDecoder.SampleRate(), ttsBitDepthInBytes, ttsChannelCount)
	if e != nil {
		return e
	}
	<-readyChan
	filePlayer := otoContext.NewPlayer(fileDecoder)
	filePlayer.Play()
	for filePlayer.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}
	return filePlayer.Close()
}

func (t *TextToSpeech) Play(name string) {
	e := t.PlayFile(t.GetFileName(name))
	if e != nil {
		log.Printf("Audio->Play ERROR: '%s'", e.Error())
	}
}

func (t *TextToSpeech) CreateFile(filename string, message string) error {
	log.Printf("Audio->CreateFile '%s'", filename)
	audioFile, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if e != nil {
		return e
	}
	ttsUrl := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(message), ttsDefaultLanguage)
	httpResponse, e := http.Get(ttsUrl)
	if e != nil {
		_ = audioFile.Close()
		return e
	}
	_, e = io.Copy(audioFile, httpResponse.Body)
	_ = httpResponse.Body.Close()
	_ = audioFile.Close()
	if e != nil {
		return e
	}
	return nil
}

func (t *TextToSpeech) Create(name string, message string) {
	log.Printf("Audio->Create '%s' for '%s'", name, message)
	e := t.CreateFile(t.GetFileName(name), message)
	if e != nil {
		log.Printf("Audio->Create ERROR: '%s'", e.Error())
	}
}

func (t *TextToSpeech) CreateTemp(message string) (string, error) {
	tempFile, e := os.CreateTemp(os.TempDir(), "*.mp3")
	if e != nil {
		return "", e
	}
	tempFileName := tempFile.Name()
	_ = tempFile.Close()
	e = t.CreateFile(tempFileName, message)
	if e != nil {
		_ = os.Remove(tempFileName)
		return "", e
	}
	return tempFileName, nil
}

//goland:noinspection GoUnusedExportedFunction,GoUnusedExportedFunction
func (t *TextToSpeech) PlayOnce(message string) {
	go func() {
		tempFileName, e := t.CreateTemp(message)
		if e != nil {
			log.Println("Audio->PlayOnce ERROR: " + e.Error())
		} else {
			_ = t.PlayFile(tempFileName)
			_ = os.Remove(tempFileName)
		}
	}()
}

func (t *TextToSpeech) GenerateAudioFiles() {
	log.Printf("Audio->GenerateAudioFiles for '%s' (%s)", t.name, t.lang)
	name := "about"
	mesg := fmt.Sprintf("Hello dear! My name is %s, I am so pleased to meet you!", t.name)
	t.Create(name, mesg)
	for i := 1; i <= 3; i++ {
		name = fmt.Sprintf("target-%02d-hours", i)
		mesg = fmt.Sprintf("%d hours left", i)
		t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 1; i < 10; i += 1 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("%d minutes left", i)
		t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 10; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("%d minutes left", i)
		t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 1; i < 10; i += 1 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("%d", i)
		t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
	for i := 10; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("%d seconds left", i)
		t.Create(name, mesg)
		time.Sleep(750 * time.Millisecond)
	}
}

func GenerateAllAudioFiles(mainAppName string, audioPath string) {
	log.Printf("GenerateAllAudioFiles for '%s' in '%s'", mainAppName, audioPath)
	NewTextToSpeech(mainAppName, audioPath, "en").GenerateAudioFiles()
}
