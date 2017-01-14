package main

import (
	"os"
	"os/exec"
	"syscall"
)

type Process struct {
	*exec.Cmd

	Name  string
	Color int

	multiterm *Multiterm
}

func NewProcess(name, command string, color int, root string, multiterm *Multiterm) (proc *Process) {
	proc = &Process{
		exec.Command("/bin/sh", "-c", command),
		name,
		color,
		multiterm,
	}

	proc.Dir = root

	proc.multiterm.Connect(proc)

	return
}

func (p *Process) writeLine(b []byte) {
	p.multiterm.WriteLine(p, b)
}

func (p *Process) writeErr(err error) {
	p.multiterm.WriteErr(p, err)
}

func (p *Process) signal(sig os.Signal) {
	group, err := os.FindProcess(-p.Process.Pid)
	if err != nil {
		p.writeErr(err)
		return
	}

	if err = group.Signal(sig); err != nil {
		p.writeErr(err)
	}
}

func (p *Process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *Process) Run() {
	p.multiterm.PipeOutput(p)
	defer p.multiterm.ClosePipe(p)

	p.writeLine([]byte("\033[1mRunning...\033[0m"))

	if err := p.Cmd.Run(); err != nil {
		p.writeErr(err)
	} else {
		p.writeLine([]byte("\033[1mProcess exited\033[0m"))
	}
}

func (p *Process) Interrupt() {
	if p.Running() {
		p.writeLine([]byte("\033[1mInterrupting...\033[0m"))
		p.signal(syscall.SIGINT)
	}
}

func (p *Process) Kill() {
	if p.Running() {
		p.writeLine([]byte("\033[1mKilling...\033[0m"))
		p.signal(syscall.SIGKILL)
	}
}
