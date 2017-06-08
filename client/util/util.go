package util

import (
	"strings"
	"time"

	"github.com/iafan/goplayspace/client/js/navigator"
)

// Schedule implements a deferred function call
// (which will typically be exectuted after DOM rendering is complete)
func Schedule(f func()) {
	time.AfterFunc(0, f)
}

// IsSafari returns true if User-Agent is a Safari browser
func IsSafari() bool {
	ua := navigator.UserAgent()
	return strings.Contains(ua, "Safari") && !strings.Contains(ua, "Chrome")
}

// IsMacOS returns true under Mac OS
func IsMacOS() bool {
	return strings.HasPrefix(navigator.Platform(), "Mac")
}

// IsIOS returns true under iOS
func IsIOS() bool {
	p := navigator.Platform()
	return strings.HasPrefix(p, "iPhone") || strings.HasPrefix(p, "iPad") || strings.HasPrefix(p, "iPod")
}
