package block

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MocaccinoOS/statusbar/blocks"
	"github.com/getlantern/systray"
)

type Notifier interface {
	Push(string, string, string, string) error
}

type Block interface {
	Render(blocks.Notifier) string
	ID() string
	Menu(blocks.Notifier, blocks.Renderer)
	Close()
}

type Renderer struct {
	sync.Mutex
	notifier blocks.Notifier
	blocks   []Block
	active   map[string]interface{}
}

func (br *Renderer) RenderMenu() {
	for _, b := range br.blocks {
		b.Menu(br.notifier, br)
	}
}

func (br *Renderer) Render() string {
	br.Lock()
	defer br.Unlock()
	var res string
	for _, b := range br.blocks {
		if _, ok := br.active[b.ID()]; ok {
			res += b.Render(br.notifier) + " "
		}
	}

	return strings.TrimSpace(res)
}

func (br *Renderer) Disable(s string) {
	br.Lock()
	defer br.Unlock()
	fmt.Println("DeActivating", s)

	delete(br.active, s)
}

func (br *Renderer) Activate(s string) {
	br.Lock()
	defer br.Unlock()
	fmt.Println("Activating", s)
	br.active[s] = nil
}

func (br *Renderer) Run() {
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for range ticker.C {
			systray.SetTitle(br.Render())
		}
	}()
}

func (br *Renderer) Close() {
	for _, b := range br.blocks {
		b.Close()
	}
}

func NewRenderer(n blocks.Notifier, b []Block) *Renderer {
	return &Renderer{
		notifier: n,
		active:   make(map[string]interface{}),
		blocks:   b,
	}
}
