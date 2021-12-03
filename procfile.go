package main

import (
	"io"
	"os"
	"regexp"
)

type procfileEntry struct {
	Name    string
	Command string
	Port    int
}

func parseProcfile(path string, portBase, portStep int) (entries []procfileEntry) {
	var f io.Reader
	switch path {
	case "-":
		f = os.Stdin
	default:
		file, err := os.Open(path)
		fatalOnErr(err)
		defer file.Close()

		f = file
	}

	re, _ := regexp.Compile(`^([\w-]+):\s+(.+)$`)
	port := portBase
	names := make(map[string]bool)

	err := scanLines(f, func(b []byte) bool {
		if len(b) == 0 {
			return true
		}

		params := re.FindStringSubmatch(string(b))
		if len(params) != 3 {
			return true
		}

		name, cmd := params[1], params[2]

		if names[name] {
			fatal("Process names must be uniq")
		}
		names[name] = true

		entries = append(entries, procfileEntry{name, cmd, port})

		port += portStep

		return true
	})

	fatalOnErr(err)

	if len(entries) == 0 {
		fatal("No entries was found in Procfile")
	}

	return
}
