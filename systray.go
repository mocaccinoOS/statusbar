package main

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

func Systray(br *BlockRenderer) {

	communityURL := systray.AddMenuItem("Community", "Opens https://community.mocaccino.org")
	issuesURL := systray.AddMenuItem("Issues", "File a new bug or a feature request")

	systray.AddSeparator()

	br.RenderMenu()

	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		defer br.Close()

		for {
			select {
			case <-communityURL.ClickedCh:
				open.Run("https://community.mocaccino.org")
			case <-issuesURL.ClickedCh:
				open.Run("https://github.com/mocaccinoOS/mocaccino/issues")
			case <-mQuit.ClickedCh:
				systray.Quit()
				fmt.Println("Quit2 now...")
				return
			}
		}
	}()
}
