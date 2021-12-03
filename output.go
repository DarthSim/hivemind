package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/term/termios"
)

type ptyPipe struct {
	pty, tty *os.File
}

type multiOutput struct {
	maxNameLength  int
	mutex          sync.Mutex
	pipes          map[*process]*ptyPipe
	printProcName  bool
	printTimestamp bool
}

func (m *multiOutput) openPipe(proc *process) (pipe *ptyPipe) {
	var err error

	pipe = m.pipes[proc]

	pipe.pty, pipe.tty, err = termios.Pty()
	fatalOnErr(err)

	proc.Stdout = pipe.tty
	proc.Stderr = pipe.tty
	proc.Stdin = pipe.tty
	proc.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	return
}

func (m *multiOutput) Connect(proc *process) {
	if len(proc.Name) > m.maxNameLength {
		m.maxNameLength = len(proc.Name)
	}

	if m.pipes == nil {
		m.pipes = make(map[*process]*ptyPipe)
	}

	m.pipes[proc] = &ptyPipe{}
}

func (m *multiOutput) PipeOutput(proc *process) {
	pipe := m.openPipe(proc)

	go func(proc *process, pipe *ptyPipe) {
		scanLines(pipe.pty, func(b []byte) bool {
			m.WriteLine(proc, b)
			return true
		})
	}(proc, pipe)
}

func (m *multiOutput) ClosePipe(proc *process) {
	if pipe := m.pipes[proc]; pipe != nil {
		pipe.pty.Close()
		pipe.tty.Close()
	}
}

func (m *multiOutput) WriteLine(proc *process, p []byte) {
	var buf bytes.Buffer

	if m.printProcName || m.printTimestamp {
		color := fmt.Sprintf("\033[1;38;5;%vm", proc.Color)

		buf.WriteString(color)

		if m.printTimestamp {
			buf.WriteString(time.Now().Format("15:04:05"))
			buf.WriteByte(' ')
		}

		if m.printProcName {
			buf.WriteString(proc.Name)

			for i := len(proc.Name); i <= m.maxNameLength; i++ {
				buf.WriteByte(' ')
			}
		}

		buf.WriteString("\033[0m| ")
	}

	buf.Write(p)
	buf.WriteByte('\n')

	m.mutex.Lock()
	defer m.mutex.Unlock()

	buf.WriteTo(os.Stdout)
}

func (m *multiOutput) WriteErr(proc *process, err error) {
	m.WriteLine(proc, []byte(
		fmt.Sprintf("\033[0;31m%v\033[0m", err),
	))
}
