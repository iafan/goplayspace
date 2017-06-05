package history

import "github.com/gopherjs/gopherjs/js"

func ReplaceState(url string) {
	js.Global.Get("history").Call("replaceState", "", "", url)
}
