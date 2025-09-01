package myfunc

import (
	"os/exec"
	"runtime"
)

func openURL(url string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", url)
	} else {
		cmd = exec.Command("xdg-open", url) // Dành cho Linux
	}
	return cmd.Start()
}