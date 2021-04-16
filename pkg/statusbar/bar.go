package statusbar

import (
	"github.com/0xAX/notificator"
	"github.com/MocaccinoOS/statusbar/icon"
	"github.com/MocaccinoOS/statusbar/pkg/block"
	"github.com/getlantern/systray"
)

type Statusbar struct {
	Options     Options
	notificator *notificator.Notificator
}

func NewBar(p ...Option) *Statusbar {
	c := &Options{}
	c.Apply(p...)

	return &Statusbar{Options: *c}
}

func (sb *Statusbar) Ready() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle(sb.Options.Title)
	systray.SetTooltip(sb.Options.Tooltip)

	sb.notificator = notificator.New(notificator.Options{
		DefaultIcon: sb.Options.NotificatorIcon,
		AppName:     sb.Options.NotificatorAppName,
	})

	br := block.NewRenderer(sb.notificator, sb.Options.Blocks)
	br.Run()

	sb.renderSystray(br)
}

func (sb *Statusbar) Close() {
	//	now := time.Now()
	//ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}
