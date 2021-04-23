package blocks

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/MocaccinoOS/statusbar/statik"
	util "github.com/MocaccinoOS/statusbar/utils"
	"github.com/zserge/lorca"
	//	"github.com/0xAX/notificator"
)

const (
	WelcomeKey string = "welcome"
)

type Welcome struct{ URL string }

func (c *Welcome) startWelcome(sentinel string) {

	go func() {
		ui, err := lorca.New("http://127.0.0.1:9910/welcome/index.html", "", 480, 320)
		if err != nil {
			log.Println("Failed starting chrome", err.Error())
			return
		}
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
		if sentinel != "" {
			util.Touch(sentinel)
		}
	}()
}

func (c *Welcome) Menu(n Notifier, r Renderer) {

	var welcomeSentinel string

	dirname, err := os.UserHomeDir()
	if err == nil {
		welcomeSentinel = filepath.Join(dirname, ".config", "mocaccino-statusbar", "welcome_displayed")
	}

	if _, err := os.Stat(welcomeSentinel); os.IsNotExist(err) {
		c.startWelcome(welcomeSentinel)
	}
}

func (c *Welcome) Close()     {}
func (c *Welcome) ID() string { return WelcomeKey }

func (c *Welcome) Render(Notifier) string {
	return ""
}
