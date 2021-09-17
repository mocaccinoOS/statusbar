package blocks

import (
	"fmt"

	"time"

	"github.com/c2h5oh/datasize"
	"github.com/getlantern/systray"
	"github.com/mackerelio/go-osstat/network"
)

const (
	NetworkStatKey string = "NetworkStat"
)

type NetworkStat struct {
	SubMenu *systray.MenuItem
}

func (c *NetworkStat) Menu(n Notifier, r Renderer, sm SessionManager) {
	devices := map[string]*systray.MenuItem{}

	var networkStat *systray.MenuItem
	if c.SubMenu != nil {
		networkStat = c.SubMenu.AddSubMenuItem("Network statistics", "")
	} else {
		networkStat = systray.AddMenuItem("Network statistics", "")
	}

	go func() {
		for {
			stats, err := network.Get()
			if err != nil {
				continue
			}
			for _, s := range stats {
				txt := fmt.Sprintf("%s: Rx %s Tx %s", s.Name,
					(datasize.ByteSize(s.RxBytes) * datasize.B).HumanReadable(),
					(datasize.ByteSize(s.TxBytes) * datasize.B).HumanReadable())
				if _, ok := devices[s.Name]; !ok {
					devices[s.Name] = networkStat.AddSubMenuItem(txt, "")
					devices[s.Name].Disable()
				} else {
					devices[s.Name].SetTitle(txt)
				}
			}
			time.Sleep(10 * time.Second)

		}
	}()
}

func (c *NetworkStat) Close()     {}
func (c *NetworkStat) ID() string { return NetworkStatKey }

func (c *NetworkStat) Render(Notifier) string {

	return ""
}
