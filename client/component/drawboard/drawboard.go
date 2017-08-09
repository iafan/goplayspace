package drawboard

import (
	"fmt"
	"math"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"

	"github.com/iafan/goplayspace/client/draw"
	"github.com/iafan/goplayspace/client/js/canvas"
	"github.com/iafan/goplayspace/client/js/document"
	"github.com/iafan/goplayspace/client/js/window"
	"github.com/iafan/goplayspace/client/util"
)

const (
	firstStepDelay = 500 * time.Millisecond
	stepDelay      = 500 * time.Millisecond
	// should be longer than `.say-bubble.animate`` CSS animation duration
	removeBubbleDelay = 5 * time.Second

	// when determining the scale of the board, how many cells should be visible
	// in each direction from the center of the board; the scale is calculated
	// based on the smallest dimension (width or height)
	stepsInEachDirection = 15

	walkFrameDistance = 2                    // distance in px along the path between animation frames
	walkFrames        = 5                    // total frames in the walk animation
	rotationFrame     = (walkFrames - 1) / 2 // middle frame index
	virtualWalkFrames = walkFrames*2 - 1     // we move back-forth between frames rather than cycle
	walkFrameSize     = 50

	boardLineWidth    = 1
	boardStrokeStyle  = "rgba(0, 0, 0, 0.05)"
	fifthStrokeStyle  = "rgba(0, 0, 0, 0.09)"
	centerStrokeStyle = "rgba(0, 0, 0, 0.16)"
)

// DrawBoard represents the drawing board with animation logic
type DrawBoard struct {
	vecty.Core
	canvas        *canvas.Canvas
	ctx           *canvas.CanvasRenderingContext2D
	gopher        *js.Object
	canvasWrapper *js.Object
	initialized   bool

	Actions draw.ActionList

	step int

	startTime  time.Time
	startX     float64
	startY     float64
	startAngle float64

	targetTime  time.Time
	targetX     float64
	targetY     float64
	targetAngle float64
	targetDist  float64

	accelerate bool
	tabDown    bool

	x, y  float64
	angle float64
	color string
	width float64

	w, h     float64
	stepSize float64
}

func (b *DrawBoard) getDOMNodes() {
	if b.canvas == nil {
		c := document.QuerySelector("canvas")
		if c != nil {
			b.canvas = &canvas.Canvas{c}
			b.ctx = b.canvas.GetContext2D()
		}
		b.gopher = document.QuerySelector(".gopher")
		b.canvasWrapper = document.QuerySelector(".canvas-wrapper")
	}
}

func (b *DrawBoard) renderBoardLines() {
	cX := b.w / 2
	cY := b.h / 2

	nX := int(cX/b.stepSize) + 1
	nY := int(cY/b.stepSize) + 1

	b.ctx.SetLineWidth(boardLineWidth)

	for x := -nX; x <= nX; x++ {
		b.ctx.SetStrokeStyle(boardStrokeStyle)
		if x%5 == 0 {
			b.ctx.SetStrokeStyle(fifthStrokeStyle)
		}
		if x == 0 {
			b.ctx.SetStrokeStyle(centerStrokeStyle)
		}
		b.ctx.BeginPath()
		b.ctx.MoveTo(cX+float64(x)*b.stepSize, 0)
		b.ctx.LineTo(cX+float64(x)*b.stepSize, b.h)
		b.ctx.Stroke()
	}

	for y := -nY; y <= nY; y++ {
		b.ctx.SetStrokeStyle(boardStrokeStyle)
		if y%5 == 0 {
			b.ctx.SetStrokeStyle(fifthStrokeStyle)
		}
		if y == 0 {
			b.ctx.SetStrokeStyle(centerStrokeStyle)
		}
		b.ctx.BeginPath()
		b.ctx.MoveTo(0, cY+float64(y)*b.stepSize)
		b.ctx.LineTo(b.w, cY+float64(y)*b.stepSize)
		b.ctx.Stroke()
	}
}

// addSpeechBubble shows the animated 'speech bubble'
// x, y are the center coordinates of the bubble in pixels
// relative to the center of the board
func (b *DrawBoard) addSpeechBubble(x, y float64, s string) {
	el := document.CreateElement("div")
	el.Set("className", "say-bubble")

	el.Set("innerHTML", s)
	b.canvasWrapper.Call("appendChild", el)

	// need to wait for the element to be rendered
	// in order to get offsetWidth / offsetHeight for centering
	util.Schedule(func() {
		elw := el.Get("offsetWidth").Float()
		elh := el.Get("offsetHeight").Float()

		cX := b.w / 2
		cY := b.h / 2

		// center the bubble around x, y point
		style := fmt.Sprintf(
			"left: %.0fpx; top: %.0fpx",
			cX+x-elw/2, cY+y-elh/2,
		)
		el.Call("setAttribute", "style", style)

		// start animation
		el.Set("className", "say-bubble animate")

		time.AfterFunc(removeBubbleDelay, func() {
			b.canvasWrapper.Call("removeChild", el)
		})
	})
}

func (b *DrawBoard) doSubStep(pos float64) {
	oldX := b.x
	oldY := b.y

	b.x = (b.targetX-b.startX)*pos + b.startX
	b.y = (b.targetY-b.startY)*pos + b.startY
	b.angle = (b.targetAngle-b.startAngle)*pos + b.startAngle

	//console.Log("x:", b.x, "y:", b.y, "angle:", b.angle)

	cX := b.w / 2
	cY := b.h / 2

	if b.color != "" {
		b.ctx.SetLineWidth(b.width)
		b.ctx.SetStrokeStyle(b.color)
		b.ctx.BeginPath()
		b.ctx.MoveTo(cX+oldX, cY+oldY)
		b.ctx.LineTo(cX+b.x, cY+b.y)
		b.ctx.Stroke()
	}

	frame := int(b.targetDist*pos/walkFrameDistance) % virtualWalkFrames

	// offset frame number by rotationFrame index
	frame = (frame + rotationFrame) % virtualWalkFrames

	if frame > walkFrames-1 {
		frame = virtualWalkFrames - frame
	}

	bgPos := -frame * walkFrameSize

	style := fmt.Sprintf(
		"transform: translateX(%.2fpx) translateY(%.2fpx) rotate(%.2fdeg); "+
			"background-position-x: %dpx;",
		b.x, b.y, b.angle,
		bgPos,
	)

	b.gopher.Call("setAttribute", "style", style)
}

func (b *DrawBoard) doStep() {
	t := time.Now()

	if b.targetTime.IsZero() || b.targetTime.Sub(t) <= 0 || b.accelerate {
		b.doSubStep(1)

		// new step
		b.step = b.step + 1

		if b.step == len(b.Actions) {
			//console.Log("Animation stopped")
			return
		}

		//console.Log("Step:", b.step)

		b.startTime = t
		b.targetTime = t

		b.startX = b.x
		b.startY = b.y
		b.startAngle = b.angle

		a := b.Actions[b.step]

		delay := stepDelay

		switch a.Kind {
		case draw.Step:
			b.targetTime = t.Add(time.Duration(float64(delay) * a.FVal))

			rad := (-90 + b.angle) * 2 * math.Pi / 360
			b.targetX = b.startX + math.Cos(rad)*b.stepSize*a.FVal
			b.targetY = b.startY + math.Sin(rad)*b.stepSize*a.FVal

			// stop accelerating only after the 'Step' event; accelerate through others
			if b.tabDown {
				b.accelerate = false
			}

		case draw.Left:
			b.targetTime = t.Add(delay)
			b.targetAngle = b.startAngle - a.FVal // sign inverted to match clock-wise CSS rotation
		case draw.Right:
			b.targetTime = t.Add(delay)
			b.targetAngle = b.startAngle + a.FVal // sign inverted to match clock-wise CSS rotation
		case draw.Color:
			b.color = a.SVal
			util.Schedule(b.doStep)
			return
		case draw.Width:
			b.width = a.FVal
			util.Schedule(b.doStep)
			return
		case draw.Say:
			b.addSpeechBubble(b.x, b.y, a.SVal)
			util.Schedule(b.doStep)
			return
		}

		b.targetDist = math.Sqrt(
			math.Pow(b.targetX-b.startX, 2) + math.Pow(b.targetY-b.startY, 2),
		)
	}

	// calculate current position
	total := b.targetTime.Sub(b.startTime)  // total duration
	passed := t.Sub(b.startTime)            // passed duration
	rel := float64(passed) / float64(total) // passed [0..1]
	b.doSubStep(rel)

	window.RequestAnimationFrame(b.doStep)
}

func (b *DrawBoard) animate() {
	b.getDOMNodes()

	// set defaults
	b.width = 2

	b.step = -1
	//console.Log("Animation started")
	time.AfterFunc(firstStepDelay, b.doStep)
}

func (b *DrawBoard) onRendered() {
	b.getDOMNodes()

	time.AfterFunc(100*time.Millisecond, func() {
		document.QuerySelector(".canvas-lightbox").Call("focus")
	})

	if !b.initialized {
		b.initialized = true
		window.AddEventListener("resize", b.onResize)
		b.onResize()

		// start the animation
		b.animate()
	}
}

func (b *DrawBoard) handleKeyDown(e *vecty.Event) {
	switch e.Get("key").String() {
	case "Shift":
		b.accelerate = true
	case "Tab":
		e.Call("preventDefault")
		if b.tabDown {
			return
		}
		b.accelerate = true
		b.tabDown = true
	default:
		//console.Log(e.Get("key").String())
	}
}

func (b *DrawBoard) handleKeyUp(e *vecty.Event) {
	switch e.Get("key").String() {
	case "Shift":
		b.accelerate = false
	case "Tab":
		e.Call("preventDefault")
		if !b.tabDown {
			return
		}
		b.tabDown = false
	}
}

func (b *DrawBoard) onResize() {
	b.w, b.h = b.canvas.GetNodeSize()
	min := b.w
	if b.h < min {
		min = b.h
	}
	b.stepSize = min / (stepsInEachDirection*2 + 1) // "+1" to add 0.5 steps around
	b.canvas.SetSize(b.w, b.h)
	b.renderBoardLines()
	b.gopher.Call("setAttribute", "style", "")
}

// SkipRender implements the vecty.Component interface.
func (b *DrawBoard) SkipRender(prev vecty.Component) bool {
	return true
}

// Render implements the vecty.Component interface.
func (b *DrawBoard) Render() *vecty.HTML {
	util.Schedule(b.onRendered)

	return elem.Div(
		vecty.ClassMap{"canvas-lightbox": true},
		vecty.Attribute("tabindex", 0),
		elem.Div(
			vecty.ClassMap{"canvas-wrapper": true},
			elem.Canvas(),
			elem.Div(
				vecty.ClassMap{
					"gopher": true,
				},
			),
			elem.Div(
				vecty.ClassMap{"statusbar-wrapper": true},
				elem.Div(
					vecty.ClassMap{
						"statusbar": true,
					},
					vecty.UnsafeHTML("<kbd>Tab</kbd> or hold <kbd>Shift</kbd> to accelerate, <kbd>Esc</kbd> to close"),
				),
			),
		),
		event.KeyDown(b.handleKeyDown),
		event.KeyUp(b.handleKeyUp),
	)
}
