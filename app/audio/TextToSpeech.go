package audio

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type TextToSpeech struct {
	language string
}

func NewTextToSpeech(language string) *TextToSpeech {
	t := TextToSpeech{}
	t.SetLanguage(language)
	return &t
}

func (t *TextToSpeech) SetLanguage(language string) *TextToSpeech {
	switch language {
	case "en":
		t.language = "en"
	default:
		t.language = "en"
	}
	return t
}

func (t *TextToSpeech) CreateFile(filename string, message string) error {
	log.Printf("TextToSpeech->CreateFile '%s'", filename)
	audioFile, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if e != nil {
		return e
	}
	ttsUrl := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(message), t.language)
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
