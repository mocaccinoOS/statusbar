package blocks

import (

	//	"github.com/0xAX/notificator"

	"github.com/getlantern/systray"
)

const (
	SeparatorKey string = "Separator"
)

type Separator struct{}

func (c *Separator) Menu(n Notifier, r Renderer, sm SessionManager) {

	systray.AddSeparator()

}

func (c *Separator) Close()     {}
func (c *Separator) ID() string { return SeparatorKey }

func (c *Separator) Render(Notifier) string {
	return ""
}
