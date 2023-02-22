package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type App struct {
	Config            *Config
	Dir               string
	DisplayDialog     DisplayDialogFunc
	LogFile           *os.File
	Logger            *log.Logger
	Processes         map[string]*Process
	Mutex             *sync.Mutex
	ChangeProcessesCh chan int
}

func (a *App) HandleError(err error) {
	a.Logger.Println(err)

	if err := a.DisplayDialog(fmt.Sprintf("%s\n\nFor more information, please see the log file: %s", err.Error(), a.Config.Core.LogFile)); err != nil {
		a.Logger.Println(err)
	}
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
		proc.Item.Check()

		defer func() {
			if proc.ProxyConfig.After != "" {
				if err := proc.RunScript(proc.ProxyConfig.After); err != nil {
					a.HandleError(err)
				}
			}
			a.DeleteProcess(proc.ProxyConfig.Name)
			proc.Item.Uncheck()
		}()

		if proc.ProxyConfig.Before != "" {
			if err := proc.RunScript(proc.ProxyConfig.Before); err != nil {
				a.HandleError(err)
				return
			}
		}

		if err := proc.Run(); isSigKillErr(err) == false {
			// If the process was killed by SIGKILL, we don't need to handle the error.
			// Because it is a normal shutdown process by clicking the menu item.
			a.HandleError(err)
		}
	}(proc)
}

func isSigKillErr(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			if status.Signaled() && status.Signal() == syscall.SIGKILL {
				return true
			}
		}
	}
	return false
}

func (a *App) DeleteProcess(name string) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	delete(a.Processes, name)
	a.ChangeProcessesCh <- len(a.Processes)
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
}
