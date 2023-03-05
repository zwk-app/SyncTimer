package logs

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type Level int

const CriticalLevel = 1
const ErrorLevel = 2
const WarningLevel = 3
const InfoLevel = 4
const DebugLevel = 5

func LevelName(level Level) string {
	list := map[Level]string{
		DebugLevel:    "DEBUG",
		InfoLevel:     "INFO",
		WarningLevel:  "WARN",
		ErrorLevel:    "ERROR",
		CriticalLevel: "CRIT",
	}
	return list[level]
}

type Logs struct {
	level    Level
	fileName string
	stdOut   bool
}

var logger *Logs = nil

func Logger() *Logs {
	if logger == nil {
		logger = new(Logs)
		logger.stdOut = false
		logger.fileName = ""
		logger.level = DebugLevel
	}
	return logger
}

func SetLevel(level Level) {
	Debug("Logs", fmt.Sprintf("SetLevel: %s", LevelName(level)), nil)
	Logger().level = level
}

func SetStdOut(stdOut bool) {
	Debug("Logs", fmt.Sprintf("SetStdOut: '%t'", stdOut), nil)
	if stdOut {
		Info("", fmt.Sprintf("Using StdOut"), nil)
		log.SetOutput(os.Stdout)
		Logger().stdOut = stdOut
	}
}

func SetFileName(fileName string) {
	Debug("Logs", fmt.Sprintf("SetFileName: '%s'", fileName), nil)
	if len(fileName) > 0 {
		f, e := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if e == nil {
			Info("", fmt.Sprintf("Using '%s'", fileName), nil)
			log.SetOutput(f)
			Logger().fileName = fileName
		} else {
			Warn("Logs", fmt.Sprintf("SetFileName: cannot write in file '%s'", fileName), e)
		}
	}
}

func contextMethod() string {
	pc, _, _, _ := runtime.Caller(4)
	runtimeContext := runtime.FuncForPC(pc).Name()
	return strings.Split(runtimeContext, "/")[1]
}

func formatError(e error) string {
	if e != nil {
		return fmt.Sprintf("%s: %s", contextMethod(), e.Error())
	}
	return ""
}

func formatLog(level Level, title string, message string) string {
	mesg := fmt.Sprintf("[%s]", LevelName(level))
	mesg = fmt.Sprintf("%-8s", mesg)
	if len(title) > 0 && len(message) > 0 {
		mesg += fmt.Sprintf("%-20s %s", title, message)
	} else if len(message) > 0 {
		mesg += fmt.Sprintf("%-20s %s", contextMethod(), message)
	} else {
		mesg = ""
	}
	return mesg
}

func logMessage(level Level, title string, message string, e error) {
	if level <= Logger().level {
		if len(message) == 0 && e != nil {
			message = e.Error()
		}
		if mesg := formatLog(level, title, message); len(mesg) > 0 {
			log.Printf("%s", mesg)
		}
		if e != nil {
			if mesg := formatLog(ErrorLevel, title, e.Error()); len(mesg) > 0 {
				log.Printf("%s", mesg)
			}
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n", formatError(e)))
		}
	}
}

func Debug(title string, message string, e error) {
	logMessage(DebugLevel, title, message, e)
}

func Info(title string, message string, e error) {
	logMessage(InfoLevel, title, message, e)
}

func Warn(title string, message string, e error) {
	logMessage(WarningLevel, title, message, e)
}

func Error(title string, message string, e error) {
	logMessage(ErrorLevel, title, message, e)
}

func CriticalExit(title string, message string, e error) {
	logMessage(CriticalLevel, title, message, e)
	os.Exit(1)
}
