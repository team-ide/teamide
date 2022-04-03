package toolbox

import (
	"errors"
	"github.com/gorilla/websocket"
)

func WSSFPTConnection(token string, ws *websocket.Conn) (err error) {
	var sshConfig *SSHConfig = sshTokenCache[token]
	if sshConfig == nil {
		err = errors.New("令牌会话丢失")
		return
	}
	SSHClient := &SSHClient{
		Token:  token,
		Config: sshConfig,
	}
	err = SSHClient.StartSftp(ws)

	return
}
