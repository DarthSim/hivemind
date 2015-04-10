package main

import (
	"fmt"
	"os"
)

func fatalOnErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(i ...interface{}) {
	fmt.Print("\033[1;31m")
	fmt.Print(i...)
	fmt.Print("\033[0m\n")

	os.Exit(1)
}
