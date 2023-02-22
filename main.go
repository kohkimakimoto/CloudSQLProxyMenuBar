package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
)

func main() {
	app := &App{
		Dir:               filepath.Join(homeDir(), ".cloudsqlproxymenubar"),
		DisplayDialog:     DisplayDialog,
		Logger:            log.Default(),
		Processes:         map[string]*Process{},
		Mutex:             new(sync.Mutex),
		ChangeProcessesCh: make(chan int),
	}

	// Construct config with default values
	config := NewConfig()
	config.Core.LogFile = filepath.Join(app.Dir, "output.log")
	app.Config = config

	// Register the app to the global variable
	g = app

	// Run systray
	systray.Run(onReady, onExit)
}

var g *App

func onReady() {
	if err := start(); err != nil {
		g.HandleError(err)
	}
}

func start() error {
	// Init app directory that is usually `~/.cloudsqlproxymenubar`
	if _, err := os.Stat(g.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(g.Dir, os.FileMode(0700)); err != nil {
			return err
		}
	}

	// Init config file
	cfgFile := filepath.Join(g.Dir, "config.toml")
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		if err := os.WriteFile(cfgFile, []byte(InitialConfig), os.FileMode(0600)); err != nil {
			return err
		}
	}
	// Load config from a file.
	if _, err := os.Stat(cfgFile); err == nil {
		if err := g.Config.LoadFromFile(cfgFile); err != nil {
			return err
		}
	}

	if g.Config.Core.LogFile != "" {
		f, err := os.Create(g.Config.Core.LogFile)
		if err != nil {
			return err
		}
		g.Logger.SetOutput(f)
		g.LogFile = f
	}
	
	systray.SetTemplateIcon(Icon, Icon)
	systray.SetTitle("SQL")

	// Construct proxy menu items
	for _, key := range g.Config.SortedProxyKeys() {
		proxyConfig := g.Config.Proxies[key]
		proxyItem := systray.AddMenuItem(proxyConfig.LabelOrName(), fmt.Sprintf("%s %s", g.Config.Core.CloudSqlProxy, proxyConfig.Arguments))
		go func(proxyConfig *ProxyConfig, proxyItem *systray.MenuItem) {
			for {
				select {
				case <-proxyItem.ClickedCh:
					g.HandleProxyAction(proxyConfig, proxyItem)
				}
			}
		}(proxyConfig, proxyItem)
	}

	systray.AddSeparator()
	githubItem := systray.AddMenuItem("GitHub Repository", "")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "")

	go func() {
		for {
			select {
			case <-githubItem.ClickedCh:
				_ = browser.OpenURL("https://github.com/kohkimakimoto/CloudSQLProxyMenuBar")
			case <-quitItem.ClickedCh:
				systray.Quit()
			case num := <-g.ChangeProcessesCh:
				if num == 0 {
					systray.SetTitle("SQL")
				} else {
					systray.SetTitle(fmt.Sprintf("SQL %d", num))
				}
			}
		}
	}()

	// Wait for interrupt signal to stop the process.
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		systray.Quit()
	}()

	return nil
}

func onExit() {
	if g.LogFile != nil {
		_ = g.LogFile.Close()
	}
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
