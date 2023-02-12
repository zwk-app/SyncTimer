package tts

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

var embeddedAudioFS *embed.FS
var embeddedAudioPath string

func SetEmbeddedAudioFS(embeddedFS *embed.FS, embeddedPath string) {
	embeddedAudioFS = embeddedFS
	embeddedAudioPath = embeddedPath
}

func CreateAudioFile(filename string, message string) (string, error) {
	log.Printf("CreateAudioFile %s", filename)
	audioFile, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if e != nil {
		return "", e
	}
	ttsUrl := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(message), ttsDefaultLanguage)
	httpResponse, e := http.Get(ttsUrl)
	if e != nil {
		_ = audioFile.Close()
		return "", e
	}
	_, e = io.Copy(audioFile, httpResponse.Body)
	_ = httpResponse.Body.Close()
	_ = audioFile.Close()
	if e != nil {
		return "", e
	}
	return audioFile.Name(), nil
}

func CreateTempAudioFile(message string) (string, error) {
	tempFile, e := os.CreateTemp(os.TempDir(), "*.mp3")
	if e != nil {
		return "", e
	}
	tempFileName := tempFile.Name()
	_ = tempFile.Close()
	_, e = CreateAudioFile(tempFileName, message)
	if e != nil {
		_ = os.Remove(tempFileName)
		return "", e
	}
	return tempFileName, nil
}

func PlayAudioFile(fileName string) error {
	embeddedFileName := embeddedAudioPath + filepath.Base(fileName)
	fileBytes, e := embeddedAudioFS.ReadFile(embeddedFileName)
	if e == nil {
		log.Printf("PlayAudioFile embedded:%s", embeddedFileName)
	} else {
		fileBytes, e = os.ReadFile(fileName)
		if e != nil {
			log.Printf("PlayAudioFile ERROR:%s", e.Error())
			return e
		}
		log.Printf("PlayAudioFile local:%s", fileName)
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

func Play(message string) {
	go func() {
		tempFileName, e := CreateTempAudioFile(message)
		if e != nil {
			log.Println("TextToSpeech : CreateTempAudioFile error: " + e.Error())
		} else {
			_ = PlayAudioFile(tempFileName)
			_ = os.Remove(tempFileName)
		}
	}()
}

type TextToSpeech struct {
	name string
	path string
	lang string
}

var TextToSpeechEngine *TextToSpeech

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

func (tts *TextToSpeech) GetFileName(name string) string {
	return tts.path + name + "." + tts.lang + ".mp3"
}

func (tts *TextToSpeech) Create(name string, message string) {
	log.Printf("TextToSpeech->Create '%s' for '%s'", name, message)
	_, e := CreateAudioFile(tts.GetFileName(name), message)
	if e != nil {
		log.Printf("TextToSpeech->Create ERROR: '%s'", e.Error())
	}
}

func (tts *TextToSpeech) Play(name string) {
	log.Printf("TextToSpeech->Play '%s'", name)
	e := PlayAudioFile(tts.GetFileName(name))
	if e != nil {
		log.Printf("TextToSpeech->Play ERROR: '%s'", e.Error())
	}
}

func (tts *TextToSpeech) GenerateAudioFiles() {
	log.Printf("TextToSpeech->GenerateAudioFiles for '%s' (%s)", tts.name, tts.lang)
	name := "about"
	mesg := fmt.Sprintf("Hello dear! My name is %s, I am so pleased to meet you!", tts.name)
	tts.Create(name, mesg)
	for i := 1; i <= 3; i++ {
		name = fmt.Sprintf("target-%02d-hours", i)
		mesg = fmt.Sprintf("Target in %d hours", i)
		tts.Create(name, mesg)
		time.Sleep(500 * time.Millisecond)
	}
	for i := 5; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-minutes", i)
		mesg = fmt.Sprintf("Target in %d minutes", i)
		tts.Create(name, mesg)
		time.Sleep(500 * time.Millisecond)
	}
	for i := 5; i <= 60; i += 5 {
		name = fmt.Sprintf("target-%02d-seconds", i)
		mesg = fmt.Sprintf("Target in %d seconds", i)
		tts.Create(name, mesg)
		time.Sleep(500 * time.Millisecond)
	}
}

func GenerateAllAudioFiles(mainAppName string, audioPath string) {
	log.Printf("GenerateAllAudioFiles for '%s' in '%s'", mainAppName, audioPath)
	NewTextToSpeech(mainAppName, audioPath, "en").GenerateAudioFiles()
}
