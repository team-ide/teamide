package terminal

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os/exec"
)

func IsWindows() bool {
	return true
}

func start(size *Size) (starter *terminalStart, err error) {
	cmd := exec.Command("conhost", "--headless")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		util.Logger.Error("cmd StdoutPipe error", zap.Error(err))
		return
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		util.Logger.Error("cmd StdinPipe error", zap.Error(err))
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		util.Logger.Error("cmd StderrPipe error", zap.Error(err))
		return
	}

	err = cmd.Start()
	if err != nil {
		util.Logger.Error("cmd Start error", zap.Error(err))
		return
	}
	starter = &terminalStart{
		Stop: func() {

			_ = stdout.Close()
			_ = stdin.Close()
			_ = stderr.Close()
			if cmd != nil && cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		},
		Write_: stdin.Write,
		Read_:  stdout.Read,
	}
	return
}
