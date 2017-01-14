package main

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

const baseColor = 32

type hivemindConfig struct {
	Procfile           string
	Root               string
	PortBase, PortStep int
	Timeout            int
}

type hivemind struct {
	conf        hivemindConfig
	output      multiOutput
	procs       []*process
	procWg      sync.WaitGroup
	done        chan bool
	interrupted chan os.Signal
}

func newHivemind(conf hivemindConfig) (h *hivemind) {
	h = &hivemind{conf: conf}
	h.output = multiOutput{}

	entries := parseProcfile(h.conf.Procfile, h.conf.PortBase, h.conf.PortStep)
	h.procs = make([]*process, len(entries))

	for i, entry := range entries {
		h.procs[i] = newProcess(entry.Name, entry.Command, baseColor+i, h.conf.Root, &h.output)
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
	case <-time.After(time.Duration(h.conf.Timeout) * time.Second):
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
