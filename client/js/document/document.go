package document

import "github.com/gopherjs/gopherjs/js"

// QuerySelector is a wrapper for document.querySelector
func QuerySelector(sel string) *js.Object {
	return js.Global.Get("document").Call("querySelector", sel)
}
