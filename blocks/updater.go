package blocks

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/MocaccinoOS/statusbar/pkg/uilibs"
	_ "github.com/MocaccinoOS/statusbar/statik"
	"github.com/getlantern/systray"
	"github.com/zserge/lorca"
	//	"github.com/0xAX/notificator"
)

const (
	UpdaterKey string = "Updater"
)

type Updater struct {
	URL string

	CheckUpgradesTimer *time.Duration
	runningUpgrade     bool
	availableUpgrade   bool
	checkingUpgrades   bool
	logfileAttached    bool
	sync.Mutex
}

func icon(icon string) string {
	return fmt.Sprintf("<i class=\"bi bi-%s\"></i>", icon)
}

func (c *Updater) message(ui lorca.UI, i ...interface{}) {
	ui.Eval(fmt.Sprintf("document.querySelector('.updates').innerHTML = `%s`",
		i...))
}

func (c *Updater) box(ui lorca.UI, i ...interface{}) {
	ui.Eval(fmt.Sprintf("document.querySelector('.box').innerHTML = `%s`",
		i...))
}

func (c *Updater) refreshUI(ui lorca.UI, sm SessionManager, w *uilibs.UITerminalWriter,
	stdout, stderr chan bool) {
	//c.Lock()
	//defer c.Unlock()
	pr := sm.Process("update", UpgradeCommand)
	if pr.IsAlive() {
		ui.Eval("$('#terminal').show()")
		c.runningUpgrade = true
		c.message(ui, "upgrades running"+icon("download"))
		c.box(ui, "Upgrades running, don't turn off your computer!")
		ui.Eval("$('#upgrade').hide()")
		ui.Eval("$('#check').hide()")
		c.availableUpgrade = false
		if !c.logfileAttached {
			sm.AttachLogfiles(pr, w, stdout, stderr)
			c.logfileAttached = true
		}
	} else if c.availableUpgrade {
		c.message(ui, icon("download"))
		ui.Eval("$('#upgrade').show()")
		c.box(ui, "Upgrades available! click on the Upgrade button to start the upgrade process. Don't turn off your computer during the process")
		c.runningUpgrade = false
	} else if c.checkingUpgrades {
		ui.Eval("$('#check').hide()")
		c.box(ui, "Please wait, checking available upgrades.")
		//ui.Eval(`$("#check i").attr("class","bi bi-arrow-clockwise");`)
		c.message(ui, `<div class="d-flex justify-content-center">
		<div class="spinner-border" role="status">
		  <span class="visually-hidden">Checking upgrades...</span>
		</div>
	  </div>
	  `)
	} else {
		//ui.Eval(`$("#check i").attr("class","bi bi-search");`)
		ui.Eval("$('#check').show()")
		c.box(ui, "All seems upgraded and running the latest software. You can close the window")

		fmt.Println("Upgrade not running")
		ui.Eval("$('#terminal').hide()")
		os.RemoveAll(pr.StateDir())

		c.message(ui, "All done "+icon("check-circle-fill"))
		c.runningUpgrade = false
	}

	//_, err :=ioutil.Stat(filepath.Join(pr.StateDir(),"error")); ioutil.
}

func (c *Updater) availableUpgrades() bool {
	//c.Lock()
	c.checkingUpgrades = true

	defer func() {
		c.checkingUpgrades = false
		//	c.Unlock()
	}()

	cmd := exec.Command("pkexec", "/bin/bash", "-c", `luet upgrade`)
	b, _ := cmd.CombinedOutput()
	if strings.Contains(string(b), "Nothing to do") {
		fmt.Println("Nothing to do")
		c.availableUpgrade = false
		return false
	} else {
		fmt.Println("upgrades available")
		c.availableUpgrade = true
		return true
	}
}

const UpgradeCommand = "pkexec /bin/bash -c 'luet upgrade --no-spinner --color=false -y'"

func (c *Updater) Menu(n Notifier, r Renderer, sm SessionManager) {
	url := systray.AddMenuItem("Software updates", "")
	//n.Push("test", "test", "", notificator.UR_CRITICAL)
	go func() {

		for range url.ClickedCh {
			go func() {
				c.logfileAttached = false
				ui, err := lorca.New("http://127.0.0.1:9910/updater/index.html", "", 600, 390)
				if err != nil {
					log.Println("Failed starting chrome", err.Error())
					return
				}
				doneUI, doneUpgrades, doneStdout, doneStderr := make(chan bool), make(chan bool), make(chan bool), make(chan bool)

				w := uilibs.NewTerminalWriter(ui)
				w.Start()
				pr := sm.Process("update", UpgradeCommand)
				c.refreshUI(ui, sm, w, doneStdout, doneStderr)

				c.availableUpgrades()

				// Bind Go functions to JS
				ui.Bind("check", func() {
					c.availableUpgrades()
				})

				ui.Bind("upgrade", func() {

					pr = sm.Process("update", UpgradeCommand)
					c.refreshUI(ui, sm, w, doneStdout, doneStderr)
					os.RemoveAll(pr.StateDir())

					c.message(ui, "Upgrade in progress")
					ui.Eval("$('#terminal').show()")
					ui.Eval("$('#upgrade').hide()")

					fmt.Println("Start upgrade", pr.StateDir())

					err := pr.Run()
					if err != nil {
						c.message(ui,
							icon("cancel"))

					}
					fmt.Println("attach logs", pr.StateDir())

					sm.AttachLogfiles(pr, w, doneStdout, doneStderr)
					c.logfileAttached = true
				})

				go func() {
					t := time.NewTicker(1 * time.Second)
					for {
						select {
						case <-t.C: // Every 100ms increate number of ticks and update UI
							c.refreshUI(ui, sm, w, doneStdout, doneStderr)
						case <-doneUI:
							return
						}
					}
				}()

				// go func() {

				// 	d := 5 * time.Minute
				// 	if c.CheckUpgradesTimer != nil {
				// 		d = *c.CheckUpgradesTimer
				// 	}
				// 	t := time.NewTicker(d)
				// 	for {
				// 		select {
				// 		case <-t.C: // Every 100ms increate number of ticks and update UI
				// 			c.availableUpgrades()
				// 		case <-doneUpgrades:
				// 			return
				// 		}
				// 	}
				// }()
				//	ui.SetBounds(lorca.Bounds{WindowState: lorca.WindowStateMaximized})
				//	if err != nil {
				//		log.Fatal(err)
				//	}
				defer func() {
					ui.Close()
					w.Close()
					doneUI <- true
					doneUpgrades <- true
					doneStderr <- true
					doneStdout <- true
				}()

				// Wait until UI window is closed
				<-ui.Done()

			}()
		}

	}()
}

func (c *Updater) Close()     {}
func (c *Updater) ID() string { return UpdaterKey }

func (c *Updater) Render(Notifier) string {
	if c.runningUpgrade {
		return "🚀 Upgrade runnng"
	}
	if c.availableUpgrade {
		return "🚀 Upgrade available"
	}
	return ""
}
