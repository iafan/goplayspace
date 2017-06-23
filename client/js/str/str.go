package str

import "unicode/utf16"

// UTF16ToUTF8Pos converts UTF-16 position (as reported by JavaScript
// calls, e.g. selectionStart / selectionEnd) to a UTF-8 position
// in the string
func UTF16ToUTF8Pos(s string, i int) int {
	return len(string(utf16.Decode(utf16.Encode([]rune(s))[:i])))
}

// UTF16ToUTF8Pos converts UTF-8 position in the string
// to a UTF16 position suitable for e.g. caret positioning
func UTF8ToUTF16Pos(s string, i int) int {
	return len(utf16.Encode([]rune(s[:i])))
}
