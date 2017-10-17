package editor

import (
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"

	"github.com/iafan/goplayspace/client/component/editor/undo"
	"github.com/iafan/goplayspace/client/js/console"
	"github.com/iafan/goplayspace/client/js/document"
	"github.com/iafan/goplayspace/client/js/textarea"
	"github.com/iafan/goplayspace/client/ranges"
	"github.com/iafan/goplayspace/client/util"
)

// saveStateTimeout defines how much time should pass after the last
// onChange event for state to be saved to undo stack
const saveStateTimeout = 500 * time.Millisecond

// Editor implements editor logic
type Editor struct {
	vecty.Core

	ta          *textarea.Textarea
	sh          *Shadow
	shiftDown   bool
	ctrlDown    bool
	metaDown    bool
	highlighted string
	selLinesCSS string
	errorsCSS   string
	warningsCSS string

	Range            *ranges.Range   `vecty:"prop"`
	HighlightingMode bool            `vecty:"prop"`
	ReadonlyMode     bool            `vecty:"prop"`
	ErrorLines       map[string]bool `vecty:"prop"`
	WarningLines     map[string]bool `vecty:"prop"`
	UndoStack        *undo.Stack     `vecty:"prop"`
	ChangeTimer      **time.Timer    // note this is a pointer to a pointer

	Highlighter     func(s string) string `vecty:"prop"`
	OnTopicChange   func(topic string)
	OnChange        func(value string)
	OnLineSelChange func(value string)
	OnKeyDown       func(e *vecty.Event)
}

// Focus sets focus to the control
func (ed *Editor) Focus() {
	if ed.ta == nil {
		console.Log("editor.Focus(): getTextarea() is nil")
		return
	}
	util.Schedule(ed.ta.Focus)
}

// GetSelection gets text selection
func (ed *Editor) GetSelection() (start, end int) {
	if ed.ta == nil {
		return -1, -1
	}
	return ed.ta.GetSelectionStart(), ed.ta.GetSelectionEnd()
}

// SetSelection sets text selection
func (ed *Editor) SetSelection(start, end int) {
	if ed.ta == nil {
		return
	}
	ed.ta.SetSelectionStart(start)
	ed.ta.SetSelectionEnd(end)
}

func (ed *Editor) updateSelectionInfo(e *vecty.Event) {
	if ed.ta == nil || ed.OnTopicChange == nil {
		return
	}
	ss := ed.ta.GetSelectionStart()
	se := ed.ta.GetSelectionEnd()
	text := ed.ta.GetValue()
	if se > len(text) {
		se = len(text)
	}
	sel := text[ss:se]

	if sel == "" {
		return
	}

	// FIXME: sel must be an alphanumeric sequence,
	// otherwise selection expansion should not be performed

	// test if there is a '.' symbol before the selection
	if ss > 0 && text[ss-1] == '.' {
		// go back until we get to non-alpha character to get the full package name
		start := ss - 2
		for i := start; i >= 0; i-- {
			ch := text[i : i+1]
			if strings.ToLower(ch) == strings.ToUpper(ch) {
				// we're at non-alpha char
				if i < start {
					// we've got a non-empty package name,
					// updating the selected text
					sel = text[i+1 : se]
				}
				break
			}
		}
	}

	ed.OnTopicChange(sel)
}

// ResizeTextarea resizes the height of the textarea
// to match the computed height of the shadow
func (ed *Editor) ResizeTextarea() {
	if ed.sh == nil || ed.ta == nil {
		return
	}

	ed.ta.SetHeight(ed.sh.GetHeight())
}

func (ed *Editor) makeHighlightedText(text string) string {
	a := strings.Split(text, "\n")
	for i, line := range a {
		a[i] = "<li>" + html.EscapeString(line) + "</li>\n"
	}

	return "<ol>\n" + strings.Join(a, "") + "</ol>"
}

// Highlight applies highlighting to the editor
func (ed *Editor) Highlight(on bool) {
	if ed.sh == nil || ed.ta == nil {
		console.Log("editor.Highlight(): getShadow() or getTextarea() is nil!")
		return
	}
	text := ed.ta.GetValue()
	ed.highlighted = ""
	if on && ed.Highlighter != nil {
		ed.highlighted = ed.Highlighter(text)
	}
	if ed.highlighted == "" {
		ed.highlighted = ed.makeHighlightedText(text)
	}
	ed.sh.SetValue(ed.highlighted)
	ed.ResizeTextarea()
}

func (ed *Editor) onChange(e *vecty.Event) {
	if ed.ta == nil {
		console.Log("editor.onChange(): getTextarea() is nil!")
		return
	}
	shouldFireSelChange := ed.Range != nil
	ed.Range = nil
	ed.WarningLines = nil
	ed.ErrorLines = nil
	ed.Highlight(ed.HighlightingMode)

	t := *ed.ChangeTimer
	if t == nil {
		t = time.AfterFunc(saveStateTimeout, ed.saveState)
		*ed.ChangeTimer = t
	} else {
		t.Stop()
		t.Reset(saveStateTimeout)
	}

	ed.fireOnChangeEvent()
	if shouldFireSelChange {
		ed.fireOnLineSelChangeEvent()
	}
}

func (ed *Editor) cancelEvent(e *vecty.Event) {
	e.Call("preventDefault")
	e.Call("stopPropagation")
}

// InsertText inserts text in place of selection
func (ed *Editor) InsertText(text string) {
	if ed.ta == nil {
		console.Log("editor.InsertText(): getTextarea() is nil!")
		return
	}
	ed.ta.InsertText(text)
	ed.onChange(nil)
}

// WrapSelection wraps selection with the provided
// starting and ending text snippets
func (ed *Editor) WrapSelection(begin, end string) {
	if ed.ta == nil {
		console.Log("editor.WrapSelection(): getTextarea() is nil!")
		return
	}
	ed.saveState()
	ed.ta.WrapSelection(begin, end)
	ed.saveState()
	ed.onChange(nil)
}

// SetText replaces the editor text
func (ed *Editor) SetText(text string) {
	if ed.ta == nil {
		console.Log("editor.SetText() getTextarea() is nil")
		return
	}
	ed.saveState()
	ed.ta.SetValue(text)
	ed.saveState()
	ed.onChange(nil)
}

// SetState replaces the editor text and sets selection
func (ed *Editor) SetState(text string, selStart, selEnd int) {
	if ed.ta == nil {
		console.Log("editor.SetState() getTextarea() is nil")
		return
	}
	ed.saveState()
	ed.ta.SetState(text, selStart, selEnd)
	ed.saveState()
	ed.onChange(nil)
}

func (ed *Editor) fireOnChangeEvent() {
	if ed.OnChange != nil {
		ed.OnChange(ed.ta.GetValue())
	}
}

func (ed *Editor) fireOnLineSelChangeEvent() {
	if ed.OnLineSelChange != nil {
		ed.OnLineSelChange(ed.Range.String())
	}
}

func (ed *Editor) resetLineSelection() {
	if ed.Range.HasSelection() {
		ed.Range.ClearSelection()
		ed.fireOnLineSelChangeEvent()
	}
}

func (ed *Editor) toggleLine(n int) {
	defer ed.fireOnLineSelChangeEvent()

	if ed.Range == nil {
		ed.Range = &ranges.Range{}
	}

	if ed.shiftDown {
		ed.Range.AddSelPoint(n)
		return
	}

	if ed.ctrlDown || ed.metaDown {
		ed.Range.ToggleLine(n)
		return
	}

	if ed.Range.IsOnlyLineSelected(n) {
		ed.Range.ToggleLine(n) // remove selection
	} else {
		ed.Range.SetRange(n, n) // reset selection to this line only
	}
}

func (ed *Editor) toggleLineSelection() {
	if ed.ta == nil {
		return
	}
	ss := ed.ta.GetSelectionStart()
	line := 1
	sel := ed.ta.GetValue()[:ss]
	for {
		i := strings.Index(sel, "\n")
		if i == -1 {
			break
		}
		line++
		sel = sel[i+1:]
	}

	ed.toggleLine(line)
}

func (ed *Editor) getIndent() int {
	if ed.ta == nil {
		return 0
	}
	ss := ed.ta.GetSelectionStart()
	s := ed.ta.GetValue()[:ss]
	i := strings.LastIndex(s, "\n")
	if i > 0 {
		s = s[i+1:]
	}
	for i = 0; i < len(s); i++ {
		if s[i] != '\t' {
			break
		}
	}
	before, _ := ed.ta.GetSymbolsAroundSelection()
	if strings.ContainsAny(before, "{([") {
		i++
	} else if before == "}" && i > 0 {
		i--
	}

	return i
}

func (ed *Editor) handleKeyDown(e *vecty.Event) {
	ed.shiftDown = e.Get("shiftKey").Bool()
	ed.ctrlDown = e.Get("ctrlKey").Bool()
	ed.metaDown = e.Get("metaKey").Bool()

	if ed.ta == nil {
		return
	}

	switch e.Get("keyCode").Int() {
	case 84: // T
		if ed.ctrlDown { // Ctrl+T
			e.Call("preventDefault")
			ed.toggleLineSelection()
			return
		}
	case 8: // Backspace
		before, after := ed.ta.GetSymbolsAroundSelection()

		insidePair := false
		switch before {
		case `"`, "'", "`":
			insidePair = before == after
		case "(":
			insidePair = after == ")"
		case "[":
			insidePair = after == "]"
		case "{":
			insidePair = after == "}"
		}

		if !insidePair {
			break
		}

		ss, se := ed.GetSelection()
		if ss != se || ss == -1 {
			break
		}

		e.Call("preventDefault")
		ed.SetSelection(ss-1, ss+1)
		ed.InsertText("")
		return
	case 9: // Tab
		e.Call("preventDefault")
		ed.InsertText("\t")
		return
	case 13: // Enter
		if !ed.shiftDown && !ed.ctrlDown && !ed.metaDown {
			e.Call("preventDefault")
			i := ed.getIndent()
			before, after := ed.ta.GetSymbolsAroundSelection()
			if before == "{" && after == "}" ||
				before == "(" && after == ")" ||
				before == "[" && after == "]" {
				iAfter := i - 1
				if iAfter < 0 {
					iAfter = 0
				}
				ed.WrapSelection(
					"\n"+strings.Repeat("\t", i),
					"\n"+strings.Repeat("\t", iAfter))
			} else {
				ed.InsertText("\n" + strings.Repeat("\t", i))
			}
			return
		}
	case 27: // Esc
		e.Call("preventDefault")
		ed.resetLineSelection()
		return
	case 89: // Y
		if ed.ctrlDown || ed.metaDown { // Ctrl+Y or Command+Y
			e.Call("preventDefault")
			ed.Redo()
			return
		}
	case 90: // Z
		if ed.ctrlDown || ed.metaDown {
			e.Call("preventDefault")
			if ed.shiftDown {
				ed.Redo() // Shift+Ctrl+Z or Shift+Command+Z
			} else {
				ed.Undo() // Ctrl+Z or Command+Z
			}
			return
		}
	}

	if ed.OnKeyDown != nil {
		ed.OnKeyDown(e)
	}
}

func (ed *Editor) handleKeyPress(e *vecty.Event) {
	if ed.ta == nil {
		return
	}
	before, after := ed.ta.GetSymbolsAroundSelection()
	canWrapQuotes := (before == "" || strings.ContainsAny(before, " \n{([:=")) &&
		(after == "" || strings.ContainsAny(after, " \n})]:="))
	canWrapBraces := after == "" || strings.ContainsAny(after, " \n})]:=")

	r := rune(e.Get("charCode").Int())
	rs := string(r)

	if canWrapQuotes {
		switch r {
		case '"', '\'', '`':
			e.Call("preventDefault")
			ed.WrapSelection(rs, rs)
		}
	}

	if canWrapBraces {
		switch r {
		case '(':
			e.Call("preventDefault")
			ed.WrapSelection("(", ")")
		case '[':
			e.Call("preventDefault")
			ed.WrapSelection("[", "]")
		case '{':
			e.Call("preventDefault")
			ed.WrapSelection("{", "}")
		}
	}

	switch r {
	case ')', ']', '}', '"', '\'', '`':
		if after != rs {
			break
		}

		ss, se := ed.GetSelection()
		if ss != se || ss == -1 {
			break
		}

		e.Call("preventDefault")
		ed.SetSelection(ss+1, ss+1)
	}
}

func (ed *Editor) handleShadowMouseDown(e *vecty.Event) {
	if e.Get("button").Int() != 0 {
		return
	}

	e.Call("preventDefault")

	ed.shiftDown = e.Get("shiftKey").Bool()
	ed.ctrlDown = e.Get("ctrlKey").Bool()
	ed.metaDown = e.Get("metaKey").Bool()

	ed.toggleLine(e.Get("target").Get("data-index").Int())
}

func (ed *Editor) handleScrollerClick(e *vecty.Event) {
	ed.Focus()
}

func (ed *Editor) afterRender() {
	list := js.Global.Get("document").Call("querySelectorAll", ".shadow ol li")
	if list == nil || list.Length() == 0 || ed.sh == nil || ed.ta == nil {
		time.AfterFunc(5*time.Millisecond, ed.afterRender)
		return
	}
	n := list.Length()
	for i := 0; i < n; i++ {
		list.Index(i).Set("onmousedown", ed.handleShadowMouseDown)
		list.Index(i).Set("data-index", i+1)
	}
}

func (ed *Editor) updateStateFromRanges() {
	ed.selLinesCSS = ""
	if ed.Range == nil {
		return
	}
	for _, r := range ed.Range.Sel {
		for i := r.Begin; i <= r.End; i++ {
			ed.selLinesCSS = ed.selLinesCSS +
				".shadow ol li:nth-child(" + strconv.Itoa(i) + ") {background: var(--sel-bgcolor)}\n" +
				".shadow ol li:nth-child(" + strconv.Itoa(i) + ")::before {background: var(--sel-bgcolor)}\n"
		}
	}
}

func (ed *Editor) updateStateFromErrors() {
	ed.errorsCSS = ""
	if ed.ErrorLines == nil {
		return
	}
	for key := range ed.ErrorLines {
		ed.errorsCSS = ed.errorsCSS + ".shadow ol li:nth-child(" + key + ") {background: var(--error-bgcolor)}\n"
	}
}

func (ed *Editor) updateStateFromWarnings() {
	ed.warningsCSS = ""
	if ed.WarningLines == nil {
		return
	}
	for key := range ed.WarningLines {
		ed.warningsCSS = ed.warningsCSS + ".shadow ol li:nth-child(" + key + ") {background: var(--warn-bgcolor)}\n"
	}
}

// Mount implements the vecty.Mounter interface.
func (ed *Editor) Mount() {
	obj := document.QuerySelector(".editor")
	if obj == nil {
		panic("Can't locate .editor")
	}
	ed.ta = &textarea.Textarea{obj}

	obj = document.QuerySelector(".shadow")
	if obj == nil {
		panic("Can't locate .shadow")
	}
	ed.sh = &Shadow{obj}
}

// Render implements the vecty.Component interface.
func (ed *Editor) Render() *vecty.HTML {
	ed.updateStateFromRanges()
	ed.updateStateFromWarnings()
	ed.updateStateFromErrors()
	util.Schedule(ed.afterRender)

	return elem.Div(
		vecty.Markup(
			vecty.Class("scroller"),
			event.MouseDown(ed.handleScrollerClick),
		),
		elem.TextArea(
			vecty.Markup(
				vecty.Class("editor"),
				vecty.MarkupIf(ed.HighlightingMode, vecty.Class("highlighted")),
				vecty.Property("autocapitalize", "off"),
				vecty.Attribute("autocomplete", "off"),
				vecty.Attribute("autocorrect", "off"),
				vecty.Property("autofocus", true),
				vecty.Property("spellcheck", false),
				vecty.Property("readonly", ed.ReadonlyMode),
				event.KeyDown(ed.handleKeyDown),
				event.KeyPress(ed.handleKeyPress),
				event.Select(ed.updateSelectionInfo),
				event.Input(ed.onChange),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("shadow"),
				vecty.UnsafeHTML(ed.highlighted),
				event.ContextMenu(ed.cancelEvent),
			),
		),
		elem.Style(
			vecty.Markup(
				vecty.UnsafeHTML(ed.selLinesCSS),
			),
		),
		elem.Style(
			vecty.Markup(
				vecty.UnsafeHTML(ed.warningsCSS),
			),
		),
		elem.Style(
			vecty.Markup(
				vecty.UnsafeHTML(ed.errorsCSS),
			),
		),
	)
}
