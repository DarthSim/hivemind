package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"sync"
	"syscall"
)

type Process struct {
	*exec.Cmd

	Name  string
	Color int
}

func NewProcess(name, command string, color int) (proc *Process) {
	proc = &Process{
		exec.Command("/bin/sh", "-c", command),
		name,
		color,
	}

	multiterm.PipeOutput(proc)

	return
}

func (p *Process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *Process) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if err := p.Cmd.Run(); err != nil {
		multiterm.FlushToStdout(p, bytes.NewBufferString(
			fmt.Sprintf("\033[1;31m%v\033[0m\n", err),
		))
	}
}

func (p *Process) Term() {
	if p.Running() {
		if err := p.Process.Signal(syscall.SIGTERM); err != nil {
			p.Stderr.Write([]byte(err.Error()))
		}
	}
}
