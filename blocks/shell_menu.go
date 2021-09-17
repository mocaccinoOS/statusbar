package blocks

import (
	"fmt"
	"time"

	util "github.com/MocaccinoOS/statusbar/utils"
	"github.com/getlantern/systray"
)

const (
	ShellMenuKey string = "Shell"
)

type ShellMenu struct {
	Name    string
	Command string
	SubMenu *systray.MenuItem
}

func (c *ShellMenu) Menu(n Notifier, r Renderer, sm SessionManager) {
	var showCommand *systray.MenuItem
	if c.SubMenu != nil {
		showCommand = c.SubMenu.AddSubMenuItem("", "")
	} else {
		showCommand = systray.AddMenuItem("", "")
	}
	showCommand.Disable()
	go func() {
		for {
			out, _ := util.Run(c.Command)
			showCommand.SetTitle(fmt.Sprintf("%s: %s", c.Name, out))
			//	systray.SetTitle("")
			time.Sleep(5 * time.Second)
		}
	}()
}

func (c *ShellMenu) Close()     {}
func (c *ShellMenu) ID() string { return ShellKey }

func (c *ShellMenu) Render(Notifier) string {
	return ""
}
