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

func (p *Process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *Process) Run(wg *sync.WaitGroup, done chan bool) {
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer func() { done <- true }()

		multiterm.PipeOutput(p)
		defer multiterm.ClosePipe(p)

		if err := p.Cmd.Run(); err != nil {
			multiterm.WriteErr(p, err)
		} else {
			multiterm.WriteLine(p, []byte("\033[1mProcess exited\033[0m"))
		}
	}(wg)
}

func (p *Process) Term() {
	go func() {
		if p.Running() {
			group, err := os.FindProcess(-p.Process.Pid)
			if err != nil {
				multiterm.WriteErr(p, err)
				return
			}

			if err = group.Signal(syscall.SIGINT); err != nil {
				multiterm.WriteErr(p, err)
			}
		}
	}()
}
