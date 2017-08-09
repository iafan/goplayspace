package document

import "github.com/gopherjs/gopherjs/js"

// Body is a wrapper for document.body
func Body() *js.Object {
	return js.Global.Get("document").Get("body")
}

// CreateElement is a wrapper for document.createElement
func CreateElement(name string) *js.Object {
	return js.Global.Get("document").Call("createElement", name)
}

// QuerySelector is a wrapper for document.querySelector
func QuerySelector(sel string) *js.Object {
	return js.Global.Get("document").Call("querySelector", sel)
}

// AddEventListener is a wrapper for document.addEventListener
func AddEventListener(params ...interface{}) {
	js.Global.Get("document").Call("addEventListener", params...)
}

// RemoveEventListener is a wrapper for document.removeEventListener
func RemoveEventListener(params ...interface{}) {
	js.Global.Get("document").Call("removeEventListener", params...)
}
