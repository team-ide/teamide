package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"strings"
	"testing"
	"time"
)

func TestSftp(t *testing.T) {
	client, err := NewClient(Config{
		Address:  "127.0.0.1:22",
		Username: "root",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}

	s, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	_, err = s.SendRequest("pty-req", true, ssh.Marshal(&ptyRequestMsg{
		Term: "xterm",
	}))
	if err != nil {
		panic(err)
	}
	_, err = s.SendRequest("shell", true, nil)
	//err = NewSSHShell(&terminal.Size{Cols: 128, Rows: 128}, s)
	if err != nil {
		panic(err)
	}
	//if err = s.RequestSubsystem("sftp"); err != nil {
	//	panic(err)
	//}
	stdout, err := s.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stdin, err := s.StdinPipe()
	if err != nil {
		panic(err)
	}

	_, err = stdin.Write([]byte(`
export PS1=''
sftp localhost -ssubsystem

`))
	if err != nil {
		panic(err)
	}
	var bs = make([]byte, 1024*1024)
	var n int
	for {
		n, err = stdout.Read(bs)
		if err != nil {
			panic(err)
		}
		str := string(bs[:n])
		fmt.Println("on read bytes:", bs[:n])
		fmt.Println("on read str:", str)
		if strings.HasSuffix(str, "password: ") {

			fmt.Println("should input password")
			_, err = stdin.Write([]byte(`123456
`))
			//break
		}
		if strings.HasPrefix(str, "Connected to") {
			break
		}

		//if strings.Contains(str, "sftp localhost") {
		//	break
		//}
	}
	//_, err = stdin.Write([]byte(`ok`))
	fmt.Println("NewClientPipe start")
	c, err := sftp.NewClientPipe(stdout, stdin, sftp.MaxPacket(1024*32))
	if err != nil {
		panic(err)
	}
	fmt.Println("NewClientPipe end")
	fmt.Println(c)
	fmt.Println(c.Getwd())
	//sftp.NewClientPipe()
	time.Sleep(time.Second * 100)
}
