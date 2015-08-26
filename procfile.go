package main

import (
	"bufio"
	"regexp"
)

type ProcfileEntry struct {
	Name    string
	Command string
}

func parseProcfile() (entries []ProcfileEntry) {
	re, _ := regexp.Compile("^(\\w+):\\s+(.+)$")

	scanner := bufio.NewScanner(config.Procfile)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			params := re.FindStringSubmatch(scanner.Text())
			if len(params) < 2 {
				fatal("Invalid process format: ", scanner.Text())
			}

			entries = append(entries, ProcfileEntry{
				params[1],
				params[2],
			})
		}
	}

	fatalOnErr(scanner.Err())

	if len(entries) == 0 {
		fatal("No entries was found in Procfile")
	}

	return
}
