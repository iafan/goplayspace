package help

import (
	"regexp"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

var pkgFuncSplitterR = regexp.MustCompile(`^(\w[\w\d_]*)\.(\w[\w\d_]*)$`)

// this variable holds iframe URL across re-renders, since Browser object
// is recreated every time; this is needed to avoid resetting previous help topic
var cachedURL string

// Browser renders iframe with appropriate topic
type Browser struct {
	vecty.Core

	Topic   string            `vecty:"prop"`
	Imports map[string]string `vecty:"prop"`
}

func (h *Browser) getURL() string {
	// determine help URL based on selected text

	if h.Imports == nil {
		panic("*HelpBrowser.Imports is nil")
	}

	godocRoot := "https://godoc.org/"

	// test for `package.ident`
	if matches := pkgFuncSplitterR.FindAllStringSubmatch(h.Topic, 1); matches != nil {
		pkg := matches[0][1]
		ident := matches[0][2]
		if h.Imports[pkg] != "" {
			pkg = h.Imports[pkg]
		}
		return godocRoot + pkg + "/#" + ident

	} else if h.Imports[h.Topic] != "" { // test for `package` name (both short and fully-qualified)
		return godocRoot + h.Imports[h.Topic] + "/#pkg-overview"

	} else if hash := langspec[h.Topic]; hash != "" { // test for keywords, built-in functions and special strings
		return langspecRoot + hash
	}
	return ""
}

func (h *Browser) getCachedURL() string {
	s := h.getURL()
	if s == "" {
		return cachedURL
	}
	cachedURL = s
	return s
}

// Render implements the vecty.Component interface.
func (h *Browser) Render() vecty.ComponentOrHTML {
	return elem.InlineFrame(
		vecty.Markup(
			vecty.Class("help-browser"),
			vecty.Property("src", h.getCachedURL()),
		),
	)
}
