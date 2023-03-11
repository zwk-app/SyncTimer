package audio

import (
	"fmt"
	"github.com/zwk-app/go-tools/logs"
	"io"
	"net/http"
	"net/url"
	"os"
)

//goland:noinspection GoNameStartsWithPackageName
type AudioTextToSpeech struct {
	language string
}

var textToSpeech *AudioTextToSpeech

func TextToSpeech() *AudioTextToSpeech {
	if textToSpeech == nil {
		textToSpeech = &AudioTextToSpeech{}
		textToSpeech.language = "en"
	}
	return textToSpeech
}

func SetLanguage(language string) {
	switch language {
	case "en":
		TextToSpeech().language = "en"
	default:
		TextToSpeech().language = "en"
	}
}

func CreateFile(filename string, message string) {
	logs.Debug("AudioTextToSpeech", fmt.Sprintf("CreateFile: '%s' in '%s'", message, filename), nil)
	audioFile, e := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if e != nil {
		logs.Error("AudioTextToSpeech", "", e)
		SetError(e)
		return
	}
	ttsUrl := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(message), TextToSpeech().language)
	httpResponse, e := http.Get(ttsUrl)
	if e != nil {
		_ = audioFile.Close()
		logs.Error("AudioTextToSpeech", "", e)
		SetError(e)
		return
	}
	_, e = io.Copy(audioFile, httpResponse.Body)
	_ = httpResponse.Body.Close()
	_ = audioFile.Close()
	if e != nil {
		logs.Error("AudioTextToSpeech", "", e)
		SetError(e)
	}
}

func CreateTemp(message string) string {
	logs.Debug("AudioTextToSpeech", fmt.Sprintf("CreateTemp: '%s'", message), nil)
	tempFile, e := os.CreateTemp(os.TempDir(), "*.mp3")
	if e != nil {
		SetError(e)
		logs.Error("AudioTextToSpeech", "", e)
		return ""
	}
	tempFileName := tempFile.Name()
	_ = tempFile.Close()
	CreateFile(tempFileName, message)
	if HasError() {
		_ = os.Remove(tempFileName)
		logs.Error("AudioTextToSpeech", "", Engine().last.error)
		return ""
	}
	return tempFileName
}
