package util

import "time"

// Schedule implements a deferred function call
// (which will typically be exectuted after DOM rendering is complete)
func Schedule(f func()) {
	time.AfterFunc(0, f)
}
