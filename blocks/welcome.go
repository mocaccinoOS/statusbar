package blocks

import (
	"log"
	//	"github.com/0xAX/notificator"

	"github.com/zserge/lorca"
)

const (
	WelcomeKey string = "welcome"
)

type Welcome struct{}

func (c *Welcome) Menu(n Notifier, r Renderer) {

	//	n.Push("test", "test", "", notificator.UR_CRITICAL)

	go func() {
		ui, err := lorca.New("https://www.mocaccino.org/", "", 480, 320)
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

func (c *Welcome) Close()     {}
func (c *Welcome) ID() string { return WelcomeKey }

func (c *Welcome) Render(Notifier) string {
	return ""
}
