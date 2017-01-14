package main

import (
	"os"
	"os/exec"
	"syscall"
)

type process struct {
	*exec.Cmd

	Name  string
	Color int

	output *multiOutput
}

func newProcess(name, command string, color int, root string, output *multiOutput) (proc *process) {
	proc = &process{
		exec.Command("/bin/sh", "-c", command),
		name,
		color,
		output,
	}

	proc.Dir = root

	proc.output.Connect(proc)

	return
}

func (p *process) writeLine(b []byte) {
	p.output.WriteLine(p, b)
}

func (p *process) writeErr(err error) {
	p.output.WriteErr(p, err)
}

func (p *process) signal(sig os.Signal) {
	group, err := os.FindProcess(-p.Process.Pid)
	if err != nil {
		p.writeErr(err)
		return
	}

	if err = group.Signal(sig); err != nil {
		p.writeErr(err)
	}
}

func (p *process) Running() bool {
	return p.Process != nil && p.ProcessState == nil
}

func (p *process) Run() {
	p.output.PipeOutput(p)
	defer p.output.ClosePipe(p)

	p.writeLine([]byte("\033[1mRunning...\033[0m"))

	if err := p.Cmd.Run(); err != nil {
		p.writeErr(err)
	} else {
		p.writeLine([]byte("\033[1mProcess exited\033[0m"))
	}
}

func (p *process) Interrupt() {
	if p.Running() {
		p.writeLine([]byte("\033[1mInterrupting...\033[0m"))
		p.signal(syscall.SIGINT)
	}
}

func (p *process) Kill() {
	if p.Running() {
		p.writeLine([]byte("\033[1mKilling...\033[0m"))
		p.signal(syscall.SIGKILL)
	}
}
