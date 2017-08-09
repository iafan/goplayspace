package canvas

import (
	"github.com/gopherjs/gopherjs/js"
)

type CanvasRenderingContext2D struct {
	*js.Object
}

// Properties

func (ctx *CanvasRenderingContext2D) SetFillStyle(style interface{}) {
	ctx.Set("fillStyle", style)
}

func (ctx *CanvasRenderingContext2D) SetLineWidth(w float64) {
	ctx.Set("lineWidth", w)
}

func (ctx *CanvasRenderingContext2D) SetStrokeStyle(style interface{}) {
	ctx.Set("strokeStyle", style)
}

// Methods

func (ctx *CanvasRenderingContext2D) BeginPath() {
	ctx.Call("beginPath")
}

func (ctx *CanvasRenderingContext2D) ClearRect(x, y, w, h float64) {
	ctx.Call("clearRect", x, y, w, h)
}

func (ctx *CanvasRenderingContext2D) FillRect(x, y, w, h float64) {
	ctx.Call("fillRect", x, y, w, h)
}

func (ctx *CanvasRenderingContext2D) Translate(x, y float64) {
	ctx.Call("translate", x, y)
}

func (ctx *CanvasRenderingContext2D) MoveTo(x, y float64) {
	ctx.Call("moveTo", x, y)
}

func (ctx *CanvasRenderingContext2D) LineTo(x, y float64) {
	ctx.Call("lineTo", x, y)
}

func (ctx *CanvasRenderingContext2D) Stroke() {
	ctx.Call("stroke")
}
