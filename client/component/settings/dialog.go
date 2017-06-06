package settings

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// Dialog contains the logic behind the settings dialog
// exposed on the applicaiton page under '.settings-dialog' class
type Dialog struct {
	vecty.Core
	//node *js.Object

	Theme            string
	TabWidth         int
	FontWeight       string
	UseWebfont       bool
	HighlightingMode bool

	OnChange func(d *Dialog)
}

/*
func (d *Dialog) getDOMNode() *js.Object {
	if d.node == nil {
		d.node = js.Global.Get("document").Call("querySelector", ".settings-dialog")
	}
	return d.node
}
*/

func (d *Dialog) updateTheme(e *vecty.Event) {
	d.Theme = e.Target.Get("value").String()
	d.fireOnChangeEvent()
}

func (d *Dialog) updateTabWidth(e *vecty.Event) {
	d.TabWidth = e.Target.Get("value").Int()
	d.fireOnChangeEvent()
}

func (d *Dialog) updateFontWeight(e *vecty.Event) {
	d.FontWeight = e.Target.Get("value").String()
	d.fireOnChangeEvent()
}

func (d *Dialog) updateUseWebfont(e *vecty.Event) {
	d.UseWebfont = e.Target.Get("value").Bool()
	d.fireOnChangeEvent()
}

func (d *Dialog) updateHighlighting(e *vecty.Event) {
	d.HighlightingMode = e.Target.Get("checked").Bool()
	d.fireOnChangeEvent()
}

func (d *Dialog) fireOnChangeEvent() {
	if d.OnChange != nil {
		d.OnChange(d)
	}
}

// Render implements the vecty.Component interface.
func (d *Dialog) Render() *vecty.HTML {
	return elem.Div(
		vecty.ClassMap{"settings-dialog": true},
		elem.Paragraph(
			elem.Div(
				vecty.Text("Theme:"),
			),
			elem.Select(
				elem.Option(
					vecty.Property("value", "space"),
					vecty.Property("selected", d.Theme == "space"),
					vecty.Text("Space"),
				),
				elem.Option(
					vecty.Property("value", "classic"),
					vecty.Property("selected", d.Theme == "classic"),
					vecty.Text("Classic"),
				),
				elem.Option(
					vecty.Property("value", "light"),
					vecty.Property("selected", d.Theme == "light"),
					vecty.Text("Light"),
				),
				elem.Option(
					vecty.Property("value", "dark"),
					vecty.Property("selected", d.Theme == "dark"),
					vecty.Text("Dark"),
				),
				event.Change(d.updateTheme),
			),
		),
		elem.Paragraph(
			elem.Div(
				vecty.Text("Tab width:"),
			),
			elem.Select(
				elem.Option(
					vecty.Property("value", "2"),
					vecty.Property("selected", d.TabWidth == 2),
					vecty.Text("2"),
				),
				elem.Option(
					vecty.Property("value", "4"),
					vecty.Property("selected", d.TabWidth == 4),
					vecty.Text("4"),
				),
				elem.Option(
					vecty.Property("value", "6"),
					vecty.Property("selected", d.TabWidth == 6),
					vecty.Text("6"),
				),
				elem.Option(
					vecty.Property("value", "8"),
					vecty.Property("selected", d.TabWidth == 8),
					vecty.Text("8"),
				),
				event.Change(d.updateTabWidth),
			),
		),
		elem.Paragraph(
			elem.Div(
				vecty.Text("Font weight:"),
			),
			elem.Select(
				elem.Option(
					vecty.Property("value", "lighter"),
					vecty.Property("selected", d.FontWeight == "lighter"),
					vecty.Text("Lighter"),
				),
				elem.Option(
					vecty.Property("value", "normal"),
					vecty.Property("selected", d.FontWeight == "normal"),
					vecty.Text("Normal"),
				),
				event.Change(d.updateFontWeight),
			),
		),
		elem.Paragraph(
			elem.Div(
				vecty.Text("‘Fira Code’ font source:"),
			),
			elem.Select(
				elem.Option(
					vecty.Property("value", ""),
					vecty.Property("selected", !d.UseWebfont),
					vecty.Text("System"),
				),
				elem.Option(
					vecty.Property("value", 1),
					vecty.Property("selected", d.UseWebfont),
					vecty.Text("Webfont"),
				),
				event.Change(d.updateUseWebfont),
			),
		),
		elem.Paragraph(
			elem.Input(
				vecty.Property("id", "highlighting"),
				vecty.Property("type", "checkbox"),
				vecty.Property("checked", d.HighlightingMode),
				event.Change(d.updateHighlighting),
			),
			elem.Label(
				vecty.Attribute("for", "highlighting"),
				vecty.Text("Syntax highlighting"),
			),
		),
	)
}
