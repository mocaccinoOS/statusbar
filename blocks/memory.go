package blocks

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/mackerelio/go-osstat/memory"
)

const (
	MemoryKey string = "memory"
)

type Memory struct{}

func (c *Memory) Menu(n Notifier, r Renderer) {
	showMemory := systray.AddMenuItemCheckbox("Show Memory Metrics", "Show Memory metrics", false)
	go func() {
		for range showMemory.ClickedCh {
			if showMemory.Checked() {
				showMemory.Uncheck()
				r.Disable(MemoryKey)
				systray.SetTitle("")
			} else {
				showMemory.Check()
				r.Activate(MemoryKey)
			}
		}
	}()
}

func (c *Memory) Close()     {}
func (c *Memory) ID() string { return MemoryKey }

func (c *Memory) Render(Notifier) string {
	now, err := memory.Get()
	if err != nil {
		return ""
	}
	totram := (now.Used * 100) / now.Total

	icon := ""
	switch {
	case totram > 80:
		icon = "☢"
	case totram > 90:
		icon = "☠"
	case totram < 10:
		icon = "🐢"
	case totram > 50:
		icon = "🚴"
	}

	return fmt.Sprintf(" Memory %d%% %s", totram, icon)
}
