package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/DarthSim/hivemind/_third_party/github.com/kr/pty"
)

type PtyPipe struct {
	pty, tty *os.File
}

type Multiterm struct {
	maxNameLength int
	mutex         sync.Mutex
	pipes         map[*Process]*PtyPipe
}

func (m *Multiterm) openPipe(proc *Process) (pipe *PtyPipe) {
	var err error

	pipe = m.pipes[proc]

	pipe.pty, pipe.tty, err = pty.Open()
	fatalOnErr(err)

	proc.Stdout = pipe.tty
	proc.Stderr = pipe.tty
	proc.Stdin = pipe.tty
	proc.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	return
}

func (m *Multiterm) Connect(proc *Process) {
	if len(proc.Name) > m.maxNameLength {
		m.maxNameLength = len(proc.Name)
	}

	if m.pipes == nil {
		m.pipes = make(map[*Process]*PtyPipe)
	}

	m.pipes[proc] = &PtyPipe{}
}

func (m *Multiterm) PipeOutput(proc *Process) {
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
