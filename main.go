//go:generate statik -f -src=./public -include=*.woff2,*.woff,*.ttf,*.png,*.jpg,*.txt,*.html,*.css,*.js

package main

import (
	"github.com/MocaccinoOS/statusbar/blocks"
	"github.com/MocaccinoOS/statusbar/pkg/statusbar"

	"github.com/getlantern/systray"
)

var Version string

func main() {
	link := systray.AddMenuItem("Links", "")
	systray.AddSeparator()
	metrics := systray.AddMenuItem("Metrics", "")
	bar := statusbar.NewBar(
		statusbar.WithAppName("MocaccinoOS Statusbar"),
		//	statusbar.WithNotificationIcon("icon/icon.png"),
		statusbar.WithBlocks(
			//&blocks.Upgrade{},

			//	&blocks.Settings{},
			&blocks.Open{
				Text:    "Community",
				SubText: "Browse Community",
				URL:     "https://community.mocaccino.org/",
				SubMenu: link,
			},
			&blocks.Open{
				Text:    "Issues",
				SubText: "File a new bug or a feature request",
				URL:     "https://github.com/mocaccinoOS/mocaccino/issues",
				SubMenu: link,
			},
			&blocks.Open{
				Text:    "Chat",
				SubText: "Join us on slack",
				URL:     "https://join.slack.com/t/luet/shared_invite/enQtOTQxMjcyNDQ0MDUxLWQ5ODVlNTI1MTYzNDRkYzkyYmM1YWE5YjM0NTliNDEzNmQwMTkxNDRhNDIzM2Y5NDBlOTZjZTYxYWQyNDE4YzY",
				SubMenu: link,
			},
			&blocks.Separator{},
			&blocks.ChromeEmbeddedOpener{
				Text:    "Documentation",
				SubText: "Browse MocaccinoOS Docs",
				URL:     "https://www.mocaccino.org/docs/",
				SubMenu: link,
			},
			&blocks.ChromeEmbeddedOpener{
				Text:    "Packages",
				SubText: "Browse MocaccinoOS Packages",
				URL:     "https://packages.mocaccino.org/",
				SubMenu: link,
			},
			&blocks.Separator{},
			&blocks.Welcome{},
			&blocks.CPU{
				SubMenu: metrics,
			},
			&blocks.Memory{
				SubMenu: metrics,
			},
			&blocks.ShellToggle{
				SubMenu: metrics,
				Name:    "Show Running processes", Prefix: "p", Command: "ps -e | wc -l"},
			//&blocks.Network{},

			&blocks.Separator{},
			&blocks.ShellMenu{
				SubMenu: metrics,
				Name:    "Running processes", Command: "ps -e | wc -l"},
			&blocks.NetworkStat{SubMenu: metrics},

			&blocks.Separator{},
			&blocks.Updater{},

			&blocks.About{Version: Version},

			&blocks.Separator{},
		),
	)
	systray.Run(bar.Ready, bar.Close)
}
