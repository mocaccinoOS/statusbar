package uilibs

import (
	"fmt"

	"github.com/zserge/lorca"
)

type UITerminalWriter struct {
	UI   lorca.UI
	Chan chan string
	Done chan bool
}

func NewTerminalWriter(ui lorca.UI) *UITerminalWriter {
	return &UITerminalWriter{UI: ui, Chan: make(chan string), Done: make(chan bool)}
}

func (ui *UITerminalWriter) Start() {
	go func() {
		for {
			select {
			case m := <-ui.Chan:
				ui.WriteString(m)
			case <-ui.Done:
				return
			}
		}
	}()
}

func (ui *UITerminalWriter) Close() {
	ui.Done <- true
}

func (ui *UITerminalWriter) Write(p []byte) (n int, err error) {
	fmt.Println("Writing", string(p))

	ui.UI.Eval(fmt.Sprintf("console.log(`%s`)",
		string(p)))
	return len(p), nil
}

func (ui *UITerminalWriter) WriteString(s string) (int, error) {
	return ui.Write([]byte(s))
}
