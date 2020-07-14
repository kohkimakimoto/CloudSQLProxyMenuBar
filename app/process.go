package app

import (
	"github.com/getlantern/systray"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
)

type Process struct {
	App         *App
	ProxyConfig *ProxyConfig
	Item        *systray.MenuItem
	Cmd         *exec.Cmd
}

func NewProcess(app *App, config *ProxyConfig, item *systray.MenuItem) *Process {
	return &Process{
		App:         app,
		ProxyConfig: config,
		Item:        item,
		Cmd:         nil,
	}
}

func (p *Process) Run() error {
	args, err := shellwords.Parse(p.ProxyConfig.Options)
	if err != nil {
		return err
	}

	cmd := exec.Command(p.App.Config.Core.CloudSqlProxy, args[0:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = app.LogFile
	cmd.Stderr = app.LogFile
	cmd.Dir = p.App.Dir

	if err := cmd.Start(); err != nil {
		return err
	}

	p.Cmd = cmd
	p.Item.Check()

	return cmd.Wait()
}

func (p *Process) Kill() error {
	if p.Cmd == nil {
		return nil
	}

	return p.Cmd.Process.Kill()
}

func (p *Process) Shutdown() {
	p.App.RemoveProcess(p.ProxyConfig.Name)
	p.Item.Uncheck()
}
