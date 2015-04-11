package main

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
)

func waitForKill() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func main() {
	var procs []*Process

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

	var procWg sync.WaitGroup

	for _, proc := range procs {
		proc.Run(&procWg)
	}

	waitForKill()

	for _, proc := range procs {
		proc.Term()
	}

	procWg.Wait()
}
