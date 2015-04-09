package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/kr/pty"
)

type Multiterm struct {
	maxNameLength int
	mutex         sync.Mutex
}

var multiterm = Multiterm{}

func (m *Multiterm) PipeOutput(proc *Process) {
	if len(proc.Name) > m.maxNameLength {
		m.maxNameLength = len(proc.Name)
	}

	pty, tty, err := pty.Open()
	if err != nil {
		log.Fatal(err)
	}

	proc.Stdout = tty
	proc.Stderr = tty
	proc.Stdin = tty
	proc.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	go func(proc *Process, pty *os.File) {
		var buf bytes.Buffer

		reader := bufio.NewReader(pty)

		for {
			b, err := reader.ReadBytes('\n')
			buf.Write(b)

			if err == nil {
				m.FlushToStdout(proc, &buf)
			}
		}
	}(proc, pty)
}

func (m *Multiterm) FlushToStdout(proc *Process, buf *bytes.Buffer) {
	color := fmt.Sprintf("\033[0;%vm", proc.Color)

	var nameBuf bytes.Buffer

	nameBuf.WriteString(color)
	nameBuf.WriteString(proc.Name)

	for nameBuf.Len()-len(color) < m.maxNameLength {
		nameBuf.Write([]byte{' '})
	}

	nameBuf.WriteString("\033[0m | ")

	m.mutex.Lock()
	defer m.mutex.Unlock()

	nameBuf.WriteTo(os.Stdout)
	buf.WriteTo(os.Stdout)
}
