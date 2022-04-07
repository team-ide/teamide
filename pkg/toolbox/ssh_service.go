package toolbox

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func WSSSHConnection(token string, ws *websocket.Conn, Logger *zap.Logger) (err error) {
	var sshConfig *SSHConfig = sshTokenCache[token]
	client := SSHClient{
		Token:  token,
		Config: sshConfig,
		ws:     ws,
		Logger: Logger,
	}
	shellClient := &SSHShellClient{
		SSHClient: client,
	}
	shellClient.start()

	return
}

func WSSFPTConnection(token string, ws *websocket.Conn, Logger *zap.Logger) (err error) {
	var sshConfig *SSHConfig = sshTokenCache[token]
	client := SSHClient{
		Token:  token,
		Config: sshConfig,
		ws:     ws,
		Logger: Logger,
	}
	sftpClient := &SSHSftpClient{
		SSHClient: client,
	}
	sftpClient.start()

	return
}
