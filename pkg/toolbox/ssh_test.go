package toolbox

import (
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"testing"
)

func TestSSHClient(t *testing.T) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Development = false
	//logger, _ := loggerConfig.Build()
	fmt.Println("newSSHClient start")
	sshClient, err := newSSHClient()
	fmt.Println("newSSHClient end")
	if err != nil {
		panic(err)
	}

	fmt.Println("NewSession start")
	sshSession, err := sshClient.NewSession()
	fmt.Println("NewSession end")
	if err != nil {
		panic(err)
	}
	//var terminalSize = TerminalSize{
	//	Cols: 100,
	//	Rows: 40,
	//}
	//err = NewSSHShell(terminalSize, sshSession)
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("Output start")
	bs, err := sshSession.Output("ll")
	fmt.Println("Output end")
	fmt.Println("pwd:", string(bs))
	if err != nil {
		panic(err)
	}

	//err = sshSession.RequestSubsystem("sftp")
	//if err != nil {
	//	logger.Error("RequestSubsystem error", zap.Error(err))
	//	panic(err)
	//}
	//
	//pw, err := sshSession.StdinPipe()
	//if err != nil {
	//	panic(err)
	//}
	//pr, err := sshSession.StdoutPipe()
	//if err != nil {
	//	panic(err)
	//}
	//
	//ftpClient, err := sftp.NewClientPipe(pr, pw)
	//if err != nil {
	//	panic(err)
	//}
	//
	//path, err := ftpClient.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(path)
}

var (
	sshConfig81 = SSHConfig{
		Type:     "tcp",
		Address:  "192.168.6.81:22",
		User:     "root",
		Password: "bxyvrv1601",
	}
	sshConfigJumps = SSHConfig{
		Type:      "tcp",
		Address:   "jumps.linkdood.cn:10022",
		User:      "zhuliang",
		Password:  "oTkQFYYienzwB3Fb",
		PublicKey: "D:\\Workspaces\\Code\\note\\工作\\zhuliang.pem",
	}
)

func newSSHClient() (client *ssh.Client, err error) {
	client, err = NewSSHClient(
		sshConfigJumps,
	)
	if err != nil {
		return
	}
	return
}
