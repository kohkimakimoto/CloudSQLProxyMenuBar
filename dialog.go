package main

import (
	"fmt"
	"os/exec"
)

type DisplayDialogFunc func(msg string) error

func DisplayDialog(msg string) error {
	osa, err := exec.LookPath("osascript")
	if err != nil {
		return err
	}
	script := fmt.Sprintf("display dialog %q with title \"CloudSQLProxyMenuBar\" with icon caution buttons {\"OK\"} default button \"OK\"", msg)
	cmd := exec.Command(osa, "-e", script)
	return cmd.Run()
}
