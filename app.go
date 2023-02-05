package main

import (
	"github.com/getlantern/systray"
	"log"
	"os"
	"sync"
)

type App struct {
	Config             *Config
	Dir                string
	NotificationSender NotificationSender
	LogFile            *os.File
	Logger             *log.Logger
	Processes          map[string]*Process
	Mutex              *sync.Mutex
	ChangeProcessesCh  chan int
}

func (a *App) HandleError(err error) {
	a.Logger.Println(err)
	a.NotificationSender.HandleError(err)
}

func (a *App) Notify(msg string) error {
	return a.NotificationSender.Notify(msg)
}

func (a *App) HandleProxyAction(config *ProxyConfig, item *systray.MenuItem) {
	proc := a.GetProcess(config.Name)
	if proc != nil {
		// This proxy is running. Try to kill it.
		a.KillProcess(proc)
	} else {
		// This proxy is not running. Try to start it.
		a.StartProcess(config, item)
	}
}

func (a *App) KillProcess(proc *Process) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	delete(a.Processes, proc.ProxyConfig.Name)
	a.ChangeProcessesCh <- len(a.Processes)

	proc.Item.Uncheck()

	if err := proc.Kill(); err != nil {
		a.HandleError(err)
	}
}

func (a *App) StartProcess(config *ProxyConfig, item *systray.MenuItem) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	proc := &Process{
		CloudSqlProxy: a.Config.Core.CloudSqlProxy,
		Dir:           a.Dir,
		LogFile:       a.LogFile,
		ProxyConfig:   config,
		Item:          item,
		Cmd:           nil,
	}

	a.Processes[proc.ProxyConfig.Name] = proc
	a.ChangeProcessesCh <- len(a.Processes)

	go func(proc *Process) {
		defer a.KillProcess(proc)
		proc.Item.Check()
		if err := proc.Run(); err != nil {
			a.HandleError(err)
		}
	}(proc)
}

func (a *App) GetProcess(name string) *Process {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	return a.Processes[name]
}

func (a *App) HandleExit() {
	for _, proc := range a.Processes {
		a.KillProcess(proc)
	}

	if a.LogFile != nil {
		_ = a.LogFile.Close()
	}
	_ = a.Notify("The CloudSQLProxyMenuBar was stopped.")
}
