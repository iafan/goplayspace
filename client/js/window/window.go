package window

import "github.com/gopherjs/gopherjs/js"

// AddEventListener is a wrapper for window.addEventListener
func AddEventListener(params ...interface{}) {
	js.Global.Get("window").Call("addEventListener", params...)
}

// RemoveEventListener is a wrapper for window.removeEventListener
func RemoveEventListener(params ...interface{}) {
	js.Global.Get("window").Call("removeEventListener", params...)
}

// RequestAnimationFrame is a wrapper for window.requestAnimationFrame
func RequestAnimationFrame(callback interface{}) {
	js.Global.Get("window").Call("requestAnimationFrame", callback)
}
