package localstorage

import (
	"strconv"

	"github.com/gopherjs/gopherjs/js"
)

// Get wraps localStorage.getItem and returns a string value
func Get(key string, def string) string {
	v := js.Global.Get("localStorage").Call("getItem", key)
	if v == nil {
		return def
	}
	return v.String()
}

// Set is a wrapper for localStorage.setItem
func Set(key string, value interface{}) {
	js.Global.Get("localStorage").Call("setItem", key, value)
}

// GetInt wraps localStorage.getItem and returns an integer value
func GetInt(key string, def int) int {
	v := js.Global.Get("localStorage").Call("getItem", key)
	if v == nil {
		return def
	}
	return v.Int()
}

// GetBool wraps localStorage.getItem and returns a boolean value
func GetBool(key string, def bool) bool {
	v := js.Global.Get("localStorage").Call("getItem", key)
	if v == nil {
		return def
	}
	b, err := strconv.ParseBool(v.String())
	if err != nil {
		return def
	}
	return b
}
