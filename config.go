package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v1"
)

var config struct {
	Procfile           *os.File
	PortBase, PortStep int
}

func init() {
	portBase := kingpin.Flag("port", "Specify a port to use as the base").Default("5000").Short('p').Int()
	portStep := kingpin.Flag("port-step", "Specify a step to increase port number").Default("100").Short('P').Int()
	procfilePath := kingpin.Arg("procfile", "Specify path to Procfile").Default("./Procfile").File()

	kingpin.Parse()

	config.Procfile = *procfilePath
	config.PortBase = *portBase
	config.PortStep = *portStep
}
