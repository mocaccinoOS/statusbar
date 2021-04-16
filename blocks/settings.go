package blocks

import (
	"log"
	"net/http"

	"github.com/rakyll/statik/fs"

	_ "github.com/MocaccinoOS/statusbar/statik"
)

const (
	SettingsKey string = "Settings"
)

type Settings struct{}

func (c *Settings) Menu(n Notifier, r Renderer) {
	uo := &ChromeEmbeddedOpener{URL: "http://127.0.0.1:8080", Text: "Settings"}
	uo.Menu(n, r)
	go func() {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		// Serve the contents over HTTP.
		http.Handle("/", http.FileServer(statikFS))
		http.ListenAndServe("127.0.0.1:8080", nil)
	}()
}

func (c *Settings) Close()     {}
func (c *Settings) ID() string { return SettingsKey }

func (c *Settings) Render(Notifier) string {
	return ""
}
