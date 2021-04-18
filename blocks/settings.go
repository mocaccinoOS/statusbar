package blocks

import (
	_ "github.com/MocaccinoOS/statusbar/statik"
)

const (
	SettingsKey string = "Settings"
)

type Settings struct{}

func (c *Settings) Menu(n Notifier, r Renderer) {
	uo := &ChromeEmbeddedOpener{URL: "http://127.0.0.1:9910/settings", Text: "Settings"}
	uo.Menu(n, r)
}

func (c *Settings) Close()     {}
func (c *Settings) ID() string { return SettingsKey }

func (c *Settings) Render(Notifier) string {
	return ""
}
