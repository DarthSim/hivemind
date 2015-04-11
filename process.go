package main

import (
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

	proc.Dir = config.Root

	return
}

func (p *Process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *Process) Run(wg *sync.WaitGroup) {
	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()

		multiterm.PipeOutput(p)
		defer multiterm.ClosePipe(p)

		if err := p.Cmd.Run(); err != nil {
			multiterm.WriteErr(p, err)
		}
	}(wg)
}

func (p *Process) Term() {
	go func() {
		if p.Running() {
			if err := p.Process.Signal(syscall.SIGTERM); err != nil {
				multiterm.WriteErr(p, err)
			}
		}
	}()
}
