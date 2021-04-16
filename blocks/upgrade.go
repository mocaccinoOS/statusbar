package blocks

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/0xAX/notificator"
	util "github.com/MocaccinoOS/statusbar/utils"
	"github.com/acarl005/stripansi"
	"github.com/briandowns/spinner"
	"github.com/getlantern/systray"
)

const (
	UpgradeKey   string = "Upgrade"
	SentinelFile string = "/run/mocaccino/upgrades_available"
)

type Upgrade struct {
	sync.Mutex
	upgradesAvailable  bool
	notificationShowed bool
	upgradesRunning    bool

	spinnerChar string
}

func (c *Upgrade) Write(b []byte) (int, error) {
	c.spinnerChar = stripansi.Strip(string(b))
	return len(b), nil
}

func (c *Upgrade) bind(m *systray.MenuItem, n Notifier) {
	go func() {
		for range m.ClickedCh {
			someSet := []string{"🌕", "🌔", "🌓", "🌑"}
			s := spinner.New(someSet, 100*time.Millisecond, spinner.WithWriter(c))
			s.Start()

			m.Disable()
			m.SetTitle("⚡Upgrades running, hold on..")
			n.Push("Upgrade starting",
				"Sit tight, upgrades are running in background. Don't reboot your computer or halt your system",
				"",
				notificator.UR_NORMAL,
			)
			c.Lock()
			c.upgradesRunning = true
			c.upgradesAvailable = false
			c.Unlock()
			if out, err := util.Sudo(fmt.Sprintf("luet upgrade -y && rm -rf %s", SentinelFile)); err == nil {
				m.Hide()
				n.Push("Upgrade done",
					"All done, sit back and relax now",
					"",
					notificator.UR_NORMAL,
				)
				c.Lock()
				c.upgradesAvailable = false
				c.Unlock()
			} else {
				n.Push("Upgrade failed",
					err.Error()+string(out),
					"",
					notificator.UR_CRITICAL,
				)
				m.Enable()
				m.SetTitle("🚀  Upgrade")
			}
			c.Lock()
			c.upgradesRunning = false
			c.notificationShowed = false
			c.Unlock()

			s.Stop()
		}
	}()
}

func (c *Upgrade) onAvailableUpgrades(n Notifier) {
	// upgrades available
	c.Lock()
	c.upgradesAvailable = true
	c.Unlock()
	if !c.notificationShowed {
		n.Push("Available upgrades!",
			"New package upgrades are available for the system",
			"",
			notificator.UR_NORMAL,
		)
		c.notificationShowed = true
		upgrade := systray.AddMenuItem("🚀  Upgrade", "Upgrade to latest available software")
		c.bind(upgrade, n)
	}
}

func (c *Upgrade) onNotAvailable(n Notifier) {
	c.Lock()
	c.upgradesAvailable = false
	c.notificationShowed = false
	c.Unlock()
	// no upgrades available, nothing to show
}

func (c *Upgrade) Menu(n Notifier, r Renderer) {
	r.Activate(UpgradeKey)
	go func() {
		for {
			if _, err := os.Stat(SentinelFile); os.IsNotExist(err) {
				c.onNotAvailable(n)
			} else if err == nil {
				c.onAvailableUpgrades(n)
			}
			time.Sleep(60 * time.Second)
		}
	}()
}

func (c *Upgrade) Close()     {}
func (c *Upgrade) ID() string { return UpgradeKey }

func (c *Upgrade) Render(Notifier) string {
	icon := ""
	c.Lock()
	switch {
	case c.upgradesAvailable:
		icon = "Upgrades available 🚀"
	case c.upgradesRunning:
		icon = c.spinnerChar // "Upgrade in progress ⌛⏳🌕🌔🌓🌑"
	}
	c.Unlock()

	return icon
}
