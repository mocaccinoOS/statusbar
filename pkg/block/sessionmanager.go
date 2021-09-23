package block

import (
	"log"
	"os"
	"path/filepath"

	"github.com/MocaccinoOS/statusbar/pkg/uilibs"
	process "github.com/mudler/go-processmanager"
	"github.com/nxadm/tail"
)

type SessionManager struct {
	Application string
}

func (sm *SessionManager) Path(s ...string) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	d := []string{dirname, sm.Application}
	d = append(d, s...)

	return filepath.Join(d...)
}

func (sm *SessionManager) Process(name, command string) *process.Process {
	return process.New(
		process.WithName("/bin/bash"),
		process.WithArgs("-c", command),
		process.WithStateDir(sm.Path("process", name)),
	)
}

func (sm *SessionManager) AttachLogfiles(pr *process.Process, w *uilibs.UITerminalWriter, doneStdout, doneStderr chan bool) {
	go func() {
		t, _ := tail.TailFile(pr.StdoutPath(), tail.Config{Follow: true})

		for {
			select {
			case line := <-t.Lines: // Every 100ms increate number of ticks and update UI
				w.Chan <- line.Text
			case <-doneStdout:
				return
			}
		}

	}()
	go func() {
		t, _ := tail.TailFile(pr.StderrPath(), tail.Config{Follow: true})
		for {
			select {
			case line := <-t.Lines: // Every 100ms increate number of ticks and update UI
				w.Chan <- line.Text
			case <-doneStderr:
				return
			}
		}
	}()
}
