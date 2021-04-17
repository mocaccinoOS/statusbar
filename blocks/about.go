package blocks

import (

	//	"github.com/0xAX/notificator"

	"fmt"

	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
)

const (
	AboutKey  string = "About"
	AboutText string = `
MocaccinoOS statusbar version %s

MocaccinoOS statusbar Copyright (C) 2021 Ettore Di Giacinto
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it
under certain conditions.

Source code available at: https://github.com/mocaccinoOS/statusbar
`
)

type About struct{ Version string }

func (c *About) Menu(n Notifier, r Renderer) {
	about := systray.AddMenuItem("About", "")
	go func() {
		for range about.ClickedCh {
			dialog.Message(fmt.Sprintf(AboutText, c.Version)).Title("About").Info()
		}
	}()
}

func (c *About) Close()     {}
func (c *About) ID() string { return AboutKey }

func (c *About) Render(Notifier) string {
	return ""
}
