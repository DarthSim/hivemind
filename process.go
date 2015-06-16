package main

import (
	"os"
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

func (p *Process) signal(sig os.Signal) {
	if p.Running() {
		group, err := os.FindProcess(-p.Process.Pid)
		if err != nil {
			multiterm.WriteErr(p, err)
			return
		}

		if err = group.Signal(sig); err != nil {
			multiterm.WriteErr(p, err)
		}
	}
}

func (p *Process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *Process) Run(wg *sync.WaitGroup, done chan bool) {
	multiterm.PipeOutput(p)
	defer multiterm.ClosePipe(p)

	if err := p.Cmd.Run(); err != nil {
		multiterm.WriteErr(p, err)
	} else {
		multiterm.WriteLine(p, []byte("\033[1mProcess exited\033[0m"))
	}
}

func (p *Process) Interrupt() {
	p.signal(syscall.SIGINT)
}

func (p *Process) Kill() {
	p.signal(syscall.SIGKILL)
}
