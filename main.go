package main

import (
	"github.com/getlantern/systray"
	"github.com/kohkimakimoto/CloudSQLProxyMenuBar/app"
	"log"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	if err := app.Boot(); err != nil {
		log.Println(err)
	}
}

func onExit() {
	app.HandleExit()
}
