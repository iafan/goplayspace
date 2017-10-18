package editor

import (
	"github.com/gopherjs/gopherjs/js"
)

// Shadow contains the logic behind the shadow syntax highlighter
// exposed on the application page under '.shadow' class
type Shadow struct {
	*js.Object
}

// GetHeight gets the height of the shadow div in pixels
func (s *Shadow) GetHeight() int {
	return s.Get("offsetHeight").Int()
}

// SetValue sets the inner HTML (highlighted code)
func (s *Shadow) SetValue(html string) {
	s.Set("innerHTML", html)
}
