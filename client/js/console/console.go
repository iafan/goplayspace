package console

import "github.com/gopherjs/gopherjs/js"

func Log(params ...interface{}) {
	js.Global.Get("console").Call("log", params...)
}

func Time(label string) {
	js.Global.Get("console").Call("time", label)
}

func TimeEnd(label string) {
	js.Global.Get("console").Call("timeEnd", label)
}
