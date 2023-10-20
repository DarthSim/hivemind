package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"

	_ "github.com/DarthSim/godotenv/autoload"
)

const version = "1.2.0"

func main() {
	var (
		conf hivemindConfig
		err  error
	)

	app := &cli.App{
		Name:     "Hivemind",
		HelpName: "hivemind",
		Usage:    "The mind to rule processes of your development environment",
		Authors: []*cli.Author{
			{
				Name:  "Sergey \"DarthSim\" Alexandrovich",
				Email: "darthsim@gmail.com",
			},
		},
		Version:   version,
		ArgsUsage: "[procfile] (Use '-' to read from stdin, Procfile path can be also set with $HIVEMIND_PROCFILE)",
		HideHelp:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "title, w",
				EnvVars:     []string{"HIVEMIND_TITLE"},
				Usage:       "Specify a title of the application",
				Destination: &conf.Title,
			},
			&cli.StringFlag{
				Name:        "processes, l",
				EnvVars:     []string{"HIVEMIND_PROCESSES"},
				Usage:       "Specify process names to launch. Divide names with comma",
				Destination: &conf.ProcNames,
			},
			&cli.IntFlag{
				Name:        "port, p",
				EnvVars:     []string{"HIVEMIND_PORT,PORT"},
				Usage:       "specify a port to use as the base",
				Value:       5000,
				Destination: &conf.PortBase,
			},
			&cli.IntFlag{
				Name:        "port-step, P",
				EnvVars:     []string{"HIVEMIND_PORT_STEP"},
				Usage:       "specify a step to increase port number",
				Value:       100,
				Destination: &conf.PortStep,
			},
			&cli.StringFlag{
				Name:        "root, d",
				EnvVars:     []string{"HIVEMIND_ROOT"},
				Usage:       "specify a working directory of application. Default: directory containing the Procfile",
				Destination: &conf.Root,
			},
			&cli.IntFlag{
				Name:        "timeout, t",
				EnvVars:     []string{"HIVEMIND_TIMEOUT"},
				Usage:       "specify the amount of time (in seconds) processes have to shut down gracefully before being brutally killed",
				Value:       5,
				Destination: &conf.Timeout,
			},
			&cli.BoolFlag{
				Name:        "no-prefix",
				EnvVars:     []string{"HIVEMIND_NO_PREFIX"},
				Usage:       "process names will not be printed if the flag is specified",
				Destination: &conf.NoPrefix,
			},
			&cli.BoolFlag{
				Name:        "print-timestamps, T",
				EnvVars:     []string{"HIVEMIND_PRINT_TIMESTAMPS"},
				Usage:       "timestamps will be printed if the flag is specified",
				Destination: &conf.PrintTimestamps,
			},
		},
		Action: func(c *cli.Context) error {
			switch c.NArg() {
			case 0:
				if path := os.Getenv("HIVEMIND_PROCFILE"); len(path) > 0 {
					conf.Procfile = path
				} else {
					conf.Procfile = "./Procfile"
				}
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
		},
	}

	app.Run(os.Args)
}
