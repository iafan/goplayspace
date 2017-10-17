package log

import (
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/iafan/goplayspace/client/api"
	"github.com/iafan/goplayspace/client/js/document"
)

var hour = 60 * time.Minute
var day = 24 * hour

// Log contains the logic behind the log panel
// exposed on the application page under '.log' class
type Log struct {
	vecty.Core
	node *js.Object

	Error  string              `vecty:"prop"`
	Events []*api.CompileEvent `vecty:"prop"`
	HasRun bool                `vecty:"prop"`
}

func (l *Log) getEvents() []vecty.MarkupOrChild {
	if len(l.Events) == 0 {
		return nil
	}
	out := make([]vecty.MarkupOrChild, len(l.Events)+1)

	var totalDelay time.Duration
	for _, evt := range l.Events {
		totalDelay += evt.Delay
	}

	format := "T+15:04:05"
	totalDays := int(totalDelay / day)

	if totalDelay < hour {
		format = "T+04:05.000"
	}
	if totalDelay < time.Minute {
		format = "T+05.000000"
	}
	if totalDelay < time.Millisecond {
		format = "T+05.000000000"
	}

	deltaTime := time.Time{}
	var deltaDuration time.Duration
	for i, evt := range l.Events {
		deltaTime = deltaTime.Add(evt.Delay)
		deltaDuration += evt.Delay
		text := deltaTime.Format(format)
		if totalDays > 0 {
			text = "D+" + strconv.Itoa(int(deltaDuration/day)) + " " + text
		}
		out[i] = elem.Div(
			vecty.Markup(
				vecty.Class(evt.Kind),
			),
			vecty.If(totalDelay > 0, elem.Span(
				vecty.Markup(
					vecty.Class("time"),
				),
				vecty.Text(text),
			)),
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
		vecty.Markup(
			vecty.Class("final"),
		),
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
	if l.node == nil {
		return
	}
	l.node.Set("scrollTop", l.node.Get("scrollHeight").Int())
}

// Mount implements the vecty.Mounter interface.
func (l *Log) Mount() {
	l.node = document.QuerySelector(".log")
	if l.node == nil {
		panic("Can't locate .log")
	}
}

// Render implements the vecty.Component interface.
func (l *Log) Render() *vecty.HTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("log"),
		),
		elem.Div(l.getEvents()...),
		elem.Div(
			vecty.Markup(
				vecty.Class("status"),
				vecty.MarkupIf(l.Error != "", vecty.Class("error")),
			),
			vecty.Text(l.getStatusText()),
		),
	)
}
