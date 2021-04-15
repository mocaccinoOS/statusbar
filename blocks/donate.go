package blocks

import (
	"fmt"

	util "github.com/MocaccinoOS/statusbar/utils"
	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
)

const (
	DonateKey       string = "donate"
	continueMessage        = `
This will start xmrig locally to donate to the MocaccinoOS project wallet when 
your CPU will be detected idleing in a 15mins interval. 

Are you sure you want to continue? `
)

type Donate struct {
}

func (c *Donate) Close() {
	StopDonate()
	OnlyIdleDisable()
}

func (c *Donate) Menu(r Renderer) {
	donateMenu := systray.AddMenuItem("Donate HW Power", "Donate HW Power to MocaccinoOS by running xmrig")
	mStartDonating := donateMenu.AddSubMenuItem("Start donating HW", "Runs xmrig")
	mStopDonating := donateMenu.AddSubMenuItem("Stop donating HW", "Stops xmrig")
	onlyIdle := donateMenu.AddSubMenuItemCheckbox("Only when idleing", "Donate only when your CPU is idleing", false)
	mStopDonating.Hide()

	go func() {

		donateShown := true

		showDonate := func() {
			mStopDonating.Hide()
			mStartDonating.Show()
			donateShown = true
		}

		hideDonate := func() {
			mStopDonating.Show()
			mStartDonating.Hide()
			donateShown = false
		}

		startdonate := func() {
			if err := StartDonate(); err != nil {
				showDonate()
				dialog.Message(err.Error()).Title("Failed").Error()
			} else {
				hideDonate()
			}
		}
		stopdonate := func() {
			if err := StopDonate(); err != nil {
				hideDonate()
				dialog.Message(err.Error()).Title("Failed").Error()
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
				dialog.Message(err.Error()).Title("Failed").Error()
				onlyIdle.Check()
			} else {
				onlyIdle.Uncheck()
			}
		}

		enableIdleDaemon := func() {
			ok := dialog.Message("%s", continueMessage).Title("Are you sure?").YesNo()
			if ok {
				if err := OnlyIdleEnable(); err != nil {
					dialog.Message(err.Error()).Title("Failed").Error()
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

func (c *Donate) String() string {
	return fmt.Sprintf("")
}

func StartDonate() error {
	if err := OnlyIdleEnable(); err != nil {
		return err
	}
	out, err := util.Sudo("yip -s reconcile /etc/mocaccino/profiles && systemctl start xmrig")
	if err != nil {
		fmt.Println("Failed!")
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
