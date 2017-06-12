package textarea

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
)

// Textarea provides wrapper functions for <textarea> object
type Textarea struct {
	*js.Object
}

// Focus sets focus to the control
func (t *Textarea) Focus() {
	t.Call("focus")
}

// GetSelectionStart returns selection start character
func (t *Textarea) GetSelectionStart() int {
	return t.Get("selectionStart").Int()
}

// SetSelectionStart sets selection start character
func (t *Textarea) SetSelectionStart(val int) {
	t.Set("selectionStart", val)
}

// GetSelectionEnd returns selection end character
func (t *Textarea) GetSelectionEnd() int {
	return t.Get("selectionEnd").Int()
}

// SetSelectionEnd sets selection end character
func (t *Textarea) SetSelectionEnd(val int) {
	t.Set("selectionEnd", val)
}

// GetValue returns current texatrea value (text)
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

	if ss > 0 {
		before = val[ss-1 : ss]
	}

	if se < len(val) {
		after = val[se : se+1]
	}

	return before, after
}

// SetState sets texatrea value (text) and selection
func (t *Textarea) SetState(val string, selStart, selEnd int) {
	t.Set("value", val)
	t.SetSelectionStart(selStart)
	t.SetSelectionEnd(selEnd)
}

// SetValue sets texatrea value (text)
// while preserving selection
func (t *Textarea) SetValue(val string) {
	t.SetState(val, t.GetSelectionStart(), t.GetSelectionEnd())
}

// SetHeight sets texatrea height in pixels
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
