package textarea

import (
	"strconv"

	"unicode/utf16"

	"github.com/gopherjs/gopherjs/js"
	"github.com/iafan/goplayspace/client/js/str"
)

// Textarea provides wrapper functions for <textarea> object
type Textarea struct {
	*js.Object
}

// Focus sets focus to the control
func (t *Textarea) Focus() {
	t.Call("focus")
}

// GetSymbolWidthsAround returns the size of the unicode symbols
// in utf8 bytes before and after a given utf8 byte range in a string
func (t *Textarea) GetSymbolWidthsAround(before, after int) (wBefore, wAfter int) {
	u16 := utf16.Encode([]rune(t.GetValue()))

	if before <= len(u16) {
		if r := utf16.Decode(u16[:before]); len(r) > 0 {
			wBefore = len(string(r[len(r)-1]))
		}
	}

	if after < len(u16) {
		if r := utf16.Decode(u16[after:]); len(r) > 0 {
			wAfter = len(string(r[0]))
		}
	}

	return
}

func (t *Textarea) utf16ToUTF8Pos(i int) int {
	return str.UTF16ToUTF8Pos(t.GetValue(), i)
}

func (t *Textarea) utf8ToUTF16Pos(i int) int {
	return str.UTF8ToUTF16Pos(t.GetValue(), i)
}

// GetSelectionStart returns selection start utf8 byte position
// Note that JavaScript strings are UTF-16-encoded,
// and selectionStart returns a position in the UTF-16 array
// representing the textarea value
func (t *Textarea) GetSelectionStart() int {
	return t.utf16ToUTF8Pos(t.Get("selectionStart").Int())
}

// SetSelectionStart sets selection start character
func (t *Textarea) SetSelectionStart(val int) {
	t.Set("selectionStart", t.utf8ToUTF16Pos(val))
}

// GetSelectionEnd returns selection end utf8 byte position
// Note that JavaScript strings are UTF-16-encoded,
// and selectionStart returns a position in the UTF-16 array
// representing the textarea value
func (t *Textarea) GetSelectionEnd() int {
	return t.utf16ToUTF8Pos(t.Get("selectionEnd").Int())
}

// SetSelectionEnd sets selection end character
func (t *Textarea) SetSelectionEnd(val int) {
	t.Set("selectionEnd", t.utf8ToUTF16Pos(val))
}

// GetValue returns current textarea value (text)
func (t *Textarea) GetValue() string {
	return t.Get("value").String()
}

// GetSymbolsAroundSelection returns one symbol before and one symbol
// after selection (or around the caret, if there's no selection);
// one or both strings can be empty if caret is placed at the beginning
// or the end of the document, or if the document is blank
func (t *Textarea) GetSymbolsAroundSelection() (before, after string) {
	ss := t.GetSelectionStart()
	se := t.GetSelectionEnd()
	val := t.GetValue()

	wb, wa := t.GetSymbolWidthsAround(ss, se)

	if ss > 0 {
		before = val[ss-wb : ss]
	}

	if se < len(val) {
		after = val[se : se+wa]
	}

	return before, after
}

// SetState sets textarea value (text) and selection
func (t *Textarea) SetState(val string, selStart, selEnd int) {
	t.Set("value", val)
	t.SetSelectionStart(selStart)
	t.SetSelectionEnd(selEnd)
}

// SetValue sets textarea value (text)
// while preserving selection
func (t *Textarea) SetValue(val string) {
	t.SetState(val, t.GetSelectionStart(), t.GetSelectionEnd())
}

// SetHeight sets textarea height in pixels
func (t *Textarea) SetHeight(val int) {
	t.Call("setAttribute", "style", "height:"+strconv.Itoa(val)+"px")
}

// InsertText replaces selection with the provided text
// And adjusts the caret pos
func (t *Textarea) InsertText(text string) {
	ss := t.GetSelectionStart()
	se := t.GetSelectionEnd()

	val := t.GetValue()
	val = val[:ss] + text + val[se:]
	t.Set("value", val)

	ss = ss + len(text)
	t.SetSelectionStart(ss)
	t.SetSelectionEnd(ss) // the same as start
}

// WrapSelection wraps selection with the provided
// starting and ending text snippets
// and places the caret before the `end` part
func (t *Textarea) WrapSelection(begin, end string) {
	ss := t.GetSelectionStart()
	se := t.GetSelectionEnd()

	val := t.GetValue()
	val = val[:ss] + begin + val[ss:se] + end + val[se:]
	t.Set("value", val)

	t.SetSelectionStart(ss + len(begin))
	t.SetSelectionEnd(se + len(begin))
}
