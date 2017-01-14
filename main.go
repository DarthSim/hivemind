package main

import (
	"os"
	"path/filepath"

	cli "gopkg.in/urfave/cli.v1"
)

const version = "1.0.1"

func main() {
	var (
		conf hivemindConfig
		err  error
	)

	app := cli.NewApp()

	app.Name = "Hivemind"
	app.HelpName = "hivemind"
	app.Usage = "The mind to rule processes of your development environment"
	app.Description = "Hivemind is a process manager for Procfile-based applications"
	app.Author = "Sergey \"DarthSim\" Alexandrovich"
	app.Email = "darthsim@gmail.com"
	app.Version = version
	app.ArgsUsage = "[procfile]"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "processes, l", Usage: "Specify process names to launch. Divide names with comma", Destination: &conf.ProcNames},
		cli.IntFlag{Name: "port, p", Usage: "specify a port to use as the base", Value: 5000, Destination: &conf.PortBase},
		cli.IntFlag{Name: "port-step, P", Usage: "specify a step to increase port number", Value: 100, Destination: &conf.PortStep},
		cli.StringFlag{Name: "root, d", Usage: "specify a working directory of application. Default: directory containing the Procfile", Destination: &conf.Root},
		cli.IntFlag{Name: "timeout, t", Usage: "specify the amount of time (in seconds) processes have to shut down gracefully before being brutally killed", Value: 5, Destination: &conf.Timeout},
	}

	app.Action = func(c *cli.Context) error {
		switch c.NArg() {
		case 0:
			conf.Procfile = "./Procfile"
		case 1:
			conf.Procfile = c.Args().First()
		default:
			fatal("Specify a single procfile")
		}

		if conf.Timeout < 1 {
			fatal("Timeout should be greater than 0")
		}

		if len(conf.Root) == 0 {
			conf.Root = filepath.Dir(conf.Procfile)
		}

		conf.Root, err = filepath.Abs(conf.Root)
		fatalOnErr(err)

		newHivemind(conf).Run()

		return nil
	}

	app.Run(os.Args)
}
