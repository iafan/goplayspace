package main

import (
	"github.com/gopherjs/vecty"
	"github.com/iafan/goplayspace/client/component/app"
	"github.com/iafan/goplayspace/client/js/localstorage"
)

func main() {
	vecty.SetTitle(app.PageTitle)

	a := &app.Application{
		Theme:            localstorage.Get("theme", "light"),
		TabWidth:         localstorage.GetInt("tab-width", 4),
		HighlightingMode: localstorage.GetBool("highlighting", true),
	}

	vecty.RenderBody(a)
	a.WaitForPageLoaded()
}
