package app

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type LogsLevel int

const LogsCritical = 1
const LogsError = 2
const LogsWarning = 3
const LogsInfo = 4
const LogsDebug = 5

func LogsLevelName(level LogsLevel) string {
	list := map[LogsLevel]string{
		LogsDebug:    "DEBUG",
		LogsInfo:     "INFO",
		LogsWarning:  "WARN",
		LogsError:    "ERROR",
		LogsCritical: "CRIT",
	}
	return list[level]
}

type Logs struct {
	App   *AppEngine
	level LogsLevel
}

func NewAppLogs(appEngine *AppEngine, level LogsLevel) *Logs {
	r := &Logs{}
	r.App = appEngine
	r.level = level
	r.uninitializedExit()
	return r.SetFileName(r.App.Config.Logs.FileName).SetStdOut(r.App.Config.Logs.StdOut)
}

func (r *Logs) uninitializedExit() {
	if r.App == nil {
		pc, _, _, _ := runtime.Caller(2)
		runtimeContext := runtime.FuncForPC(pc).Name()
		runtimeShortContext := strings.Split(runtimeContext, "/")[1]
		_, _ = os.Stderr.WriteString(fmt.Sprintf("%s: Critical: AppEngine is not set", runtimeShortContext))
		os.Exit(1)
	}
}

func (r *Logs) SetStdOut(stdOut bool) *Logs {
	r.uninitializedExit()
	r.Debug("AppLogs", fmt.Sprintf("SetStdOut: '%t'", stdOut), nil)
	if stdOut {
		r.Info("", fmt.Sprintf("Using StdOut"), nil)
		log.SetOutput(os.Stdout)
	}
	r.App.Config.Logs.StdOut = stdOut
	return r
}

func (r *Logs) SetFileName(fileName string) *Logs {
	r.uninitializedExit()
	r.Debug("AppLogs", fmt.Sprintf("SetFileName: '%s'", fileName), nil)
	if len(fileName) > 0 {
		f, e := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if e == nil {
			r.Info("", fmt.Sprintf("Using '%s'", fileName), nil)
			log.SetOutput(f)
		} else {
			LogsCriticalErrorExit("AppLogs", fmt.Sprintf("SetFileName: cannot write in file %s", fileName), e)
		}
	}
	r.App.Config.Logs.FileName = fileName
	return r
}

func LogsFormat(level LogsLevel, title string, message string, e error) string {
	mesg := fmt.Sprintf("[%s]", LogsLevelName(level))
	mesg = fmt.Sprintf("%-8s", mesg)
	if len(title) >= 0 {
		mesg += fmt.Sprintf("%-16s %s", title, message)
	} else {
		mesg += message
	}
	if level <= LogsCritical && e != nil {
		pc, _, _, _ := runtime.Caller(2)
		runtimeContext := runtime.FuncForPC(pc).Name()
		runtimeShortContext := strings.Split(runtimeContext, "/")[1]
		mesg += "\n" + runtimeShortContext + ": " + e.Error()
	}
	return mesg
}

func (r *Logs) append(level LogsLevel, title string, message string, e error) {
	if r.level <= level {
		mesg := LogsFormat(level, title, message, e)
		log.Printf("%s", mesg)
		if e != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("%s\n%s", mesg, e.Error()))
		}
	}
}

func (r *Logs) Error(title string, message string, e error) {
	r.append(LogsError, title, message, e)
}

func (r *Logs) Warn(title string, message string, e error) {
	r.append(LogsWarning, title, message, e)
}

func (r *Logs) Info(title string, message string, e error) {
	r.append(LogsInfo, title, message, e)
}

func (r *Logs) Debug(title string, message string, e error) {
	r.append(LogsDebug, title, message, e)
}

func LogsCriticalErrorExit(title string, message string, e error) {
	_, _ = os.Stderr.WriteString(LogsFormat(LogsCritical, title, message, e))
	os.Exit(1)
}
