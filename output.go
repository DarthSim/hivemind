package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/kr/pty"
)

type PtyPipe struct {
	pty, tty *os.File
}

type Multiterm struct {
	maxNameLength int
	mutex         sync.Mutex
	pipes         map[*Process]*PtyPipe
}

var multiterm = Multiterm{}

func (m *Multiterm) openPipe(proc *Process) (pipe *PtyPipe) {
	pty, tty, err := pty.Open()
	fatalOnErr(err)

	pipe = &PtyPipe{pty, tty}

	if m.pipes == nil {
		m.pipes = make(map[*Process]*PtyPipe)
	}

	m.pipes[proc] = pipe

	proc.Stdout = tty
	proc.Stderr = tty
	proc.Stdin = tty
	proc.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	return
}

func (m *Multiterm) PipeOutput(proc *Process) {
	if len(proc.Name) > m.maxNameLength {
		m.maxNameLength = len(proc.Name)
	}

	pipe := m.openPipe(proc)

	go func(proc *Process, pipe *PtyPipe) {
		scanner := bufio.NewScanner(pipe.pty)

		for scanner.Scan() {
			m.WriteLine(proc, scanner.Bytes())
		}
	}(proc, pipe)
}

func (m *Multiterm) ClosePipe(proc *Process) {
	if pipe := m.pipes[proc]; pipe != nil {
		pipe.pty.Close()
		pipe.tty.Close()
	}
}

func (m *Multiterm) WriteLine(proc *Process, p []byte) {
	var buf bytes.Buffer

	color := fmt.Sprintf("\033[1;%vm", proc.Color)

	buf.WriteString(color)
	buf.WriteString(proc.Name)

	for buf.Len()-len(color) < m.maxNameLength {
		buf.WriteByte(' ')
	}

	buf.WriteString("\033[0m | ")
	buf.Write(p)
	buf.WriteByte('\n')

	m.mutex.Lock()
	defer m.mutex.Unlock()

	buf.WriteTo(os.Stdout)
}

func (m *Multiterm) WriteErr(proc *Process, err error) {
	m.WriteLine(proc, []byte(
		fmt.Sprintf("\033[1;31m%v\033[0m", err),
	))
}
