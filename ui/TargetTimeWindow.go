package ui

import (
	"SyncTimer/timer"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
)

var targetWindowInitialized = false
var targetWindowContainer *fyne.Container
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

func TargetWindowOnClose() {
	log.Printf("TargetWindowOnClose")
	appEngine.Fyne.MainWindow.Canvas().SetOnTypedKey(nil)
	ShowMainWindow()
}

func TargetConfirmButtonOnClick() {
	log.Printf("TargetConfirmButtonOnClick")
	appEngine.SetTargetTime(targetInput.Text)
	TargetWindowOnClose()
}

func TargetInputAppend(i int) {
	targetInput.SetText(fmt.Sprintf("%s%d", targetInput.Text, i))
}

func TargetInputBackspace() {
	targetInput.SetText(targetInput.Text[:len(targetInput.Text)-1])
}

func TargetWindowOnKeyEvent(k *fyne.KeyEvent) {
	log.Printf("TargetWindowOnKeyEvent: '%s'", k.Name)
	switch k.Name {
	case fyne.KeyEscape:
		TargetWindowOnClose()
	case fyne.KeyReturn, fyne.KeyEnter:
		TargetConfirmButtonOnClick()
	case fyne.Key0:
		TargetInputAppend(0)
	case fyne.Key1:
		TargetInputAppend(1)
	case fyne.Key2:
		TargetInputAppend(2)
	case fyne.Key3:
		TargetInputAppend(3)
	case fyne.Key4:
		TargetInputAppend(4)
	case fyne.Key5:
		TargetInputAppend(5)
	case fyne.Key6:
		TargetInputAppend(6)
	case fyne.Key7:
		TargetInputAppend(7)
	case fyne.Key8:
		TargetInputAppend(8)
	case fyne.Key9:
		TargetInputAppend(9)
	case fyne.KeyBackspace:
		TargetInputBackspace()
	}

}

func TargetWindowContent() *fyne.Container {
	log.Printf("TargetWindowContent")

	if !targetWindowInitialized {

		/* Middle content */
		targetInput = NewTargetEntry()
		targetInput.OnChanged = TargetInputOnChange
		targetInput.SetPlaceHolder("hh[mm[ss]]")
		targetContainer := container.NewCenter(targetInput)

		/* Bottom buttons */
		targetCancelButton := widget.NewButton("Cancel", TargetWindowOnClose)
		targetCancelButton.SetIcon(theme.CancelIcon())
		targetConfirmButton := widget.NewButton("Confirm", TargetConfirmButtonOnClick)
		targetConfirmButton.SetIcon(theme.ConfirmIcon())
		bottomContainer := container.NewGridWithRows(1, targetCancelButton, targetConfirmButton)

		targetWindowContainer = container.NewBorder(nil, bottomContainer, nil, nil, targetContainer)

		targetWindowInitialized = true
	}
	return targetWindowContainer
}
func ShowTargetWindow() {
	log.Printf("ShowTargetWindow")
	appEngine.Fyne.MainWindow.SetContent(TargetWindowContent())
	targetInput.SetText("")
	appEngine.Fyne.MainWindow.Show()
	appEngine.Fyne.MainWindow.Canvas().SetOnTypedKey(TargetWindowOnKeyEvent)
}
