package statusbar

import (
	"fmt"

	"github.com/MocaccinoOS/statusbar/pkg/block"
	"github.com/getlantern/systray"
)

func (sb *Statusbar) renderSystray(br *block.Renderer) {
	br.RenderMenu()
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		defer br.Close()

		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}
