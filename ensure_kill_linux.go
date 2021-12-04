//go:build linux
// +build linux

package main

import "syscall"

func ensureKill(p *process) {
	p.SysProcAttr.Pdeathsig = syscall.SIGKILL
}
