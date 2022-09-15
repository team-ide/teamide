package terminal

import (
	"os/exec"
)

func IsWindows() bool {
	return false
}

func getCmd() (cmd *exec.Cmd) {
	cmd = exec.Command("bash")
	return
}
