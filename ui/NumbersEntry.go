package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type NumbersEntry struct {
	widget.Entry
}

func NewNumbersEntry() *NumbersEntry {
	numbersEntry := &NumbersEntry{}
	numbersEntry.ExtendBaseWidget(numbersEntry)
	return numbersEntry
}

func (t *NumbersEntry) MinSize() fyne.Size {
	return fyne.NewSize(128, 36)
}

func (t *NumbersEntry) TypedRune(r rune) {
	if r >= '0' && r <= '9' {
		t.Entry.TypedRune(r)
	}
}

func (t *NumbersEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		t.Entry.TypedShortcut(shortcut)
		return
	}
	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		t.Entry.TypedShortcut(shortcut)
	}
}

func (t *NumbersEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}
