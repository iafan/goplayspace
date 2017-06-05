package hash

import (
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/iafan/goplayspace/client/js/history"
)

// Hash contains the state parsed from URL hash
type Hash struct {
	ID     string
	Ranges string

	isUpdating bool

	OnChange func(h *Hash)
}

func (h *Hash) updateAddressBar() {
	h.isUpdating = true
	history.ReplaceState(h.url())
	h.isUpdating = false
}

func (h *Hash) url() string {
	if h.ID == "" && h.Ranges == "" {
		return "/"
	}

	if h.ID != "" && h.Ranges == "" {
		return "/#" + h.ID
	}

	return "/#" + h.ID + "," + h.Ranges
}

// Reset resets the hash properties
func (h *Hash) Reset() {
	h.ID = ""
	h.Ranges = ""
	h.updateAddressBar()
}

// SetID sets ID part and updates state (URL in the address bar)
func (h *Hash) SetID(id string) {
	h.ID = id
	h.updateAddressBar()
}

// SetRanges sets Ranges part and updates state (URL in the address bar)
func (h *Hash) SetRanges(ranges string) {
	h.Ranges = ranges
	h.updateAddressBar()
}

func (h *Hash) onHashChange() {
	if h.isUpdating {
		return
	}
	h.parse()
	if h.OnChange != nil {
		h.OnChange(h)
	}
}

func (h *Hash) parse() {
	s := js.Global.Get("window").Get("location").Get("hash").String()[1:]
	h.ID = s
	if tokens := strings.SplitN(s, ",", 2); len(tokens) > 1 {
		h.ID = tokens[0]
		h.Ranges = tokens[1]
	}
}

// New returns a new Hash instance filled with values
// from window.localtion.hash
func New(onChange func(h *Hash)) *Hash {
	h := &Hash{
		OnChange: onChange,
	}
	h.parse()
	js.Global.Get("window").Set("onhashchange", h.onHashChange)
	return h
}

// FIXME: we need to somehow communicate changes from editor's Ranges object down to Hash.Ranges
