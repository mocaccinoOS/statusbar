package blocks

import (

	//	"github.com/0xAX/notificator"

	"github.com/getlantern/systray"
)

const (
	TextKey string = "Text"
)

type Text struct {
	Text    string
	SubText string
	URL     string
	SubMenu *systray.MenuItem
}

func (c *Text) Menu(n Notifier, r Renderer, sm SessionManager) {
	t := systray.AddMenuItem(c.Text, c.SubText)
	t.Disable()
}

func (c *Text) Close()     {}
func (c *Text) ID() string { return TextKey }

func (c *Text) Render(Notifier) string {
	return ""
}
