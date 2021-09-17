package blocks

import (
	"fmt"

	util "github.com/MocaccinoOS/statusbar/utils"
	"github.com/getlantern/systray"
)

const (
	ShellKey string = "Shell"
)

type ShellToggle struct {
	Name    string
	Prefix  string
	Command string
	SubMenu *systray.MenuItem
}

func (c *ShellToggle) Menu(n Notifier, r Renderer, sm SessionManager) {
	var showCommand *systray.MenuItem
	if c.SubMenu != nil {
		showCommand = c.SubMenu.AddSubMenuItemCheckbox(c.Name, "", false)
	} else {
		showCommand = systray.AddMenuItemCheckbox(c.Name, "", false)
	}
	go func() {
		for range showCommand.ClickedCh {
			if showCommand.Checked() {
				showCommand.Uncheck()
				r.Disable(ShellKey)
				//systray.SetTitle("")
			} else {
				showCommand.Check()
				r.Activate(ShellKey)
			}
		}
	}()
}

func (c *ShellToggle) Close()     {}
func (c *ShellToggle) ID() string { return ShellKey }

func (c *ShellToggle) Render(Notifier) string {
	out, _ := util.Run(c.Command)
	return fmt.Sprintf("%s: %s", c.Prefix, out)
}
