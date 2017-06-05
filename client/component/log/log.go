package log

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/iafan/goplayspace/client/api"
)

// Log contains the logic behind the log panel
// exposed on the applicaiton page under '.log' class
type Log struct {
	vecty.Core
	node *js.Object

	Error  string
	Events []*api.CompileEvent
	HasRun bool
}

func (l *Log) getDOMNode() *js.Object {
	if l.node == nil {
		l.node = js.Global.Get("document").Call("querySelector", ".log")
	}
	return l.node
}

func (l *Log) getEvents() []vecty.MarkupOrComponentOrHTML {
	if len(l.Events) == 0 {
		return nil
	}
	out := make([]vecty.MarkupOrComponentOrHTML, len(l.Events)+1)

	for i, evt := range l.Events {
		out[i] = elem.Div(
			vecty.ClassMap{evt.Kind: true},
			//vecty.Text("["+strconv.Itoa(evt.Delay)+"] "+evt.Message), // FIXME: show timings only when necessary
			vecty.Text(evt.Message),
		)
	}

	final := ""
	if l.HasRun {
		final = "Program exited."
		if len(l.Events) == 0 {
			final = "Program exited producing no output."
		}
	}

	out[len(out)-1] = elem.Div(
		vecty.ClassMap{"final": true},
		vecty.Text(final),
	)
	return out
}

func (l *Log) getStatusText() string {
	if l.Error != "" {
		return l.Error
	}
	return "Syntax OK"
}

// ScrollToBottom scrolls log area to the bottom
func (l *Log) ScrollToBottom() {
	if l.getDOMNode() == nil {
		return
	}
	l.node.Set("scrollTop", l.node.Get("scrollHeight").Int())
}

// Render implements the vecty.Component interface.
func (l *Log) Render() *vecty.HTML {
	return elem.Div(
		vecty.ClassMap{"log": true},
		elem.Div(l.getEvents()...),
		elem.Div(
			vecty.ClassMap{
				"status": true,
				"error":  l.Error != "",
			},
			vecty.Text(l.getStatusText()),
		),
	)
}
