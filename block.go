package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MocaccinoOS/statusbar/blocks"
	"github.com/getlantern/systray"
)

type Block interface {
	String() string
	ID() string
	Menu(blocks.Renderer)
	Close()
}

type BlockRenderer struct {
	sync.Mutex
	blocks []Block
	active map[string]interface{}
}

func (br *BlockRenderer) RenderMenu() {
	for _, b := range br.blocks {
		b.Menu(br)
	}
}

func (br *BlockRenderer) Render() string {
	br.Lock()
	defer br.Unlock()
	var res string
	for _, b := range br.blocks {
		if _, ok := br.active[b.ID()]; ok {
			res += b.String() + " "
		}
	}

	return strings.TrimSpace(res)
}

func (br *BlockRenderer) Disable(s string) {
	br.Lock()
	defer br.Unlock()
	fmt.Println("DeActivating", s)

	delete(br.active, s)
}

func (br *BlockRenderer) Activate(s string) {
	br.Lock()
	defer br.Unlock()
	fmt.Println("Activating", s)
	br.active[s] = nil
}

func (br *BlockRenderer) Run() {
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				systray.SetTitle(br.Render())
			}
		}
	}()
}

func (br *BlockRenderer) Close() {
	for _, b := range br.blocks {
		b.Close()
	}
}

func Renderer() *BlockRenderer {
	return &BlockRenderer{active: make(map[string]interface{}), blocks: []Block{
		&blocks.Donate{},
		&blocks.CPU{},
		&blocks.Memory{},
	}}
}
