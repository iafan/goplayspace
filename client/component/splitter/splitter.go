package splitter

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/iafan/goplayspace/client/js/document"
)

// Splitter type definitions
const (
	LeftPane = iota
	RightPane
	TopPane
	BottomPane
)

// Splitter contains the column/row
// (aka vertical/horizontal) splitter component
type Splitter struct {
	vecty.Core

	Selector         string
	OppositeSelector string
	MinSizePercent   float64
	Type             int
	InvertSizeDelta  bool

	isMoving      bool
	node          *js.Object
	movingPane    *js.Object
	oppositePane  *js.Object
	origSize      int
	origScreenPos int
	parentSize    int
	styleParam    string
	screenCoord   string
	pos           float64
	prevPos       float64
}

func (s *Splitter) onMouseDown(e *vecty.Event) {
	e.Call("preventDefault")

	s.node = e.Get("target")
	s.node.Get("classList").Call("add", "moving")
	js.Global.Get("document").Get("body").Get("classList").Call("add", "moving")

	s.movingPane = document.QuerySelector(s.Selector)
	if s.movingPane == nil {
		panic("Can't find node using '" + s.Selector + "' query selector")
	}

	if s.OppositeSelector != "" {
		s.oppositePane = document.QuerySelector(s.OppositeSelector)
		if s.oppositePane == nil {
			panic("Can't find node using '" + s.OppositeSelector + "' query selector")
		}
	}

	s.styleParam = "height"
	s.screenCoord = "screenY"
	sizeParam := "offsetHeight"
	if s.Type == LeftPane || s.Type == RightPane {
		s.styleParam = "width"
		s.screenCoord = "screenX"
		sizeParam = "offsetWidth"
	}

	s.origSize = s.movingPane.Get(sizeParam).Int()
	s.origScreenPos = e.Get(s.screenCoord).Int()
	s.parentSize = s.movingPane.Get("parentNode").Get(sizeParam).Int()

	s.isMoving = true

	document.AddEventListener("mousemove", s.onDocumentMouseMove)
	document.AddEventListener("mouseup", s.onDocumentMouseUp)
}

func (s *Splitter) updatePos() {
	if s.pos == s.prevPos {
		return
	}
	s.prevPos = s.pos

	s.movingPane.Call("setAttribute", "style", s.styleParam+":"+strconv.FormatFloat(s.pos, 'f', 4, 32)+"%")

	if s.oppositePane != nil {
		s.oppositePane.Call("setAttribute", "style", s.styleParam+":"+strconv.FormatFloat(100-s.pos, 'f', 4, 32)+"%")
	}

	js.Global.Call("requestAnimationFrame", s.updatePos)
}

func (s *Splitter) onDocumentMouseMove(e *vecty.Event) {
	e.Call("preventDefault")

	screenPos := e.Get(s.screenCoord).Int()

	coeff := 1
	if s.Type == BottomPane || s.Type == RightPane {
		coeff = -1
	}

	size := s.origSize + coeff*(screenPos-s.origScreenPos)
	f := 100 * float64(size) / float64(s.parentSize)

	if f < s.MinSizePercent {
		f = s.MinSizePercent
	}
	if f > 100-s.MinSizePercent {
		f = 100 - s.MinSizePercent
	}
	s.pos = f

	js.Global.Call("requestAnimationFrame", s.updatePos)
}

func (s *Splitter) onDocumentMouseUp(e *vecty.Event) {
	e.Call("preventDefault")

	document.RemoveEventListener("mousemove", s.onDocumentMouseMove)
	document.RemoveEventListener("mouseup", s.onDocumentMouseUp)
	s.node.Get("classList").Call("remove", "moving")
	js.Global.Get("document").Get("body").Get("classList").Call("remove", "moving")

	s.isMoving = false
}

// Render implements the vecty.Component interface.
func (s *Splitter) Render() *vecty.HTML {
	return elem.Div(
		vecty.ClassMap{
			"splitter": true,
			"col":      s.Type == LeftPane || s.Type == RightPane,
			"row":      s.Type == TopPane || s.Type == BottomPane,
		},
		event.MouseDown(s.onMouseDown),
	)
}
