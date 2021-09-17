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
	sync.Mutex
}

func icon(icon string) string {
	return fmt.Sprintf("<i class=\"bi bi-%s\"></i>", icon)
}

func (c *Updater) message(ui lorca.UI, i ...interface{}) {
	ui.Eval(fmt.Sprintf("document.querySelector('.updates').innerHTML = `%s`",
		i...))
}

func (c *Updater) refreshUI(ui lorca.UI, sm SessionManager) {
	//c.Lock()
	//defer c.Unlock()
	pr := sm.Process("update", UpgradeCommand)
	if pr.IsAlive() {
		ui.Eval("$('#terminal').show()")
		c.runningUpgrade = true
		c.message(ui, icon("download"))
		ui.Eval("$('#upgrade').hide()")
		c.availableUpgrade = false
	} else if c.availableUpgrade {
		c.message(ui, icon("download"))
		ui.Eval("$('#upgrade').show()")

		c.runningUpgrade = false
	} else if c.checkingUpgrades {
		ui.Eval("$('#check').hide()")

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
	//cmd := exec.Command("/bin/bash", "-c", `LUET_NOLOCK=true luet upgrade`)
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
				ui, err := lorca.New("http://127.0.0.1:9910/updater/index.html", "", 600, 390)
				if err != nil {
					log.Println("Failed starting chrome", err.Error())
					return
				}

				w := uilibs.NewTerminalWriter(ui)
				w.Start()
				pr := sm.Process("update", UpgradeCommand)
				c.refreshUI(ui, sm)
				doneUI, doneUpgrades := make(chan bool), make(chan bool)

				// Bind Go functions to JS
				ui.Bind("check", func() {
					c.availableUpgrades()
				})

				ui.Bind("upgrade", func() {

					pr = sm.Process("update", UpgradeCommand)
					c.refreshUI(ui, sm)
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

					sm.AttachLogfiles(pr, w)
				})

				go func() {
					t := time.NewTicker(1 * time.Second)
					for {
						select {
						case <-t.C: // Every 100ms increate number of ticks and update UI
							c.refreshUI(ui, sm)
						case <-doneUI:
							return
						}
					}
				}()

				go func() {

					d := 5 * time.Minute
					if c.CheckUpgradesTimer != nil {
						d = *c.CheckUpgradesTimer
					}
					t := time.NewTicker(d)
					for {
						select {
						case <-t.C: // Every 100ms increate number of ticks and update UI
							c.availableUpgrades()
						case <-doneUpgrades:
							return
						}
					}
				}()
				//	ui.SetBounds(lorca.Bounds{WindowState: lorca.WindowStateMaximized})
				//	if err != nil {
				//		log.Fatal(err)
				//	}
				defer func() {
					ui.Close()
					w.Close()
					doneUI <- true
					doneUpgrades <- true
				}()

				c.availableUpgrades()
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
