package editor

import "github.com/iafan/goplayspace/client/component/editor/undo"

func (ed *Editor) getStateAsUndoEntry() *undo.Entry {
	return &undo.Entry{
		Text:     ed.ta.GetValue(),
		SelStart: ed.ta.GetSelectionStart(),
		SelEnd:   ed.ta.GetSelectionEnd(),
	}
}

func (ed *Editor) saveState() {
	if ed.getTextarea() == nil {
		return
	}

	text := ed.ta.GetValue()
	if text == "" || text == ed.InitialValue {
		return
	}

	ss := ed.ta.GetSelectionStart()
	se := ed.ta.GetSelectionEnd()

	if state := ed.UndoStack.CurrentState(); state == nil ||
		state.Text != text || state.SelStart != ss || state.SelEnd != se {
		ed.UndoStack.Push(ed.getStateAsUndoEntry())
	}
}

// Undo does one undo step
func (ed *Editor) Undo() {
	if ed.getTextarea() == nil {
		return
	}

	ed.saveState()

	if !ed.UndoStack.CanUndo() {
		return
	}

	entry := ed.UndoStack.Undo()
	if entry != nil {
		ed.ta.SetValue(entry.Text)
		ed.ta.SetSelectionStart(entry.SelStart)
		ed.ta.SetSelectionEnd(entry.SelEnd)
	} else {
		ed.ta.SetValue("")
	}
	ed.onChange(nil)
}

// Redo does one redo step
func (ed *Editor) Redo() {
	if ed.getTextarea() == nil {
		return
	}

	ed.saveState()

	if !ed.UndoStack.CanRedo() {
		return
	}

	entry := ed.UndoStack.Redo()
	ed.ta.SetValue(entry.Text)
	ed.ta.SetSelectionStart(entry.SelStart)
	ed.ta.SetSelectionEnd(entry.SelEnd)
	ed.onChange(nil)
}
