package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

const baseColor = 32

type Hivemind struct {
	procs       []*Process
	procWg      sync.WaitGroup
	done        chan bool
	interrupted chan os.Signal
}

func NewHivemind() (h *Hivemind) {
	h = &Hivemind{}
	h.createProcesses()
	return
}

func (h *Hivemind) createProcesses() {
	entries := parseProcfile()
	h.procs = make([]*Process, len(entries))

	for i, entry := range entries {
		port := config.PortBase + config.PortStep*i
		h.procs[i] = NewProcess(
			entry.Name,
			strings.Replace(entry.Command, "$PORT", strconv.Itoa(port), -1),
			baseColor+i,
		)
	}
}

func (h *Hivemind) runProcess(proc *Process) {
	h.procWg.Add(1)

	go func() {
		defer h.procWg.Done()
		defer func() { h.done <- true }()

		proc.Run()
	}()
}

func (h *Hivemind) waitForDoneOrInterrupt() {
	select {
	case <-h.done:
	case <-h.interrupted:
	}
}

func (h *Hivemind) waitForTimeoutOrInterrupt() {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(config.Timeout) * time.Second)
		timeout <- true
	}()

	select {
	case <-timeout:
	case <-h.interrupted:
	}
}

func (h *Hivemind) waitForExit() {
	h.waitForDoneOrInterrupt()

	for _, proc := range h.procs {
		go proc.Interrupt()
	}

	h.waitForTimeoutOrInterrupt()

	for _, proc := range h.procs {
		go proc.Kill()
	}
}

func (h *Hivemind) Run() {
	h.done = make(chan bool, len(h.procs))

	h.interrupted = make(chan os.Signal)
	signal.Notify(h.interrupted, os.Interrupt, os.Kill)

	for _, proc := range h.procs {
		h.runProcess(proc)
	}

	go h.waitForExit()

	h.procWg.Wait()
}
