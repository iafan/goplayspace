About Go Play Space
===================

Go Play Space is an experimental alternative [Go Playground](https://play.golang.org)
frontend that is built in Go itself (using [GopherJS](https://github.com/gopherjs/gopherjs)),
a Go&rarr;JavaScript transpiler, and [Vecty](https://github.com/gopherjs/vecty),
a React-like frontend library for GopherJS).

![demo](https://cloud.githubusercontent.com/assets/1728158/26770686/b3f0a4d0-496d-11e7-8be2-9ab88e856b8c.gif)

### Try it yourself: [https://goplay.space &rarr;](https://goplay.space/)

Main differences from the official Go Playground:

1. Syntax highlighting
2. Golang help lookup: double-click on e.g. <code>package</code> keyword
   or <code>Println</code> function name in source code, and you will see
   the relevant help topic.
3. Live syntax error checking
4. Error line highlighting (both for syntax errors and for errors
   returned from the compiler)
5. Ability to highlight lines and blocks of code (like on Github, but better!) —
   just click on the line numbers. Use <kbd>Shift</kbd> and <kbd>Ctrl</kbd>
   to modify the selection
6. Keyboard shortcuts (see button captions). Note that both
   <kbd>Cmd</kbd> and <kbd>Ctrl</kbd> keys work interchangeably
7. Support for several UI themes
8. Support for [Fira Code](https://github.com/tonsky/FiraCode) font
   (if it is installed in your system)
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

Credits
-------

Gopher vector logo by [Takuya Ueda](https://twitter.com/tenntenn),
licensed under the Creative Commons 3.0 Attributions license and based
on original artwork by [Renee French](http://reneefrench.blogspot.com/).
See https://github.com/golang-samples/gopher-vector

Go proverbs: [Rob Pike](https://twitter.com/rob_pike)
