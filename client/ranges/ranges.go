package ranges

import (
	"strconv"
	"strings"

	"github.com/iafan/goplayspace/client/js/console"
)

// LineRange represents the range of selected lines
type LineRange struct {
	Begin int
	End   int
}

func (r *LineRange) String() string {
	if r.Begin == r.End {
		return strconv.Itoa(r.Begin)
	}
	return strconv.Itoa(r.Begin) + "-" + strconv.Itoa(r.End)
}

// Range maintains the list of ranges
type Range struct {
	Sel      []*LineRange
	selPoint int
}

// String generates the string representation of the Range object
func (r *Range) String() string {
	if r == nil || len(r.Sel) == 0 {
		return ""
	}

	a := make([]string, len(r.Sel))
	for i := 0; i < len(r.Sel); i++ {
		a[i] = r.Sel[i].String()
	}
	return strings.Join(a, ",")
}

// HasSelection returns true if there's some selection
func (r *Range) HasSelection() bool {
	return len(r.Sel) > 0
}

// ClearSelection clears selection
func (r *Range) ClearSelection() {
	if r.HasSelection() {
		r.Sel = nil
		r.selPoint = 0
	}
}

// AddRange adds a selection range
// potentially consolidating other ranges
func (r *Range) AddRange(begin, end int) {
	r.selPoint = end

	out := make([]*LineRange, 0)

	l := len(r.Sel)

	// Special case: if the previous list of ranges is empty,
	// return the range with the new selection
	if l == 0 {
		r.Sel = append(out, &LineRange{begin, end})
		return
	}

	// Special case: if the new range goes after the last selection range
	// and doesn't touch it, just append the new range to the slice
	if begin > r.Sel[l-1].End+1 {
		r.Sel = append(r.Sel, &LineRange{begin, end})
		return
	}

	for i := 0; i < l; i++ {
		b := r.Sel[i].Begin
		e := r.Sel[i].End
		//console("[:1]", "i:", i, "b:", b, "e:", e, "begin:", begin, "end:", end)

		// If the new range goes before the current range and doesn't touch it,
		// append it to selections, then append the tail and return
		if end < b-1 {
			out = append(out, &LineRange{begin, end})
			r.Sel = append(out, r.Sel[i:]...)
			return
		}

		// If the new range goes after the current range and doesn't touch it,
		// append current selection to the output but don't modify the range
		if begin > e+1 {
			out = append(out, r.Sel[i])
			continue
		}

		// The new range overlaps or touches the existing one,
		// adjust the new range to consume the existing one

		if b < begin {
			begin = b
		}
		if e > end {
			end = e
		}

		// The range

		// if it's the last range, add it to the output
		if i == l-1 {
			out = append(out, &LineRange{begin, end})
		}
	}

	r.Sel = out
}

// SetRange resets the selection with one given range
func (r *Range) SetRange(begin, end int) {
	r.Sel = nil
	r.AddRange(begin, end)
}

// IsLineSelected returns true if the line number is in any of the
// selection ranges
func (r *Range) IsLineSelected(n int) bool {
	for i := 0; i < len(r.Sel); i++ {
		if n >= r.Sel[i].Begin && n <= r.Sel[i].End {
			return true
		}
	}
	return false
}

// IsOnlyLineSelected returns true if the line number is the only
// line selected
func (r *Range) IsOnlyLineSelected(n int) bool {
	return len(r.Sel) == 1 && r.Sel[0].Begin == n && r.Sel[0].End == n
}

// ToggleLine either adds or removes the line from selection
func (r *Range) ToggleLine(n int) {
	if r.IsLineSelected(n) {
		r.RemoveRange(n, n)
	} else {
		r.AddRange(n, n)
	}
}

// AddSelPoint adds selection between previous
// selection point and the one provided
// (to be used with Shift+clicking to select ranges with the mouse)
func (r *Range) AddSelPoint(end int) {
	if r.selPoint == 0 {
		r.AddRange(end, end)
		return
	}

	begin := r.selPoint
	if begin > end {
		begin, end = end, begin
	}
	r.AddRange(begin, end)
}

// RemoveRange removes selection range
// potentially breaking apart or adjusting other ranges
func (r *Range) RemoveRange(begin, end int) {
	r.selPoint = 0
	out := make([]*LineRange, 0)
	for i := 0; i < len(r.Sel); i++ {
		b := r.Sel[i].Begin
		e := r.Sel[i].End

		// if exclusion region completely covers
		// the selection, skip the selection
		if begin <= b && end >= e {
			continue
		}

		// if exclusion region completely surrounded
		// by the selection, break selection into two
		if begin > b && end < e {
			out = append(out, &LineRange{b, begin - 1})
			out = append(out, &LineRange{end + 1, e})
			continue
		}

		// if exclusion region completely falls outside,
		// output the existing one
		if begin > e || end < b {
			out = append(out, &LineRange{b, e})
			continue
		}

		// otherwise, adjust the region from either side
		if begin > b {
			e = begin - 1
		}

		if end < e {
			b = end + 1
		}
		out = append(out, &LineRange{b, e})
	}

	r.Sel = out
}

// Parse parses string representation of ranges
func (r *Range) Parse(s string) error {
	if s == "" {
		return nil
	}

	if tokens := strings.Split(s, ","); len(tokens) > 0 {
		r.Sel = nil

		for i := 0; i < len(tokens); i++ {
			a := strings.SplitN(tokens[i], "-", 2)
			start, err := strconv.Atoi(a[0])
			if err != nil {
				return err
			}
			end := start
			if len(a) > 1 {
				end, err = strconv.Atoi(a[1])
				if err != nil {
					return err
				}
			}
			r.AddRange(start, end)
		}
	}
	return nil
}

// New returns a new range object from string representation of ranges
func New(s string) *Range {
	r := &Range{}
	if err := r.Parse(s); err != nil {
		console.Log("ranges.New(): parse error:", err)
	}
	return r
}
