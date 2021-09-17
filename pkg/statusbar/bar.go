package statusbar

import (
	"log"
	"net/http"

	"github.com/0xAX/notificator"
	"github.com/MocaccinoOS/statusbar/icon"
	"github.com/MocaccinoOS/statusbar/pkg/block"
	_ "github.com/MocaccinoOS/statusbar/statik"
	"github.com/getlantern/systray"
	"github.com/rakyll/statik/fs"
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

func (sb *Statusbar) Serve() {
	go func() {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}

		// Serve the contents over HTTP.
		http.Handle("/", http.FileServer(statikFS))
		http.ListenAndServe("127.0.0.1:9910", nil)
	}()
}

func (sb *Statusbar) Ready() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle(sb.Options.Title)
	systray.SetTooltip(sb.Options.Tooltip)

	sb.notificator = notificator.New(notificator.Options{
		DefaultIcon: sb.Options.NotificatorIcon,
		AppName:     sb.Options.NotificatorAppName,
	})

	br := block.NewRenderer(sb.notificator, sb.Options.Blocks, &block.SessionManager{Application: sb.Options.NotificatorAppName})
	br.Run()
	sb.Serve()

	sb.renderSystray(br)
}

func (sb *Statusbar) Close() {
	//	now := time.Now()
	//ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
}
