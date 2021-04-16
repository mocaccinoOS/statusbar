//go:generate statik -f -src=./public -include=*.jpg,*.txt,*.html,*.css,*.js

package main

import (
	"github.com/MocaccinoOS/statusbar/blocks"
	"github.com/MocaccinoOS/statusbar/pkg/statusbar"

	"github.com/getlantern/systray"
)

func main() {
	bar := statusbar.NewBar(
		statusbar.WithAppName("MocaccinoOS Statusbar"),
		//	statusbar.WithNotificationIcon("icon/icon.png"),
		statusbar.WithBlocks(
			&blocks.Upgrade{},

	//		&blocks.Settings{},
			&blocks.Open{
				Text:    "Community",
				SubText: "Browse Community",
				URL:     "https://community.mocaccino.org/",
			},
			&blocks.Open{
				Text:    "Issues",
				SubText: "File a new bug or a feature request",
				URL:     "https://github.com/mocaccinoOS/mocaccino/issues",
			},
			&blocks.Open{
				Text:    "Chat",
				SubText: "Join us on slack",
				URL:     "https://join.slack.com/t/luet/shared_invite/enQtOTQxMjcyNDQ0MDUxLWQ5ODVlNTI1MTYzNDRkYzkyYmM1YWE5YjM0NTliNDEzNmQwMTkxNDRhNDIzM2Y5NDBlOTZjZTYxYWQyNDE4YzY",
			},
			&blocks.Separator{},
			&blocks.ChromeEmbeddedOpener{
				Text:    "Documentation",
				SubText: "Browse MocaccinoOS Docs",
				URL:     "https://www.mocaccino.org/docs/",
			},
			&blocks.ChromeEmbeddedOpener{
				Text:    "Packages",
				SubText: "Browse MocaccinoOS Packages",
				URL:     "https://packages.mocaccino.org/",
			},
			&blocks.Separator{},
			&blocks.Welcome{},
			&blocks.Donate{},
			&blocks.CPU{},
			&blocks.Memory{},
		),
	)
	systray.Run(bar.Ready, bar.Close)
}
