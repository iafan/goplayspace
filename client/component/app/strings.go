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

	<h2>Features</h2>

	<ol>
		<li><strong>New:</strong> an embedded <a href="#draw">Turtle graphics mode</a>. Introduce Go to your kids, or have fun yourself!</li>
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

	<h2>Feedback</h2>

	<p>Feel free to provide your feedback, suggestions or bug reports in the <a href="https://github.com/iafan/goplayspace/issues">official bug tracker</a>, or message <a href="https://gophers.slack.com/messages/@afan/">@afan</a> in the <a href="https://gophersinvite.herokuapp.com/">Gophers Slack channel</a>.</p>
`

var drawHelpHTML = `
	<h1>Turtle Graphics Mode</h1>

	<p>Go Play Space supports the <a href="https://en.wikipedia.org/wiki/Turtle_graphics" target="_top">Turtle graphics</a> mode to help visualize algorithms and make learning experience more fun.</p>

	<p>Whenever Go code is executed and produces some console output, this output is parsed, and found control commands are interpreted on a drawing board. You can start with something basic like one <code>fmt.println</code> statement that prints all the commands sequentially,
	and then build your own functional API and algorithms — in Go — that would produce the desired sequence of control commands.</p>

	<h2>Try These Examples</h2>

	<ol>
		<li><a href="#wT_eZWJT69">Star</a> — an example on how to draw a star</li>
		<li><a href="#4GFA2un9jL">House</a> — an example on drawing a house with multiple colors</li>
		<li><a href="#61SJKVrWwj">Tree</a> — an example on using recursion to draw a tree</li>
		<li><a href="#S6FsspIE6d">Colored squares</a> — an example of using a simple “API” wrapper functions and <em>for</em> loops</li>
		<li><a href="#lAca11gTvc">Colored squares (Russian)</a> — the same “Colored squares” example above, but with function/variable names in Russian. This example demonstrates how you can offer your kids a drawing API in their spoken language of choice.</li>
	</ol>

	<h2>Control Commands Reference</h2>

	<ol>
		<li><code>draw mode</code> — triggers the draw mode; any commands before this line are ignored; gopher starts at the center of the board</li>
		<li><code>left</code> — equivalent to <code>left 90</code>: turn 90 degrees counter-clockwise</li>
		<li><code>left <em>N</em></code> — N degrees counter-clockwise [0..360]; fractional values allowed</li>
		<li><code>right</code> — equivalent to <code>right 90</code>: turn 90 degrees clockwise</li>
		<li><code>right <em>N</em></code> — N degrees clockwise [0..360]; fractional values allowed</li>
		<li><code>forward</code> — make one grid step</li>
		<li><code>forward <em>N</em></code> — make <N> grid steps; fractional values allowed</li>
		<li><code>color off</code> — turn color off (leave no trace; this is the default state)</li>
		<li><code>color <em>Color</em></code> — set stroke color; any web color would work, e.g. “red”, “black”, “#f5c0e2”, “rgba(0,0,0,0.3)”</li>
		<li><code>width <em>N</em></code> — set brush width to <N> pixels; fractional values allowed; default is 2</li>
		<li><code>say <em>Message</em></code> — shows a temporary ’balloon’ message</li>
	</ol>

	<h2>How Commands Are Interpreted</h2>

	<p>Once your Go code produces the full console output, this output is split into lines, each line is trimmed and then checked against the available commands. Any output before <code>draw mode</code> line is ignored. There must be only one command per line. Lines that don’t match any of the commands above are ignored, so you can safely print debug message lines along with the actual commands.</p>

	<h2>Feedback</h2>

	<p>Feel free to provide your feedback, suggestions or bug reports in the <a href="https://github.com/iafan/goplayspace/issues">official bug tracker</a>, or message <a href="https://gophers.slack.com/messages/@afan/">@afan</a> in the <a href="https://gophersinvite.herokuapp.com/">Gophers Slack channel</a>.</p>

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

var initialDrawCode = `package main

import (
	"fmt"
)

func main() {
	fmt.Println(` + "`" + `
		draw mode
		color black
		forward 2
		right
		forward 2
		right
		forward 2
		right
		forward 2
	` + "`" + `)
}
`
var initialDrawCaretPos = 64
