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

func (p *process) writeLine(b []byte, style StringStyle) {
	p.output.WriteLine(p, b, style)
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

	p.writeLine([]byte("Running..."), StyleBold)

	if err := p.Cmd.Run(); err != nil {
		p.writeErr(err)
	} else {
		p.writeLine([]byte("Process exited"), StyleBold)
	}
}

func (p *process) Interrupt() {
	if p.Running() {
		p.writeLine([]byte("Interrupting..."), StyleBold)
		p.signal(syscall.SIGINT)
	}
}

func (p *process) Kill() {
	if p.Running() {
		p.writeLine([]byte("Killing..."), StyleBold)
		p.signal(syscall.SIGKILL)
	}
}
