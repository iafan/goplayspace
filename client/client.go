package main

import (
	"github.com/gopherjs/vecty"
	"github.com/iafan/goplayspace/client/component/app"
	"github.com/iafan/goplayspace/client/js/localstorage"
)

func main() {
	vecty.SetTitle(app.PageTitle)

	theme := localstorage.Get("theme")
	if theme == "" {
		theme = "light"
	}

	a := &app.Application{
		Theme: theme,
	}

	vecty.RenderBody(a)
	a.WaitForPageLoaded()
}
