package util

import (
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// Schedule implements a deferred function call
// (which will typically be exectuted after DOM rendering is complete)
func Schedule(f func()) {
	time.AfterFunc(0, f)
}

// IsSafari returns true if User-Agent is a Safari browser
func IsSafari() bool {
	ua := js.Global.Get("navigator").Get("userAgent").String()
	return strings.Contains(ua, "Safari") && !strings.Contains(ua, "Chrome")
}
