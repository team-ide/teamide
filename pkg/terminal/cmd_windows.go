package terminal

import (
	"os/exec"
)

func IsWindows() bool {
	return true
}

func getCmd() (cmd *exec.Cmd) {
	cmd = exec.Command("conhost", "--headless")
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return
}
