package blocks

import (
	"github.com/MocaccinoOS/statusbar/pkg/uilibs"
	process "github.com/mudler/go-processmanager"
)

type Renderer interface {
	Activate(string)
	Disable(string)
}
type Notifier interface {
	Push(string, string, string, string) error
}

type SessionManager interface {
	Path(s ...string) string
	Process(name, command string) *process.Process
	AttachLogfiles(pr *process.Process, w *uilibs.UITerminalWriter, d, e chan bool)
}
