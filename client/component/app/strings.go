package app

import "time"

// PageTitle is set as a <title> once the app is loaded
var PageTitle = "The Go Play Space"

var helpHTML = `
	<h1>About Go Play Space</h1>
	<p>This is an experimental alternative <a href="https://play.golang.org">Go Playground</a> frontend that is built in Go itself
	(using <a href="https://github.com/gopherjs/gopherjs">GopherJS</a>, a Go&rarr;JavaScript transpiler, and
	<a href="https://github.com/gopherjs/vecty">Vecty</a>, a React-like frontend library for GopherJS).</p>

	<p>
		<a href="https://github.com/iafan/goplayspace">View source code on GitHub</a>
		<iframe class="github-button" src="https://ghbtns.com/github-btn.html?user=iafan&amp;repo=goplayspace&amp;type=watch&amp;count=true" allowtransparency="true" frameborder="0" scrolling="0" width="110" height="20"></iframe>
	</p>

	<p>Main differences from the official Go Playground:</p>

	<ol>
		<li>Syntax highlighting, auto-closing braces and quotes, proper undo/redo, auto indentation</li>
		<li>Smart help lookup: double-click on e.g. <code>package</code> keyword or <code>Println</code> function name in source code,
		and you will see the relevant help topic. Try it!</li>
		<li>Live syntax error checking</li>
		<li>Error line highlighting (both for syntax errors and for errors returned from the compiler)</li>
		<li>Ability to highlight lines and blocks of code (like on Github, but better!) — just click on the line numbers. Use <kbd>Shift</kbd> and <kbd>Ctrl</kbd> to modify the selection</li>
		<li>Keyboard shortcuts (see button captions)</li>
		<li>Support for several UI themes and UI tweaks (see the Settings button)</li>
		<li>Support for <a href="https://github.com/tonsky/FiraCode">Fira Code</a> font (either the one installed in your system or a webfont)</li>
		<li><code>go imports</code> is always run before running your code, so you don't usually have to worry
		about imports at all</li>
	</ol>

	<p>
		Code execution is proxied to the official Go Playground, so your programs will work the same.
		Shared snippets are also stored on golang.org servers.
		Any requests for content removal should be directed to <a href="mailto:security@golang.org">security@golang.org</a>.
		Please include the URL and the reason for the request.
	</p>
`

var blankTemplatePos = 29

// template has the tabulation symbol on purpose!
var blankTemplate = `package main

func main() {
	
}
`

var proverbs = []string{
	"Don't communicate by sharing memory, share memory by communicating.",
	"Concurrency is not parallelism.",
	"Channels orchestrate; mutexes serialize.",
	"The bigger the interface, the weaker the abstraction.",
	"Make the zero value useful.",
	"interface{} says nothing.",
	"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
	"A little copying is better than a little dependency.",
	"Syscall must always be guarded with build tags.",
	"Cgo must always be guarded with build tags.",
	"Cgo is not Go.",
	"With the unsafe package there are no guarantees.",
	"Clear is better than clever.",
	"Reflection is never clear.",
	"Errors are values.",
	"Don't just check errors, handle them gracefully.",
	"Design the architecture, name the components, document the details.",
	"Documentation is for users.",
	"Don't panic.",
}

var greeting = proverbs[time.Now().Unix()%int64(len(proverbs))]
var initialCaretPos = 61 + len(greeting)

var initialCode = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("` + greeting + `")
}
`
