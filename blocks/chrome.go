package blocks

import (
	"log"
	//	"github.com/0xAX/notificator"

	"github.com/getlantern/systray"
	"github.com/zserge/lorca"
)

const (
	ChromeEmbeddedOpenerKey string = "ChromeEmbeddedOpener"
)

type ChromeEmbeddedOpener struct {
	Text    string
	SubText string
	URL     string
}

func (c *ChromeEmbeddedOpener) Menu(n Notifier, r Renderer) {

	url := systray.AddMenuItem(c.Text, c.SubText)

	//n.Push("test", "test", "", notificator.UR_CRITICAL)
	go func() {

		for range url.ClickedCh {
			go func() {
				ui, err := lorca.New(c.URL, "", 480, 320)
				// ui, err := lorca.New("data:text/html,"+url.PathEscape(`
				// <html>
				// 	<head><title>Hello</title></head>
				// 	<body><h1>Hello, world!</h1></body>
				// </html>
				// `), "", 480, 320)

				ui.SetBounds(lorca.Bounds{WindowState: lorca.WindowStateMaximized})
				if err != nil {
					log.Fatal(err)
				}
				defer ui.Close()
				// Wait until UI window is closed
				<-ui.Done()
			}()
		}

	}()
}

func (c *ChromeEmbeddedOpener) Close()     {}
func (c *ChromeEmbeddedOpener) ID() string { return ChromeEmbeddedOpenerKey }

func (c *ChromeEmbeddedOpener) Render(Notifier) string {
	return ""
}
