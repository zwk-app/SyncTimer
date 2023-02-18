package ui

import (
	"SyncTimer/timer"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
)

var targetTimeDialogInitialized = false
var targetTimeDialog dialog.Dialog
var targetTimeDialogFormItems []*widget.FormItem
var targetInput *TargetEntry

type TargetEntry struct {
	widget.Entry
}

func NewTargetEntry() *TargetEntry {
	targetEntry := &TargetEntry{}
	targetEntry.ExtendBaseWidget(targetEntry)
	return targetEntry
}

func (i *TargetEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		i.Entry.TypedRune(r)
	}
}

func (i *TargetEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		i.Entry.TypedShortcut(shortcut)
		return
	}
	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		i.Entry.TypedShortcut(shortcut)
	}
}

func (i *TargetEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func TargetInputOnChange(s string) {
	valid := timer.CheckTimeString(s)
	sLen := len(s)
	switch sLen {
	case 1, 2, 4, 6:
		if !valid {
			targetInput.SetText(s[:len(s)-1])
		}
	default:
		if sLen > 6 {
			targetInput.SetText(s[0:6])
		}
	}
}

func TargetTimeDialogOnExit(b bool) {
	log.Println("TargetTimeDialogOnExit")
	if b {
		valid := appEngine.Timer.Object.SetTargetString(targetInput.Text)
		if !valid {
			log.Printf("TargetTimeDialogOnExit invalid: %s", targetInput.Text)
		}
	}
}

func TargetTimeDialogInit() {
	log.Println("TargetTimeDialogInit")
	if !targetTimeDialogInitialized {
		targetInput = NewTargetEntry()
		targetInput.OnChanged = TargetInputOnChange
		targetTimeDialogFormItems = append(targetTimeDialogFormItems, widget.NewFormItem("hh[mm[ss]]", targetInput))
		targetTimeDialog = dialog.NewForm("Target time", "OK", "Cancel",
			targetTimeDialogFormItems, TargetTimeDialogOnExit, appEngine.Fyne.MainWindow)
		targetTimeDialogInitialized = true
	}
}

func TargetTimeDialogShow() {
	log.Println("TargetTimeDialogShow")
	TargetTimeDialogInit()
	targetInput.SetText("")
	targetTimeDialog.Show()
}
