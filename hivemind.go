package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
)

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
	color := 32
	port := config.PortBase

	for _, entry := range parseProcfile("Procfile") {
		h.procs = append(
			h.procs,
			NewProcess(
				entry.Name,
				strings.Replace(entry.Command, "$PORT", strconv.Itoa(port), -1),
				color,
			),
		)

		color++
		port += config.PortStep
	}

	return
}

func (h *Hivemind) runProcess(proc *Process) {
	h.procWg.Add(1)

	go func() {
		defer h.procWg.Done()
		defer func() { h.done <- true }()

		proc.Run(&h.procWg, h.done)
	}()
}

func (h *Hivemind) waitForExit() {
	for {
		var exit bool

		select {
		case <-h.done:
			exit = true
		case <-h.interrupted:
			exit = true
		}

		if exit {
			for _, proc := range h.procs {
				go proc.Interrupt()
			}

			break
		}
	}

	<-h.interrupted

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
