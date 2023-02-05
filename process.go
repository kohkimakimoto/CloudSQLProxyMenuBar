package main

import (
	"github.com/getlantern/systray"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
)

type Process struct {
	CloudSqlProxy string
	Dir           string
	LogFile       *os.File
	ProxyConfig   *ProxyConfig
	Item          *systray.MenuItem
	Cmd           *exec.Cmd
}

func (p *Process) Run() error {
	args, err := shellwords.Parse(p.ProxyConfig.Arguments)
	if err != nil {
		return err
	}

	cmd := exec.Command(p.CloudSqlProxy, args[0:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = p.LogFile
	cmd.Stderr = p.LogFile
	cmd.Dir = p.Dir

	if err := cmd.Start(); err != nil {
		return err
	}
	p.Cmd = cmd
	return cmd.Wait()
}

func (p *Process) Kill() error {
	if p.Cmd == nil {
		return nil
	}
	return p.Cmd.Process.Kill()
}
