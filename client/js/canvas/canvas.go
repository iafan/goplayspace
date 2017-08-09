package canvas

import (
	"github.com/gopherjs/gopherjs/js"
)

type Canvas struct {
	*js.Object
}

// GetContext2D returns a '2d' context for the canvas
func (c *Canvas) GetContext2D() *CanvasRenderingContext2D {
	return &CanvasRenderingContext2D{c.Call("getContext", "2d")}
}

func (c *Canvas) GetNodeSize() (w, h float64) {
	w = c.Get("offsetWidth").Float()
	h = c.Get("offsetHeight").Float()
	return
}

func (c *Canvas) SetSize(w, h float64) {
	c.Set("width", w)
	c.Set("height", h)
}
