package main

import (
	"path/filepath"

	"github.com/alecthomas/kingpin"
)

var config struct {
	Procfile           string
	Root               string
	PortBase, PortStep int
	Timeout            int
}

func init() {
	var err error

	portBase := kingpin.Flag("port", "Specify a port to use as the base").Default("5000").Short('p').Int()
	portStep := kingpin.Flag("port-step", "Specify a step to increase port number").Default("100").Short('P').Int()
	root := kingpin.Flag("root", "Specify a working directory of application. Default: directory containing the Procfile").Short('d').String()
	timeout := kingpin.Flag("timeout", "Specify the amount of time (in seconds) processes have to shut down gracefully before being brutally killed").Default("5").Short('t').Int()
	procfile := kingpin.Arg("procfile", "Specify a Procfile to load").Default("./Procfile").String()

	kingpin.Parse()

	config.Procfile = *procfile
	config.PortBase = *portBase
	config.PortStep = *portStep
	config.Timeout = *timeout

	if config.Timeout < 1 {
		fatal("Timeout should be greater than 0")
	}

	if len(*root) > 0 {
		config.Root = *root
	} else {
		config.Root = filepath.Dir(config.Procfile)
	}

	config.Root, err = filepath.Abs(config.Root)
	fatalOnErr(err)
}
