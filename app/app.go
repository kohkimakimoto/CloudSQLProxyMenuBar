package app

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/kohkimakimoto/CloudSQLProxyMenuBar/assets"
	"github.com/pkg/browser"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

type App struct {
	Dir               string
	Config            *Config
	LogFile           *os.File
	Processes         map[string]*Process
	Mutex             *sync.Mutex
	ChangeProcessesCh chan int
}

var app = &App{
	Dir:               filepath.Join(userHomeDir(), ".cloudsqlproxymenubar"),
	Config:            nil,
	LogFile:           nil,
	Processes:         map[string]*Process{},
	Mutex:             new(sync.Mutex),
	ChangeProcessesCh: make(chan int),
}

func Boot() error {
	// init config and app
	cfg := NewConfig()
	app.Config = cfg
	app.Config.Core.LogFile = filepath.Join(app.Dir, "output.log")

	if _, err := os.Stat(app.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(app.Dir, os.FileMode(0700)); err != nil {
			return err
		}
	}

	cfgFile := filepath.Join(app.Dir, "config.toml")
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if err := ioutil.WriteFile(cfgFile, []byte(InitialConfig), os.FileMode(0600)); err != nil {
			return err
		}
	}

	// Load config from a file.
	if _, err := os.Stat(cfgFile); err == nil {
		if err := cfg.Load(cfgFile); err != nil {
			return err
		}
	}

	// config logger for error.
	if cfg.Core.LogFile != "" {
		f, err := os.Create(cfg.Core.LogFile)
		if err != nil {
			return err
		}
		log.SetOutput(f)
		app.LogFile = f
	}

	// init menu title.
	systray.SetIcon(assets.MustAsset("CloudSQL.png"))
	systray.SetTitle("SQL")

	// construct proxy menu items.
	for _, key := range app.Config.SortedProxyKeys() {
		proxyConfig := app.Config.Proxies[key]
		proxyItem := systray.AddMenuItem(proxyConfig.NameForItem(), proxyConfig.TooltipForItem())
		go func(proxyConfig *ProxyConfig, proxyItem *systray.MenuItem) {
			for {
				select {
				case <-proxyItem.ClickedCh:
					if err := app.handleProxyAction(proxyConfig, proxyItem); err != nil {
						log.Println(err)
					}
				}
			}
		}(proxyConfig, proxyItem)
	}

	// construct general events handlers.
	systray.AddSeparator()
	githubItem := systray.AddMenuItem("GitHub Repository", "")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "")
	go func() {
		for {
			select {
			case <-githubItem.ClickedCh:
				browser.OpenURL("https://github.com/kohkimakimoto/CloudSQLProxyMenuBar")
			case <-quitItem.ClickedCh:
				systray.Quit()
			case num := <-app.ChangeProcessesCh:
				if num == 0 {
					systray.SetTitle("SQL")
				} else {
					systray.SetTitle(fmt.Sprintf("SQL %d", num))
				}
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the process.
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.Printf("Received signal: %v", sig)
		systray.Quit()
	}()

	return nil
}

func (a *App) handleProxyAction(config *ProxyConfig, item *systray.MenuItem) error {
	proc := a.GetProcess(config.Name)
	if proc != nil {
		// this proxy has already been running. you are trying to stop it.
		defer proc.Shutdown()
		if err := proc.Kill(); err != nil {
			return err
		}
	} else {
		// this proxy is not active. you are trying to start it.
		// create new process
		proc = NewProcess(a, config, item)
		a.SetProcess(proc)

		go func(proc *Process) {
			defer proc.Shutdown()
			if err := proc.Run(); err != nil {
				log.Println(err)
			}
		}(proc)
	}

	return nil
}

func (a *App) SetProcess(proc *Process) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	a.Processes[proc.ProxyConfig.Name] = proc

	a.ChangeProcessesCh <- len(a.Processes)
}

func (a *App) GetProcess(name string) *Process {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	return a.Processes[name]
}

func (a *App) RemoveProcess(name string) {
	a.Mutex.Lock()
	defer a.Mutex.Unlock()

	delete(a.Processes, name)

	a.ChangeProcessesCh <- len(a.Processes)
}

func HandleExit() {
	log.Println("Shutting down...")

	for _, proc := range app.Processes {
		if err := proc.Kill(); err != nil {
			log.Println(err)
		}
	}

	log.Println("Done shut down process.")

	if app.LogFile != nil {
		app.LogFile.Close()
	}
}
