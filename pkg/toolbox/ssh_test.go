package toolbox

import (
	"fmt"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"teamide/pkg/util"
	"testing"
)

func TestSSHClient(t *testing.T) {
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

	err = sshSession.RequestSubsystem("sftp linkdood@192.168.6.81:22")
	if err != nil {
		util.Logger.Error("RequestSubsystem error", zap.Error(err))
		panic(err)
	}
	//
	pw, err := sshSession.StdinPipe()
	if err != nil {
		panic(err)
	}
	pr, err := sshSession.StdoutPipe()
	if err != nil {
		panic(err)
	}
	ftpClient, err := sftp.NewClientPipe(pr, pw)
	if err != nil {
		util.Logger.Error("NewClientPipe error", zap.Error(err))
		panic(err)
	}

	path, err := ftpClient.Getwd()
	if err != nil {
		util.Logger.Error("Getwd error", zap.Error(err))
		panic(err)
	}
	fmt.Println(path)
}

var (
	sshConfig81 = SSHConfig{
		Type:     "tcp",
		Address:  "192.168.6.81:22",
		Username: "root",
		Password: "bxyvrv1601",
	}
	sshConfigJumps = SSHConfig{
		Type:      "tcp",
		Address:   "jumps.linkdood.cn:10022",
		Username:  "zhuliang",
		Password:  "oTkQFYYienzwB3Fb",
		PublicKey: "D:\\Workspaces\\Code\\note\\工作\\zhuliang.pem",
	}
)

func newSSHClient() (client *ssh.Client, err error) {
	client, err = NewSSHClient(
		sshConfig81,
	)
	if err != nil {
		return
	}
	return
}
