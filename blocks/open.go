package blocks

import (

	//	"github.com/0xAX/notificator"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

const (
	OpenKey string = "Open"
)

type Open struct {
	Text    string
	SubText string
	URL     string
}

func (c *Open) Menu(n Notifier, r Renderer) {
	url := systray.AddMenuItem(c.Text, c.SubText)
	go func() {
		for range url.ClickedCh {
			open.Run(c.URL)
		}
	}()
}

func (c *Open) Close()     {}
func (c *Open) ID() string { return OpenKey }

func (c *Open) Render(Notifier) string {
	return ""
}
