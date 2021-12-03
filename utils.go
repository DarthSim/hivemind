package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func fatalOnErr(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(i ...interface{}) {
	fmt.Fprint(os.Stderr, "hivemind: ")
	fmt.Fprintln(os.Stderr, i...)
	os.Exit(1)
}

func splitAndTrim(str string) (res []string) {
	split := strings.Split(str, ",")
	for _, s := range split {
		s = strings.Trim(s, " ")
		if len(s) > 0 {
			res = append(res, s)
		}
	}
	return
}

func stringsContain(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func scanLines(r io.Reader, callback func([]byte) bool) error {
	var (
		err      error
		line     []byte
		isPrefix bool
	)

	reader := bufio.NewReader(r)
	buf := new(bytes.Buffer)

	for {
		line, isPrefix, err = reader.ReadLine()
		if err != nil {
			break
		}

		buf.Write(line)

		if !isPrefix {
			if !callback(buf.Bytes()) {
				return nil
			}
			buf.Reset()
		}
	}
	if err != io.EOF && err != io.ErrClosedPipe {
		return err
	}
	return nil
}
