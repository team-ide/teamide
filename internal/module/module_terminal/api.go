package module_terminal

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"teamide/internal/base"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/terminal"
	"teamide/pkg/util"
)

type api struct {
	*worker
}

func NewApi(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *api {
	return &api{
		worker: NewWorker(toolboxService_, nodeService_),
	}
}

var (
	// Terminal 权限

	// Power 文件管理器 基本 权限
	Power           = base.AppendPower(&base.PowerAction{Action: "terminal", Text: "工具", ShouldLogin: false, StandAlone: true})
	PowerWebsocket  = base.AppendPower(&base.PowerAction{Action: "terminal_websocket", Text: "工具", ShouldLogin: true, StandAlone: true})
	PowerClose      = base.AppendPower(&base.PowerAction{Action: "terminal_close", Text: "工具", ShouldLogin: true, StandAlone: true})
	PowerKet        = base.AppendPower(&base.PowerAction{Action: "terminal_key", Text: "工具", ShouldLogin: true, StandAlone: true})
	PowerChangeSize = base.AppendPower(&base.PowerAction{Action: "terminal_change_size", Text: "工具", ShouldLogin: true, StandAlone: true})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal"}, Power: Power, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/key"}, Power: PowerKet, Do: this_.key})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/websocket"}, Power: PowerWebsocket, Do: this_.websocket, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/changeSize"}, Power: PowerChangeSize, Do: this_.changeSize})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/close"}, Power: PowerClose, Do: this_.close})

	return
}

func (this_ *api) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}

func (this_ *api) key(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	res = util.UUID()
	return
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
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
	place := c.Query("place")
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	placeId := c.Query("placeId")
	if placeId == "" {
		err = errors.New("placeId获取失败")
		return
	}
	cols, _ := strconv.Atoi(c.Query("cols"))
	rows, _ := strconv.Atoi(c.Query("rows"))
	width, _ := strconv.Atoi(c.Query("width"))
	height, _ := strconv.Atoi(c.Query("height"))
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	service := this_.GetService(key)
	if service != nil {
		err = errors.New("会话[" + key + "]已存在")
		return
	}

	err = this_.Start(key, place, placeId, &terminal.Size{
		Cols:   cols,
		Rows:   rows,
		Width:  width,
		Height: height,
	}, ws)
	if err != nil {
		_ = ws.Close()
		return
	}

	res = base.HttpNotResponse
	return
}

type Request struct {
	Key string `json:"key,omitempty"`
	*terminal.Size
}

func (this_ *api) close(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}
	this_.stopService(request.Key)
	return
}

func (this_ *api) changeSize(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}
	service := this_.GetService(request.Key)
	if service != nil {
		err = service.ChangeSize(request.Size)
	}
	return
}
