package module_terminal

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"net/url"
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
	PowerUpload     = base.AppendPower(&base.PowerAction{Action: "terminal_upload", Text: "工具", ShouldLogin: true, StandAlone: true})
	PowerDownload   = base.AppendPower(&base.PowerAction{Action: "terminal_download", Text: "工具", ShouldLogin: true, StandAlone: true})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal"}, Power: Power, Do: this_.index})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/key"}, Power: PowerKet, Do: this_.key})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/websocket"}, Power: PowerWebsocket, Do: this_.websocket, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/changeSize"}, Power: PowerChangeSize, Do: this_.changeSize})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/close"}, Power: PowerClose, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/upload"}, Power: PowerUpload, Do: this_.upload})
	apis = append(apis, &base.ApiWorker{Apis: []string{"terminal/download"}, Power: PowerDownload, Do: this_.download, IsGet: true})

	return
}

func (this_ *api) index(_ *base.RequestBean, _ *gin.Context) (res interface{}, err error) {
	return
}

func (this_ *api) key(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service, err := this_.createService(request.Place, request.PlaceId)
	if err != nil {
		return
	}

	data := make(map[string]interface{})

	data["isWindows"], err = service.IsWindows()
	if err != nil {
		return
	}
	data["key"] = util.UUID()
	res = data
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
		Cols: cols,
		Rows: rows,
	}, ws)
	if err != nil {
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	res = base.HttpNotResponse
	return
}

type Request struct {
	Place   string `json:"place,omitempty"`
	PlaceId string `json:"placeId,omitempty"`
	Key     string `json:"key,omitempty"`
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

func (this_ *api) upload(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	key := c.PostForm("key")
	if key == "" {
		err = errors.New("key获取失败")
		return
	}
	mF, err := c.MultipartForm()
	if err != nil {
		return
	}
	fileList := mF.File["file"]

	if len(fileList) == 0 {
		err = errors.New("upload file is not defined")
		return
	}
	err = this_.SetFileHeadersChan(key, fileList)

	return
}

func (this_ *api) download(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	data := map[string]string{}

	err = c.Bind(&data)
	if err != nil {
		return
	}

	key := data["key"]
	if key == "" {
		err = errors.New("key获取失败")
		return
	}
	name := data["name"]
	if name == "" {
		err = errors.New("name获取失败")
		return
	}
	size := data["size"]
	if size == "" {
		err = errors.New("size获取失败")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(name)))
	c.Header("Content-Length", size)
	c.Header("download-file-name", name)

	err = this_.SetWriterChan(key, c.Writer)
	if err != nil {
		return
	}
	c.Status(http.StatusOK)
	return
}
