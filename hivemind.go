package main

import (
	"os"
	"os/signal"
	"sync"
	"time"
	"golang.org/x/crypto/ssh/terminal"
)

const baseColor = 32

type hivemindConfig struct {
	Procfile           string
	ProcNames          string
	Root               string
	PortBase, PortStep int
	Timeout            int
	ForceTTY		   bool
	LogFormat		   string
}

type hivemind struct {
	output      *multiOutput
	procs       []*process
	procWg      sync.WaitGroup
	done        chan bool
	interrupted chan os.Signal
	timeout     time.Duration
}

func newHivemind(conf hivemindConfig) (h *hivemind) {
	h = &hivemind{timeout: time.Duration(conf.Timeout) * time.Second}

	// Enable or disable colors based on TTY detection or user request
	colorizeOutput := terminal.IsTerminal(int(os.Stdout.Fd())) || conf.ForceTTY

	h.output = &multiOutput{
		ColorizeOutput: colorizeOutput,
		LogFormat: LogFormat(conf.LogFormat),
	}

	entries := parseProcfile(conf.Procfile, conf.PortBase, conf.PortStep)
	h.procs = make([]*process, 0)

	procNames := splitAndTrim(conf.ProcNames)

	for i, entry := range entries {
		if len(procNames) == 0 || stringsContain(procNames, entry.Name) {
			h.procs = append(h.procs, newProcess(entry.Name, entry.Command, baseColor+i, conf.Root, h.output))
		}
	}

	return
}

func (h *hivemind) runProcess(proc *process) {
	h.procWg.Add(1)

	go func() {
		defer h.procWg.Done()
		defer func() { h.done <- true }()

		proc.Run()
	}()
}

func (h *hivemind) waitForDoneOrInterrupt() {
	select {
	case <-h.done:
	case <-h.interrupted:
	}
}

func (h *hivemind) waitForTimeoutOrInterrupt() {
	select {
	case <-time.After(h.timeout):
	case <-h.interrupted:
	}
}

func (h *hivemind) waitForExit() {
	h.waitForDoneOrInterrupt()

	for _, proc := range h.procs {
		go proc.Interrupt()
	}

	h.waitForTimeoutOrInterrupt()

	for _, proc := range h.procs {
		go proc.Kill()
	}
}

func (h *hivemind) Run() {
	h.done = make(chan bool, len(h.procs))

	h.interrupted = make(chan os.Signal)
	signal.Notify(h.interrupted, os.Interrupt, os.Kill)

	for _, proc := range h.procs {
		h.runProcess(proc)
	}

	go h.waitForExit()

	h.procWg.Wait()
}
