package localstorage

import "github.com/gopherjs/gopherjs/js"

// Get is a wrapper for localStorage.getItem
func Get(key string) string {
	v := js.Global.Get("localStorage").Call("getItem", key)
	if v == nil {
		return ""
	}
	return v.String()
}

// Set is a wrapper for localStorage.setItem
func Set(key, value string) {
	js.Global.Get("localStorage").Call("setItem", key, value)
}
