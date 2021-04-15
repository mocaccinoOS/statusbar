package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/MocaccinoOS/statusbar/icon"
	"github.com/getlantern/systray"
)

func main() {
	onExit := func() {
		now := time.Now()
		ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("")
	systray.SetTooltip("Mocaccino OS Statusbar")

	br := Renderer()
	br.Run()

	Systray(br)
}
