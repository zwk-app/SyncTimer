package audio

import (
	"bytes"
	"embed"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type AudioPlayer struct {
	embeddedFS      *embed.FS
	embeddedPath    string
	channelCount    int
	bitDepthInBytes int
}

func NewAudioPlayer(embeddedFS *embed.FS, embeddedAudioPath string) *AudioPlayer {
	p := AudioPlayer{}
	p.embeddedFS = embeddedFS
	p.embeddedPath = embeddedAudioPath
	p.channelCount = 2
	p.bitDepthInBytes = 2
	return &p
}

func (p *AudioPlayer) playMp3Content(fileBytes []byte) error {
	mp3Decoder, e := mp3.NewDecoder(bytes.NewReader(fileBytes))
	if e != nil {
		return e
	}
	otoContext, readyChan, e := oto.NewContext(mp3Decoder.SampleRate(), p.channelCount, p.bitDepthInBytes)
	if e != nil {
		return e
	}
	<-readyChan
	mp3Player := otoContext.NewPlayer(mp3Decoder)
	mp3Player.Play()
	for mp3Player.IsPlaying() {
		time.Sleep(200 * time.Millisecond)
	}
	return mp3Player.Close()
}

func (p *AudioPlayer) Play(fileName string) error {
	embeddedFileName := path.Clean(p.embeddedPath + string(os.PathSeparator) + filepath.Base(fileName))
	_, e := os.Stat(embeddedFileName)
	if e == nil {
		embeddedFileBytes, e := p.embeddedFS.ReadFile(embeddedFileName)
		if e == nil {
			log.Printf("AudioPlayer->Play embedded:%s", embeddedFileName)
			return p.playMp3Content(embeddedFileBytes)
		}
	}
	_, e = os.Stat(fileName)
	if e == nil {
		localFileBytes, e := os.ReadFile(fileName)
		if e == nil {
			log.Printf("AudioPlayer->Play local:%s", fileName)
			return p.playMp3Content(localFileBytes)
		}
	}
	return e
}
