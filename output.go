package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"syscall"
	"encoding/json"

	"github.com/kr/pty"
)

type stringStyle int

// ColorErr is the color used for writing errors with styled text
const ColorErr int = 31

const (
	styleNone  stringStyle = iota
	styleBold
	styleError
)

// Helper struct for outputting line-delimited JSON
type jsonlog struct {
	Process	string	`json:"name"`
	Line	string	`json:"line"`
}

// LogFormat specifies the output format of logging.
type LogFormat string

const (
	logFormatText LogFormat = "text"
	logFormatJSON           = "json"
)

type ptyPipe struct {
	pty, tty *os.File
}

type multiOutput struct {
	ColorizeOutput bool
	LogFormat LogFormat

	maxNameLength int
	mutex         sync.Mutex
	pipes         map[*process]*ptyPipe
}

func (m *multiOutput) openPipe(proc *process) (pipe *ptyPipe) {
	var err error

	pipe = m.pipes[proc]

	pipe.pty, pipe.tty, err = pty.Open()
	fatalOnErr(err)

	proc.Stdout = pipe.tty
	proc.Stderr = pipe.tty
	proc.Stdin = pipe.tty
	proc.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}

	return
}

func (m *multiOutput) Connect(proc *process) {
	if len(proc.Name) > m.maxNameLength {
		m.maxNameLength = len(proc.Name)
	}

	if m.pipes == nil {
		m.pipes = make(map[*process]*ptyPipe)
	}

	m.pipes[proc] = &ptyPipe{}
}

func (m *multiOutput) PipeOutput(proc *process) {
	pipe := m.openPipe(proc)

	go func(proc *process, pipe *ptyPipe) {
		scanner := bufio.NewScanner(pipe.pty)

		for scanner.Scan() {
			m.WriteLine(proc, scanner.Bytes(), styleNone)
		}
	}(proc, pipe)
}

func (m *multiOutput) ClosePipe(proc *process) {
	if pipe := m.pipes[proc]; pipe != nil {
		pipe.pty.Close()
		pipe.tty.Close()
	}
}

func (m *multiOutput) textLogLine(proc *process, p []byte, style stringStyle) bytes.Buffer {
	var buf bytes.Buffer

	var colorCode int
	if style == styleError {
		colorCode = ColorErr
	} else {
		colorCode = proc.Color
	}

	color := fmt.Sprintf("\033[1;%vm", colorCode)

	if m.ColorizeOutput {
		buf.WriteString(color)
	}

	buf.WriteString(proc.Name)

	if m.ColorizeOutput {
		for buf.Len()-len(color) < m.maxNameLength {
			buf.WriteByte(' ')
		}
	} else {
		for buf.Len() < m.maxNameLength {
			buf.WriteByte(' ')
		}
	}

	if m.ColorizeOutput {
		buf.WriteString("\033[0m | ")
	} else {
		buf.WriteString(" | ")
	}

	if style == styleBold && m.ColorizeOutput {
		buf.WriteString("\033[1m")
	}

	buf.Write(p)

	if style == styleBold {
		buf.WriteString("\033[0m")
	}

	buf.WriteByte('\n')
	return buf
}

func (m *multiOutput) jsonLogLine(proc *process, p []byte) bytes.Buffer {
	var buf bytes.Buffer

	data := jsonlog{
		Process: proc.Name,
		Line: string(p),
	}

	line, _ := json.Marshal(&data)
	buf.Write(line)
	buf.WriteByte('\n')
	return buf
}

func (m *multiOutput) WriteLine(proc *process, p []byte, style stringStyle) {

	var buf bytes.Buffer

	switch m.LogFormat {
	case logFormatText:
		buf = m.textLogLine(proc, p, style)
	case logFormatJSON:
		buf = m.jsonLogLine(proc, p)
	default:
		panic("BUG: Log Format was invalid!")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	buf.WriteTo(os.Stdout)
}

func (m *multiOutput) WriteErr(proc *process, err error) {
	if m.ColorizeOutput {
		m.WriteLine(proc, []byte(
				fmt.Sprintf("%v", err)), styleError)
	} else {
		m.WriteLine(proc, []byte(fmt.Sprintf("%v", err) ), styleError)
	}

}
