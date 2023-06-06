package guac

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"net"
)

// NewGuacamoleTunnel
// scheme: this.scheme,
// hostname: this.hostname,
// port: this.port,
// 'ignore-cert': this.ignoreCert,
// security: this.security,
// username: this.user,
// password: this.pass
func NewGuacamoleTunnel(guacadAddr, protocol, host, port, user, password, uuid string, w, h, dpi int) (s *SimpleTunnel, err error) {
	config := NewGuacamoleConfiguration()
	config.ConnectionID = uuid
	config.Protocol = protocol
	config.OptimalScreenHeight = h
	config.OptimalScreenWidth = w
	config.OptimalResolution = dpi
	config.AudioMimetypes = []string{"audio/L16", "rate=44100", "channels=2"}
	config.Parameters = map[string]string{
		"scheme":      protocol,
		"hostname":    host,
		"port":        port,
		"ignore-cert": "true",
		"security":    "",
		"username":    user,
		"password":    password,
	}
	addr, err := net.ResolveTCPAddr("tcp", guacadAddr)
	if err != nil {
		util.Logger.Error("NewGuacamoleTunnel ResolveTCPAddr error", zap.Error(err))
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		util.Logger.Error("NewGuacamoleTunnel DialTCP error", zap.Error(err))
		return nil, err
	}
	stream := NewStream(conn, SocketTimeout)
	// 这一步才是初始化 rdp/vnc guacd 并认证资产的身份
	err = stream.Handshake(config)
	if err != nil {
		return nil, err
	}
	tunnel := NewSimpleTunnel(stream)
	return tunnel, nil
}
