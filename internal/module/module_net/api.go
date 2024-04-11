package module_net

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	goSSH "golang.org/x/crypto/ssh"
	"net/http"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/ssh"
)

type api struct {
	*module_toolbox.ToolboxService
}

func NewApi(toolboxService_ *module_toolbox.ToolboxService) *api {
	return &api{
		ToolboxService: toolboxService_,
	}
}

var (
	// Terminal 权限

	// Power 文件管理器 基本 权限
	Power          = base.AppendPower(&base.PowerAction{Action: "connection", Text: "网络链接", ShouldLogin: true, StandAlone: true})
	websocketPower = base.AppendPower(&base.PowerAction{Action: "websocket", Text: "网络链接WebSocket", ShouldLogin: true, StandAlone: true, Parent: Power})
	check          = base.AppendPower(&base.PowerAction{Action: "check", Text: "网络链接测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower     = base.AppendPower(&base.PowerAction{Action: "close", Text: "网络链接关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
	keyPower       = base.AppendPower(&base.PowerAction{Action: "key", Text: "网络链接Key", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: keyPower, Do: this_.key})
	apis = append(apis, &base.ApiWorker{Power: websocketPower, Do: this_.websocket, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Power: check, Do: this_.check})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *Config, err error) {
	config = &Config{}
	sshConfig, err := this_.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	if sshConfig != nil {
		var sshClient *goSSH.Client
		sshClient, err = ssh.NewClient(*sshConfig)
		if err != nil {
			util.Logger.Error("create net conn service ssh NewClient error", zap.Error(err))
			return
		}
		config.SSHClient = sshClient
	}
	return
}

func (this_ *api) key(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := createService(config)
	if err != nil {
		return
	}
	service.Key = util.GetUUID()
	setService(service.Key, service)
	data := make(map[string]interface{})
	data["key"] = service.Key
	res = data
	return
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  32 * 1024,
	WriteBufferSize: 32 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (this_ *api) websocket(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	if request.JWT == nil || request.JWT.UserId == 0 {
		err = errors.New("登录用户获取失败")
		return
	}
	key := c.Query("key")
	if key == "" {
		err = errors.New("key获取失败")
		return
	}
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	service := getService(key)
	if service == nil {
		err = errors.New("会话[" + key + "]不存在")

		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("service not found:"+err.Error()))
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	err = service.start(ws)
	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("start error:"+err.Error()))
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	res = base.HttpNotResponse
	return
}

type Request struct {
	Key string `json:"key"`
}

func (this_ *api) check(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := createService(config)
	if err != nil {
		return
	}
	err = service.init()
	if err != nil {
		return
	}
	service.stop()
	return
}

func (this_ *api) close(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service := getService(request.Key)
	if service != nil {
		service.stop()
	}
	return
}
