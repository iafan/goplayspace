package navigator

import "github.com/gopherjs/gopherjs/js"

// Platform is a wrapper for `navigator.platform`
func Platform() string {
	return js.Global.Get("navigator").Get("platform").String()
}

// UserAgent is a wrapper for `navigator.userAgent`
func UserAgent() string {
	return js.Global.Get("navigator").Get("userAgent").String()
}
