package main

import (
	"os"
	"path/filepath"

	"github.com/DarthSim/hivemind/_third_party/github.com/alecthomas/kingpin"
)

var config struct {
	Procfile           *os.File
	Root               string
	PortBase, PortStep int
}

func init() {
	var err error

	portBase := kingpin.Flag("port", "Specify a port to use as the base").Default("5000").Short('p').Int()
	portStep := kingpin.Flag("port-step", "Specify a step to increase port number").Default("100").Short('P').Int()
	root := kingpin.Flag("root", "Specify a working directory of application. Default: directory containing the Procfile").Short('d').String()
	procfile := kingpin.Arg("procfile", "Specify a Procfile to load").Default("./Procfile").File()

	kingpin.Parse()

	config.Procfile = *procfile
	config.PortBase = *portBase
	config.PortStep = *portStep

	if len(*root) > 0 {
		config.Root = *root
	} else {
		config.Root = filepath.Dir(config.Procfile.Name())
	}

	config.Root, err = filepath.Abs(config.Root)
	fatalOnErr(err)
}
