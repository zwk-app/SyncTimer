package audio

import (
	"bytes"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/zwk-app/go-tools/logs"
	"os"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type AudioPlayer struct {
	channelCount    int
	bitDepthInBytes int
	last            struct {
		error error
	}
}

var player *AudioPlayer

func Player() *AudioPlayer {
	if player == nil {
		player = &AudioPlayer{}
		player.channelCount = 2
		player.bitDepthInBytes = 2
	}
	return player
}

func PlayMp3Content(fileBytes []byte) {
	mp3Decoder, e := mp3.NewDecoder(bytes.NewReader(fileBytes))
	if e != nil {
		logs.Error("AudioPlayer", "", e)
		SetError(e)
		return
	}
	otoContext, readyChan, e := oto.NewContext(mp3Decoder.SampleRate(), Player().channelCount, Player().bitDepthInBytes)
	if e != nil {
		logs.Error("AudioPlayer", "", e)
		SetError(e)
		return
	}
	<-readyChan
	mp3Player := otoContext.NewPlayer(mp3Decoder)
	mp3Player.Play()
	for mp3Player.IsPlaying() {
		time.Sleep(200 * time.Millisecond)
	}
	SetError(mp3Player.Close())
}

func PlayMp3File(fileName string) {
	fileBytes, e := os.ReadFile(fileName)
	if e != nil {
		logs.Error("AudioPlayer", "", e)
		SetError(e)
		return
	}
	PlayMp3Content(fileBytes)
}
