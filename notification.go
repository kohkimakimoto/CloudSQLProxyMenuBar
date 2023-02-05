package main

import "github.com/gen2brain/beeep"

type NotificationSender interface {
	HandleError(err error)
	Notify(message string) error
}

type notificationSender struct{}

func (n *notificationSender) HandleError(err error) {
	_ = beeep.Alert("CloudSQLProxyMenuBar: Error", err.Error(), "")
}

func (n *notificationSender) Notify(message string) error {
	return beeep.Notify("CloudSQLProxyMenuBar", message, "")
}
