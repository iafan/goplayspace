About Go Play Space
===================

Go Play Space is an experimental alternative [Go Playground](https://play.golang.org)
frontend that is built in Go itself (using [GopherJS](https://github.com/gopherjs/gopherjs)),
a Go&rarr;JavaScript transpiler, and [Vecty](https://github.com/gopherjs/vecty),
a React-like frontend library for GopherJS).

Try it: [https://goplay.space](https://goplay.space/)
=====

![demo](https://cloud.githubusercontent.com/assets/1728158/26770686/b3f0a4d0-496d-11e7-8be2-9ab88e856b8c.gif)

Go Play Space supports the [Turtle graphics mode](https://goplay.space/#draw) to help visualize algorithms and make learning experience more fun.

![draw mode demo](https://user-images.githubusercontent.com/1728158/29108287-12dfd26a-7c93-11e7-966d-713356bce3d3.gif)

Features
--------

1. Syntax highlighting, auto-closing braces and quotes, proper undo/redo, auto indentation
2. Smart help lookup: double-click on e.g. <code>package</code> keyword
   or <code>Println</code> function name in source code, and you will see
   the relevant help topic.
3. Live syntax error checking
4. Error line highlighting (both for syntax errors and for errors
   returned from the compiler)
5. Ability to highlight lines and blocks of code (like on Github, but better!) —
   just click on the line numbers. Use <kbd>Shift</kbd> and <kbd>Ctrl</kbd>
   to modify the selection
6. Keyboard shortcuts (see button captions)
7. Support for several UI themes
8. Support for [Fira Code](https://github.com/tonsky/FiraCode) font
   (either the one installed in your system or a webfont)
9. `go imports` is always run before running your code, so you don't usually
   have to worry about imports at all

Code execution is proxied to the official Go Playground, so your programs will work the same.
Shared snippets are also stored on golang.org servers.

Running Locally
----------------

Download the package:

```sh
$ go get -u github.com/iafan/goplayspace/...
```

Compile both client-side code and server binary:
```sh
$ cd $GOPATH/src/github.com/iafan/goplayspace/bin
$ ./build-client && ./build-server
```

Run the server:

```sh
$ ./goplayspace
```

Then open http://localhost:8080/ in your browser.

Troubleshooting
---------------

If you have trouble compiling the client, please make sure you have the latest version of GopherJS installed by running `go get -u github.com/gopherjs/gopherjs` (see #6)

Feedback
--------

Feel free to provide your feedback, suggestions or bug reports here in the <a href="https://github.com/iafan/goplayspace/issues">bug tracker</a>, or message [@afan](https://gophers.slack.com/messages/@afan/) in the [Gophers Slack channel](https://gophersinvite.herokuapp.com/).

Credits
-------

Gopher vector logo by [Takuya Ueda](https://twitter.com/tenntenn),
licensed under the Creative Commons 3.0 Attributions license and based
on original artwork by [Renee French](http://reneefrench.blogspot.com/).
See https://github.com/golang-samples/gopher-vector

Go proverbs: [Rob Pike](https://twitter.com/rob_pike)
