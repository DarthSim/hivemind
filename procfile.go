package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

type ProcfileEntry struct {
	Name    string
	Command string
}

func parseProcfile(path string) (entries []ProcfileEntry) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	re, _ := regexp.Compile("^(\\w+):\\s+(.+)$")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			params := re.FindStringSubmatch(scanner.Text())
			if len(params) < 2 {
				log.Fatal("Invalid process format: ", scanner.Text())
			}

			entries = append(entries, ProcfileEntry{
				params[1],
				params[2],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 {
		log.Fatal("No entries was found in Procfile")
	}

	return
}
