package resources

import (
	"embed"
	"fmt"
	"github.com/zwk-app/zwk-tools/logs"
	"strings"
)

type FileResources struct {
	embedded struct {
		FS        *embed.FS
		audioPath string
		imagePath string
	}
}

var resources *FileResources

func Resources() *FileResources {
	if resources == nil {
		resources = &FileResources{}
		resources.embedded.FS = nil
		resources.embedded.audioPath = EmbeddedPath("res/audio")
		resources.embedded.imagePath = EmbeddedPath("res/images")
	}
	return resources
}

func SetEmbedded(embeddedFs *embed.FS) {
	Resources().embedded.FS = embeddedFs
}

func EmbeddedPath(embeddedPath string) string {
	if !strings.HasSuffix(embeddedPath, "/") {
		embeddedPath += "/"
	}
	return embeddedPath
}

func ReadFile(fileName string) []byte {
	logs.Debug("Resources", fmt.Sprintf("ReadFile: '%s'", fileName), nil)
	if Resources().embedded.FS == nil {
		logs.CriticalExit("Resources", "ReadFile: Embedded FS not set", nil)
	}
	b, e := Resources().embedded.FS.ReadFile(fileName)
	if e != nil {
		logs.Error("Resources", "ReadFile", e)
		return nil
	}
	return b
}

func ReadAudio(fileName string) []byte {
	logs.Debug("Resources", fmt.Sprintf("ReadAudioFile: '%s'", fileName), nil)
	return ReadFile(Resources().embedded.audioPath + fileName)
}

func ReadImage(fileName string) []byte {
	logs.Debug("Resources", fmt.Sprintf("ReadAudioFile: '%s'", fileName), nil)
	return ReadFile(Resources().embedded.imagePath + fileName)
}

func AlarmSoundNames() []string {
	return []string{
		"navy-01-wanking-of-combat",
		"navy-02-breakfast",
		"navy-12-lunch-time",
		"navy-14-wake-up",
		"navy-22-wanking-of-combat",
	}
}

func AlarmSoundTitles() []string {
	return []string{
		"Wanking of Combat (Morning)",
		"Breakfast Time",
		"Lunch Time",
		"Wake Up",
		"Wanking of Combat (Evening)",
	}
}

func AlarmSoundTitle(alarmSoundName string) string {
	for i, v := range AlarmSoundNames() {
		if v == alarmSoundName {
			return AlarmSoundTitles()[i]
		}
	}
	return ""
}

func AlarmSoundName(alarmSoundTitle string) string {
	for i, v := range AlarmSoundTitles() {
		if v == alarmSoundTitle {
			return AlarmSoundNames()[i]
		}
	}
	return ""
}
