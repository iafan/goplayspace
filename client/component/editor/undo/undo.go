package undo

// Entry holds information about single undo entry
type Entry struct {
	Text     string
	SelStart int
	SelEnd   int
}

// Stack holds undo and redo entries
type Stack struct {
	entries []*Entry
	pos     int
	end     int
}

// Push pushes new state to the stack
func (s *Stack) Push(entry *Entry) {
	if s.pos == len(s.entries)-1 {
		// shift stack items to the left
		for i := 1; i < len(s.entries); i++ {
			s.entries[i-1] = s.entries[i]
		}
	} else {
		s.pos = s.pos + 1
	}
	s.entries[s.pos] = entry
	s.end = s.pos // clear redo stack
}

// CanUndo returns true if the undo list is not empty
func (s *Stack) CanUndo() bool {
	return s.pos > 0
}

// CanRedo returns true if the redo list is not empty
func (s *Stack) CanRedo() bool {
	return s.pos < s.end
}

// CurrentState returns current state
// (entry current position is pointing to)
func (s *Stack) CurrentState() *Entry {
	if s.pos == -1 {
		return nil
	}
	return s.entries[s.pos]
}

// Undo does one undo step: moves current position one step back
// and returns the top entry on the undo stack
// (or nil if there are no more entries)
func (s *Stack) Undo() *Entry {
	if !s.CanUndo() {
		return nil
	}
	s.pos = s.pos - 1
	return s.CurrentState()
}

// Redo does one redo step: moves current position one step forward
// from redo to undo stack and returns the top entry on the undo stack
// (or nil if there are no more entries)
func (s *Stack) Redo() *Entry {
	if !s.CanRedo() {
		return nil
	}
	s.pos = s.pos + 1
	return s.CurrentState()
}

// NewStack initializes and returns a Stack instance
func NewStack(maxSize uint) *Stack {
	if maxSize < 1 {
		panic("undo.NewStack: maxSize should be a positive number")
	}
	return &Stack{
		entries: make([]*Entry, maxSize+1), // +1 because the current state takes one spot
		pos:     -1,
		end:     -1,
	}
}
