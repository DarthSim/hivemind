package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ProcfileEntry struct {
	Name    string
	Command string
}

func parseProcfile(path string, portBase, portStep int) (entries []ProcfileEntry) {
	re, _ := regexp.Compile("^(\\w+):\\s+(.+)$")

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

		entries = append(entries, ProcfileEntry{
			name,
			strings.Replace(cmd, "$PORT", strconv.Itoa(port), -1),
		})

		port += portStep
	}

	fatalOnErr(scanner.Err())

	if len(entries) == 0 {
		fatal("No entries was found in Procfile")
	}

	return
}
