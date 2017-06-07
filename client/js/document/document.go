package document

import "github.com/gopherjs/gopherjs/js"

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
