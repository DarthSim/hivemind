package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
)

func createProcesses() (procs []*Process) {
	color := 32
	port := config.PortBase

	for _, entry := range parseProcfile("Procfile") {
		procs = append(
			procs,
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

func interrupted() chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}

func detectExit(done chan bool) {
	for {
		var exit bool

		select {
		case <-done:
			exit = config.ExitTogether
		case <-interrupted():
			exit = true
		}

		if exit {
			break
		}
	}
}

func main() {
	var procWg sync.WaitGroup

	procs := createProcesses()

	done := make(chan bool, len(procs))

	for _, proc := range procs {
		proc.Run(&procWg, done)
	}

	go func() {
		detectExit(done)

		for _, proc := range procs {
			proc.Term()
		}
	}()

	procWg.Wait()
}
