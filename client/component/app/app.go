package app

import (
	"encoding/json"
	"go/format"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"strings"
	"time"

	"honnef.co/go/js/xhr"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/iafan/syntaxhighlight"

	"github.com/iafan/goplayspace/client/api"
	"github.com/iafan/goplayspace/client/component/drawboard"
	"github.com/iafan/goplayspace/client/component/editor"
	"github.com/iafan/goplayspace/client/component/editor/undo"
	"github.com/iafan/goplayspace/client/component/help"
	"github.com/iafan/goplayspace/client/component/log"
	"github.com/iafan/goplayspace/client/component/settings"
	"github.com/iafan/goplayspace/client/component/splitter"
	"github.com/iafan/goplayspace/client/draw"
	"github.com/iafan/goplayspace/client/hash"
	"github.com/iafan/goplayspace/client/js/console"
	"github.com/iafan/goplayspace/client/js/localstorage"
	"github.com/iafan/goplayspace/client/js/window"
	"github.com/iafan/goplayspace/client/ranges"
	"github.com/iafan/goplayspace/client/util"
)

const maxUndoStackSize uint = 50

const idDrawPage = "draw"

// Application implements the main application view
type Application struct {
	vecty.Core

	editor *editor.Editor
	log    *log.Log

	Input   string
	Topic   string
	Imports map[string]string

	// Settings
	Theme            string
	TabWidth         int
	FontWeight       string
	UseWebfont       bool
	HighlightingMode bool

	Hash      *hash.Hash
	snippetID string

	modifierKey          string
	isLoading            bool
	isCompiling          bool
	isSharing            bool
	isDrawingMode        bool
	hasCompilationErrors bool
	needRender           bool
	showSettings         bool
	showDrawHelp         bool

	// Log properties
	hasRun bool
	err    string
	events []*api.CompileEvent

	// Draw mode properties
	actions draw.ActionList

	// Editor properties
	warningLines map[string]bool
	errorLines   map[string]bool
	undoStack    *undo.Stack
	changeTimer  *time.Timer
}

func (a *Application) rerenderIfNeeded() {
	if !a.needRender {
		return
	}
	a.needRender = false
	vecty.Rerender(a)
}

func (a *Application) wantRerender(reason string) {
	//console.Log("want rerender:", reason)
	a.needRender = true
	util.Schedule(a.rerenderIfNeeded)
}

func (a *Application) onEditorTopicChange(topic string) {
	a.Topic = topic
	a.wantRerender("onEditorTopicChange")
}

func (a *Application) onEditorKeyDown(e *vecty.Event) {
	ctrlDown := e.Get("ctrlKey").Bool()
	metaDown := e.Get("metaKey").Bool()

	switch e.Get("keyCode").Int() {
	case 83: // S
		if ctrlDown || metaDown { // Ctrl+S or Cmd+S
			e.Call("preventDefault")
			a.doFormat()
		}
	case 13: // Enter
		if ctrlDown || metaDown { // Ctrl+Enter or Cmd+Enter
			e.Call("preventDefault")
			if a.err != "" || a.isCompiling {
				return
			}
			a.doRun()
		}
	}
}

var compileErrorLineExtractorR = regexp.MustCompile(`\/main\.go:(\d+):\s`)
var fmtErrorLineExtractorR = regexp.MustCompile(`(?m)^(\d+):(\d+):\s`)

var domMonitorInterval = 5 * time.Millisecond

func (a *Application) onLineSelChange(state string) {
	if a.isLoading || a.Hash.Ranges == state {
		return
	}
	a.Hash.SetRanges(state)
	a.wantRerender("onLineSelChange")
}

func (a *Application) runButtonClick(e *vecty.Event) {
	a.doRun()
}

func (a *Application) doRun() {
	a.isCompiling = true
	//a.doFormat()
	go a.doRunAsync()
}

func (a *Application) doRunAsync() {
	defer a.doRunAsyncComplete()

	a.hasRun = true

	bodyBytes, err := xhr.Send("POST", "/compile", []byte(a.Input))
	if err != nil {
		a.err = err.Error()
		return
	}

	compileResponse := api.CompileResponse{}

	err = json.Unmarshal(bodyBytes, &compileResponse)
	if err != nil {
		a.err = err.Error()
		return
	}

	a.err = compileResponse.Errors
	a.events = compileResponse.Events
	a.hasCompilationErrors = a.err != ""

	if compileResponse.Body != nil {
		a.setEditorText(*compileResponse.Body)
	}

	// extract line numbers from compilation error message

	if matches := compileErrorLineExtractorR.FindAllStringSubmatch(compileResponse.Errors, -1); matches != nil {
		a.errorLines = make(map[string]bool)
		for _, m := range matches {
			a.errorLines[m[1]] = true
		}
	}

	// parse gopher commands
	if !a.hasCompilationErrors {
		output := make([]string, len(a.events))
		for i := range a.events {
			output[i] = a.events[i].Message
		}
		a.actions = draw.ParseString(strings.Join(output, "\n"))
		a.isDrawingMode = len(a.actions) > 0
	}
}

func (a *Application) doRunAsyncComplete() {
	a.isCompiling = false
	a.wantRerender("doRunAsyncComplete")
	util.Schedule(func() { a.log.ScrollToBottom() })
}

func (a *Application) shareButtonClick(e *vecty.Event) {
	a.doShare()
}

func (a *Application) doShare() {
	a.isSharing = true
	a.doFormat()
	go a.doShareAsync()
}

func (a *Application) doShareAsync() {
	defer a.doShareAsyncComplete()

	bodyBytes, err := xhr.Send("POST", "/share", []byte(a.Input))
	if err != nil {
		a.err = err.Error()
		return
	}

	a.snippetID = string(bodyBytes) // already 'loaded'
	a.Hash.SetID(a.snippetID)
}

func (a *Application) doShareAsyncComplete() {
	a.isSharing = false
	a.wantRerender("doShareAsyncComplete")
}

func (a *Application) updateStateFromHash(h *hash.Hash) (canLoad bool) {
	if h.ID == idDrawPage {
		a.showDrawHelp = true
		return false
	}

	return true
}

func (a *Application) onHashChange(h *hash.Hash) {
	defer a.wantRerender("onHashChange")

	if a.isLoading || h.ID == "" {
		return
	}

	if a.updateStateFromHash(h) {
		a.doLoad(h.ID)
	}
}

func (a *Application) doLoad(id string) {
	if id == a.snippetID || id == "" {
		return
	}
	a.isLoading = true
	go a.doLoadAsync(id)
}

func (a *Application) doLoadAsync(id string) {
	defer a.doLoadAsyncComplete(id)

	req := xhr.NewRequest("GET", "/load?"+id)
	err := req.Send(nil)
	//bodyBytes, err := xhr.Send("GET", "/load?"+id, nil)
	if err != nil {
		a.err = err.Error()
		return
	}
	if req.Status != 200 {
		a.err = req.ResponseText
		return
	}

	a.setEditorText(req.ResponseText)
	// setting new text will cause OnChange event,
	// and hash will be reset; so update it afterwards
	a.Hash.ID = id
}

func (a *Application) doLoadAsyncComplete(id string) {
	a.isLoading = false
	a.snippetID = id
	a.wantRerender("doLoadAsyncComplete")
}

func (a *Application) formatButtonClick(e *vecty.Event) {
	a.doFormat()
}

func (a *Application) format(text string) (string, error) {
	if a.Input == "" {
		return "", nil
	}

	//console.Time("format")
	bytes, err := format.Source([]byte(a.Input))
	//console.TimeEnd("format")

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (a *Application) doFormat() {
	defer util.Schedule(a.editor.Focus)
	a.wantRerender("doFormat")

	text, err := a.format(a.Input)
	if err != nil {
		a.err = err.Error()
	} else {
		a.err = ""
		a.setEditorText(text)
	}
}

func (a *Application) setEditorText(text string) {
	if a.Input == text {
		return
	}
	a.Input = text
	a.parseAndReportErrors(text)
	a.editor.SetText(text)
	util.Schedule(a.editor.Focus)
}

func (a *Application) setEditorState(text string, selStart, selEnd int) {
	if a.Input == text {
		return
	}
	a.Input = text
	a.parseAndReportErrors(text)
	a.editor.SetState(text, selStart, selEnd)
	util.Schedule(a.editor.Focus)
}

func (a *Application) onEditorValueChange(text string) {
	if a.Input == text {
		return
	}
	a.Input = text
	a.parseAndReportErrors(text)
	a.Hash.Reset()
	a.wantRerender("onEditorValueChange")
}

func (a *Application) parseAndReportErrors(text string) {
	a.err = ""
	a.warningLines = nil
	a.errorLines = nil
	a.hasCompilationErrors = false

	if text == "" {
		a.setEditorState(blankTemplate, blankTemplatePos, blankTemplatePos)
	}

	// parse source code to get list of imports and parsing error, if any;
	// note that we don't clear the list of imports since we want to
	// keep the previously known good mapping even if there are parsing errors

	fset := token.NewFileSet()
	//console.Time("parse")
	f, err := parser.ParseFile(fset, "", a.Input, parser.AllErrors)
	//console.TimeEnd("parse")

	a.Imports = make(map[string]string)
	if f != nil {
		for _, imp := range f.Imports {
			var name string
			path := strings.Trim(imp.Path.Value, `"`)
			if imp.Name != nil {
				name = imp.Name.Name
			} else {
				name = path
				if i := strings.LastIndex(path, "/"); i >= -1 {
					name = path[i+1:]
				}
			}

			// FIXME: should we somehow deal with '.' and '_' import names?

			if name != "." && name != "_" {
				a.Imports[name] = path // short package name
			}
			if path != "." && path != "_" && path != name {
				a.Imports[path] = path // full package name
			}
		}
	}

	if err != nil {
		a.err = err.Error()

		// extract line numbers from parser error message

		if matches := fmtErrorLineExtractorR.FindAllStringSubmatch(a.err, -1); matches != nil {
			a.warningLines = make(map[string]bool)
			for _, m := range matches {
				a.warningLines[m[1]] = true
			}
		}
	}
}

// highlight function is used to highlight source code in the editor
func (a *Application) highlight(text string) string {
	//console.Time("highlight")
	//defer console.TimeEnd("highlight")
	hbytes, err := syntaxhighlight.AsHTML([]byte(text), syntaxhighlight.OrderedList())
	if err != nil {
		console.Log("Highlight error:", err)
		a.err = err.Error()
		return ""
	}
	return string(hbytes)
}

func (a *Application) getGlobalState() (out string) {
	out = "ok"
	if a.err != "" {
		out = "warning"
		if a.hasCompilationErrors {
			out = "error"
		}
	}
	return
}

func (a *Application) getFiraFontCSS(weight, suffix string) string {
	return `@font-face {
	font-family: 'Fira Code';
	font-weight: ` + weight + `;
	src: url('https://raw.githubusercontent.com/tonsky/FiraCode/master/distr/woff2/FiraCode-` + suffix + `.woff2') format('woff2');
}`
}

func (a *Application) getOverrideCSS() (out string) {
	if a.UseWebfont {
		if a.FontWeight == "normal" {
			out += a.getFiraFontCSS(a.FontWeight, "Regular")
		} else {
			out += a.getFiraFontCSS(a.FontWeight, "Light")
		}
	}

	out += `.editor, .shadow, .log {
	font-weight: ` + a.FontWeight + `;
}`
	return out
}

func (a *Application) updateTheme(val string) {
	a.Theme = val
	localstorage.Set("theme", val)
	a.wantRerender("updateTheme")
}

func (a *Application) updateTabWidth(val int) {
	a.TabWidth = val
	localstorage.Set("tab-width", val)
	a.wantRerender("updateTabWidth")
}

func (a *Application) updateFontWeight(val string) {
	a.FontWeight = val
	localstorage.Set("font-weight", val)
	a.wantRerender("updateFontWeight")
}

func (a *Application) updateUseWebfont(val bool) {
	a.UseWebfont = val
	localstorage.Set("use-webfont", val)
	a.wantRerender("updateUseWebfont")
}

func (a *Application) updateHighlighting(on bool) {
	a.HighlightingMode = on
	localstorage.Set("highlighting", on)
	a.editor.Highlight(on)
	a.wantRerender("updateHighlighting")
}

func (a *Application) onSettingsChange(d *settings.Dialog) {
	if d.Theme != a.Theme {
		a.updateTheme(d.Theme)
	}

	if d.TabWidth != a.TabWidth {
		a.updateTabWidth(d.TabWidth)
	}

	if d.FontWeight != a.FontWeight {
		a.updateFontWeight(d.FontWeight)
	}

	if d.UseWebfont != a.UseWebfont {
		a.updateUseWebfont(d.UseWebfont)
	}

	if d.HighlightingMode != a.HighlightingMode {
		a.updateHighlighting(d.HighlightingMode)
	}
}

func (a *Application) formatShortcutPressed(e interface{}) {
	a.err = "formatShortcutPressed()"
	a.wantRerender("formatShortcutPressed")
}

// WaitForPageLoaded waits till page is loaded (editor is ready)
// and then calls onPageLoaded()
func (a *Application) WaitForPageLoaded() {
	if a.editor.IsReady() {
		a.onPageLoaded()
	} else {
		time.AfterFunc(domMonitorInterval, a.WaitForPageLoaded)
	}
}

func (a *Application) onPageLoaded() {
	switch a.Hash.ID {
	case "":
		a.setEditorState(initialCode, initialCaretPos, initialCaretPos)
	case idDrawPage:
		a.setEditorState(initialDrawCode, initialDrawCaretPos, initialDrawCaretPos)
		fallthrough
	default:
		a.onHashChange(a.Hash)
	}
	window.AddEventListener("resize", a.onResize)
}

func (a *Application) onResize() {
	a.editor.ResizeTextarea()
}

func (a *Application) settingsButtonClick(e *vecty.Event) {
	a.showSettings = !a.showSettings
	a.wantRerender("settingsButtonClick")
}

func (a *Application) handleKeyDown(e *vecty.Event) {
	switch e.Get("key").String() {
	case "Escape":
		if a.isDrawingMode {
			a.isDrawingMode = false
			a.wantRerender("isDrawingMode switched off")
			util.Schedule(a.editor.Focus)
		}
	default:
		//console.Log(e.Get("key").String())
	}
}

// Render renders the application
func (a *Application) Render() *vecty.HTML {
	//console.Time("app:render")
	//defer console.TimeEnd("app:render")

	if a.Hash == nil {
		a.Hash = hash.New(a.onHashChange)
		a.updateStateFromHash(a.Hash)
	}

	if a.undoStack == nil {
		a.undoStack = undo.NewStack(maxUndoStackSize)
	}

	if a.modifierKey == "" {
		a.modifierKey = "Ctrl"
		if util.IsMacOS() {
			a.modifierKey = "⌘"
		}
	}

	a.editor = &editor.Editor{
		Highlighter:      a.highlight,
		OnChange:         a.onEditorValueChange,
		OnLineSelChange:  a.onLineSelChange,
		OnTopicChange:    a.onEditorTopicChange,
		OnKeyDown:        a.onEditorKeyDown,
		WarningLines:     a.warningLines,
		ErrorLines:       a.errorLines,
		Range:            ranges.New(a.Hash.Ranges),
		ChangeTimer:      &a.changeTimer,
		UndoStack:        a.undoStack,
		HighlightingMode: a.HighlightingMode,
		ReadonlyMode:     a.isDrawingMode,
	}

	a.log = &log.Log{
		Error:  a.err,
		Events: a.events,
		HasRun: a.hasRun,
	}

	tabWidthClass := "tabwidth-" + strconv.Itoa(a.TabWidth)

	return elem.Body(
		vecty.ClassMap{
			"safari":           util.IsSafari(),
			"ios":              util.IsIOS(),
			"drawingmode":      a.isDrawingMode,
			a.Theme:            true,
			tabWidthClass:      true,
			a.getGlobalState(): true,
		},
		elem.Div(
			vecty.ClassMap{"header": true},
			elem.Div(
				vecty.ClassMap{"logo": true},
			),
			elem.Div(
				vecty.ClassMap{"menu": true},
				elem.Span(
					vecty.ClassMap{"title": true},
					vecty.UnsafeHTML("The Go<br/>Play Space"),
				),
				elem.Button(
					vecty.UnsafeHTML("Run <cmd>"+a.modifierKey+"+↵</cmd>"),
					vecty.Property("disabled", a.err != "" || a.isCompiling),
					event.Click(a.runButtonClick),
				),
				elem.Button(
					vecty.UnsafeHTML("Format <cmd>"+a.modifierKey+"+S</cmd>"),
					vecty.Property("disabled", a.err != ""),
					event.Click(a.formatButtonClick),
				),
				elem.Button(
					vecty.UnsafeHTML("Share"),
					vecty.Property("disabled", a.isSharing || a.Hash.ID != ""),
					event.Click(a.shareButtonClick),
				),
			),
			elem.Div(
				vecty.ClassMap{"settings": true},
				elem.Button(
					vecty.UnsafeHTML("Settings"),
					event.Click(a.settingsButtonClick),
				),
			),
		),
		elem.Div(
			vecty.ClassMap{"body-wrapper": true},
			elem.Div(
				vecty.ClassMap{"content-wrapper": true},
				a.editor,
				elem.Div(
					vecty.ClassMap{"help-wrapper": true},
					vecty.If(a.Topic == "" && !a.showDrawHelp, elem.Div(
						vecty.ClassMap{"help": true},
						vecty.UnsafeHTML(helpHTML),
					)),
					vecty.If(a.Topic == "" && a.showDrawHelp, elem.Div(
						vecty.ClassMap{"help": true},
						vecty.UnsafeHTML(drawHelpHTML),
					)),
					vecty.If(a.Topic != "", &help.Browser{
						Imports: a.Imports,
						Topic:   a.Topic,
					}),
					&splitter.Splitter{
						Selector:         ".help-wrapper",
						OppositeSelector: ".scroller",
						Type:             splitter.RightPane,
						MinSizePercent:   2,
						OnChange:         a.onResize,
					},
				),
			),
			elem.Div(
				vecty.ClassMap{"log-wrapper": true},
				a.log,
				&splitter.Splitter{
					Selector:         ".log-wrapper",
					OppositeSelector: ".content-wrapper",
					Type:             splitter.BottomPane,
					MinSizePercent:   2,
				},
			),
		),
		vecty.If(a.showSettings, &settings.Dialog{
			Theme:            a.Theme,
			TabWidth:         a.TabWidth,
			FontWeight:       a.FontWeight,
			UseWebfont:       a.UseWebfont,
			HighlightingMode: a.HighlightingMode,
			OnChange:         a.onSettingsChange,
		}),
		vecty.If(a.isDrawingMode, &drawboard.DrawBoard{
			Actions: a.actions,
		}),
		elem.Style(
			vecty.UnsafeHTML(a.getOverrideCSS()),
		),
		event.KeyDown(a.handleKeyDown),
	)
}
