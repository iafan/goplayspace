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

// SetValue sets texatrea value (text)
func (t *Textarea) SetValue(val string) {
	ss := t.GetSelectionStart()
	se := t.GetSelectionEnd()
	t.Set("value", val)
	t.SetSelectionStart(ss)
	t.SetSelectionEnd(se)
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

	t.SetSelectionStart(ss + len(text))
	t.SetSelectionEnd(ss + len(text)) // the same as start
}
