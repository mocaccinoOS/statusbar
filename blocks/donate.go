package blocks

import (
	"fmt"

	"github.com/0xAX/notificator"
	util "github.com/MocaccinoOS/statusbar/utils"
	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
)

const (
	DonateKey       = "donate"
	continueMessage = `
This will start xmrig locally to donate to the MocaccinoOS project wallet when 
your CPU will be detected idleing in a 15mins interval. 

Are you sure you want to continue? `
)

type Donate struct {
	donating bool
}

func (c *Donate) Close() {
	// Following requires interaction. TODO: find a way to disable without poking user
	// StopDonate()
	// OnlyIdleDisable()
}

func (c *Donate) Menu(n Notifier, r Renderer) {
	donateMenu := systray.AddMenuItem("Donate HW Power", "Donate HW Power to MocaccinoOS by running xmrig")
	mStartDonating := donateMenu.AddSubMenuItem("Start donating HW", "Runs xmrig")
	mStopDonating := donateMenu.AddSubMenuItem("Stop donating HW", "Stops xmrig")
	onlyIdle := donateMenu.AddSubMenuItemCheckbox("Only when idleing", "Donate only when your CPU is idleing", false)
	mStopDonating.Hide()

	go func() {

		donateShown := true

		showDonate := func() {
			//	c.donating = false
			mStopDonating.Hide()
			mStartDonating.Show()
			donateShown = true
		}

		hideDonate := func() {
			//	c.donating = true
			mStopDonating.Show()
			mStartDonating.Hide()
			donateShown = false
		}

		startdonate := func() {
			if err := StartDonate(); err != nil {
				showDonate()
				n.Push("Failed to start donating", err.Error(), "", notificator.UR_CRITICAL)
				//dialog.Message(err.Error()).Title("Failed").Error()
			} else {
				hideDonate()
			}
		}
		stopdonate := func() {
			if err := StopDonate(); err != nil {
				hideDonate()
				n.Push("Failed to stop donating", err.Error(), "", notificator.UR_CRITICAL)

				//	dialog.Message(err.Error()).Title("Failed").Error()
			} else {
				showDonate()
			}
		}

		donateToggle := func() {
			if donateShown {
				startdonate()
			} else {
				stopdonate()
			}
		}

		disableIdleDaemon := func() {
			if err := OnlyIdleDisable(); err != nil {
				n.Push("Failed to stop Idle service", err.Error(), "", notificator.UR_CRITICAL)
				onlyIdle.Check()
			} else {
				onlyIdle.Uncheck()
			}
		}

		enableIdleDaemon := func() {
			ok := dialog.Message("%s", continueMessage).Title("Are you sure?").YesNo()
			if ok {
				if err := OnlyIdleEnable(); err != nil {
					n.Push("Failed to enable Idle service", err.Error(), "", notificator.UR_CRITICAL)
					onlyIdle.Uncheck()
				} else {
					onlyIdle.Check()
				}
			} else {
				onlyIdle.Uncheck()
			}
		}

		for {
			select {
			case <-onlyIdle.ClickedCh:
				if onlyIdle.Checked() {
					disableIdleDaemon()
				} else {
					enableIdleDaemon()
				}
			case <-mStopDonating.ClickedCh:
				donateToggle()
			case <-mStartDonating.ClickedCh:
				donateToggle()
			}
		}
	}()
}

func (c *Donate) ID() string { return DonateKey }

func (c *Donate) Render(Notifier) string {
	if c.donating {
		return "🌟"
	}
	return ""
}

func StartDonate() error {
	if err := OnlyIdleEnable(); err != nil {
		return err
	}
	out, err := util.Sudo("yip -s reconcile /etc/mocaccino/profiles && systemctl start xmrig")
	if err != nil {
		fmt.Println("Failed starting to donate!")
	}
	fmt.Println(out)
	return err
}

func StopDonate() error {
	if err := OnlyIdleDisable(); err != nil {
		return err
	}
	out, err := util.Sudo("systemctl stop xmrig")
	if err != nil {
		fmt.Println("Failed!")
	}
	fmt.Println(out)
	return err
}

func OnlyIdleEnable() error {
	out, err := util.Sudo("luet install -y system-profile/donate")
	if err != nil {
		fmt.Println("Failed!")
	}
	fmt.Println(out)
	return err
}
func OnlyIdleDisable() error {
	out, err := util.Sudo("luet uninstall -y system-profile/donate")
	fmt.Println(out)
	if err != nil {
		fmt.Println("Failed!")
	}
	return err
}
