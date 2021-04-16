package blocks

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/mackerelio/go-osstat/cpu"
)

const (
	CPUKey string = "cpu"
)

type CPU struct {
	c *cpu.Stats
}

func (c *CPU) Close() {
}

func (c *CPU) Menu(n Notifier, r Renderer) {
	showCPU := systray.AddMenuItemCheckbox("Show CPU Metrics", "Show CPU metrics", false)

	go func() {
		for range showCPU.ClickedCh {
			if showCPU.Checked() {
				showCPU.Uncheck()
				r.Disable(CPUKey)
			} else {
				showCPU.Check()
				r.Activate(CPUKey)
			}
		}
	}()
}

func (c *CPU) ID() string { return CPUKey }
func (c *CPU) Render(Notifier) string {
	now, err := cpu.Get()
	if err != nil {
		return ""
	}
	if c.c == nil {
		c.c = now
		return ""
	}

	total := float64(now.Total - c.c.Total)
	userCpu := float64(now.User-c.c.User) / total * 100
	systemCpu := float64(now.System-c.c.System) / total * 100
	//	idleCpu := float64(now.Idle-before.Idle) / total * 100
	totCpu := (userCpu + systemCpu) / 2

	icon := ""
	switch {
	case totCpu > 80:
		icon = "☢"
	case totCpu > 90:
		icon = "☠"
	case totCpu < 10:
		icon = "🐢"
	case totCpu > 50:
		icon = "🚴"
	}

	c.c = now

	return fmt.Sprintf(" CPU %1.f%% %s", totCpu, icon)
}
