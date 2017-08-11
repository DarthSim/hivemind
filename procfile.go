package main

import (
	"bufio"
	"os"
	"regexp"
)

type procfileEntry struct {
	Name    string
	Command string
	Port    int
}

func parseProcfile(path string, portBase, portStep int) (entries []procfileEntry) {
	re, _ := regexp.Compile(`^([\w-]+):\s+(.+)$`)

	f, err := os.Open(path)
	fatalOnErr(err)

	port := portBase
	names := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		params := re.FindStringSubmatch(scanner.Text())
		if len(params) != 3 {
			continue
		}

		name, cmd := params[1], params[2]

		if names[name] {
			fatal("Process names must be uniq")
		}
		names[name] = true

		entries = append(entries, procfileEntry{name, cmd, port})

		port += portStep
	}

	fatalOnErr(scanner.Err())

	if len(entries) == 0 {
		fatal("No entries was found in Procfile")
	}

	return
}
