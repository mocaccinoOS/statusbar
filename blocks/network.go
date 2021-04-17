package blocks

import (
	"fmt"
	util "github.com/MocaccinoOS/statusbar/utils"

	"github.com/c2h5oh/datasize"
	"github.com/getlantern/systray"
	"github.com/mackerelio/go-osstat/network"
)

const (
	NetworkKey string = "Network"
)

type Network struct{}

func (c *Network) Menu(n Notifier, r Renderer) {
	showNetwork := systray.AddMenuItemCheckbox("Show Network Statistics", "Show Network Statistics", false)

	go func() {
		for range showNetwork.ClickedCh {
			if showNetwork.Checked() {
				showNetwork.Uncheck()
				r.Disable(NetworkKey)
				systray.SetTitle("")
			} else {
				showNetwork.Check()
				r.Activate(NetworkKey)
			}
		}
	}()
}

func (c *Network) Close()     {}
func (c *Network) ID() string { return NetworkKey }

func (c *Network) Render(Notifier) string {
	stats, err := network.Get()
	if err != nil {
		return ""
	}
	str := ""
	for _, s := range stats {
		txBytes, _ := util.Run(fmt.Sprintf("S=10; F=/sys/class/net/%s/statistics/tx_bytes; X=`cat $F`; sleep $S; Y=`cat $F`; BPS=$(((Y-X)/S)); echo $BPS", s.Name))
		rxBytes, _ := util.Run(fmt.Sprintf("S=10; F=/sys/class/net/%s/statistics/rx_bytes; X=`cat $F`; sleep $S; Y=`cat $F`; BPS=$(((Y-X)/S)); echo $BPS", s.Name))

		var tx datasize.ByteSize
		err := tx.UnmarshalText([]byte(txBytes))
		if err != nil {
			continue
		}
		var rx datasize.ByteSize
		err = rx.UnmarshalText([]byte(rxBytes))
		if err != nil {
			continue
		}
		str += fmt.Sprintf("%s: %s %s", s.Name,
			tx.HumanReadable(),
			rx.HumanReadable())
	}
	return str
}
