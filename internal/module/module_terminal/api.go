package module_terminal

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/ssh"
	"teamide/pkg/terminal"
)

type api struct {
	*WorkerFactory
}

func NewApi(toolboxService_ *module_toolbox.ToolboxService, nodeService_ *module_node.NodeService) *api {
	return &api{
		WorkerFactory: NewWorkerFactory(toolboxService_, nodeService_),
	}
}

var (
	// Terminal 权限

	// Power 文件管理器 基本 权限
	Power                = base.AppendPower(&base.PowerAction{Action: "terminal", Text: "终端", ShouldLogin: true, StandAlone: true})
	websocketPower       = base.AppendPower(&base.PowerAction{Action: "websocket", Text: "终端WebSocket", ShouldLogin: true, StandAlone: true, Parent: Power})
	check                = base.AppendPower(&base.PowerAction{Action: "check", Text: "终端测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower           = base.AppendPower(&base.PowerAction{Action: "close", Text: "终端关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
	keyPower             = base.AppendPower(&base.PowerAction{Action: "key", Text: "终端Key", ShouldLogin: true, StandAlone: true, Parent: Power})
	changeSizePower      = base.AppendPower(&base.PowerAction{Action: "changeSize", Text: "终端窗口大小变更", ShouldLogin: true, StandAlone: true, Parent: Power})
	uploadWebsocketPower = base.AppendPower(&base.PowerAction{Action: "uploadWebsocket", Text: "终端上传WebSocket", ShouldLogin: true, StandAlone: true, Parent: Power})
	getLogs              = base.AppendPower(&base.PowerAction{Action: "getLogs", Text: "getLogs", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteLog            = base.AppendPower(&base.PowerAction{Action: "deleteLog", Text: "deleteLog", ShouldLogin: true, StandAlone: true, Parent: Power})
	cleanLog             = base.AppendPower(&base.PowerAction{Action: "cleanLog", Text: "cleanLog", ShouldLogin: true, StandAlone: true, Parent: Power})
	downloadLog          = base.AppendPower(&base.PowerAction{Action: "downloadLog", Text: "downloadLog", ShouldLogin: true, StandAlone: true, Parent: Power})
	systemInfo           = base.AppendPower(&base.PowerAction{Action: "system/info", Text: "system", ShouldLogin: true, StandAlone: true, Parent: Power})
	systemMonitor        = base.AppendPower(&base.PowerAction{Action: "system/monitor", Text: "system", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: keyPower, Do: this_.key})
	apis = append(apis, &base.ApiWorker{Power: websocketPower, Do: this_.websocket, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Power: changeSizePower, Do: this_.changeSize})
	apis = append(apis, &base.ApiWorker{Power: check, Do: this_.check})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Power: uploadWebsocketPower, Do: this_.uploadWebsocket, IsWebSocket: true})
	apis = append(apis, &base.ApiWorker{Power: getLogs, Do: this_.getLogs})
	apis = append(apis, &base.ApiWorker{Power: deleteLog, Do: this_.deleteLog})
	apis = append(apis, &base.ApiWorker{Power: cleanLog, Do: this_.cleanLog})
	apis = append(apis, &base.ApiWorker{Power: downloadLog, Do: this_.downloadLog})
	apis = append(apis, &base.ApiWorker{Power: systemInfo, Do: this_.systemInfo})
	apis = append(apis, &base.ApiWorker{Power: systemMonitor, Do: this_.systemMonitor, NotRecodeLog: true})

	return
}

func (this_ *api) key(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service, _, err := this_.createService(&CreateParam{
		place:    request.Place,
		placeId:  request.PlaceId,
		workerId: request.WorkerId,
		lastUser: request.LastUser,
		lastDir:  request.LastDir,
	})
	if err != nil {
		return
	}

	data := make(map[string]interface{})

	data["isWindows"], err = service.service.IsWindows()
	if err != nil {
		return
	}
	data["key"] = util.GetUUID()
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
	place := c.Query("place")
	if place == "" {
		err = errors.New("place获取失败")
		return
	}
	placeId := c.Query("placeId")
	workerId := c.Query("workerId")
	if workerId == "" {
		err = errors.New("workerId获取失败")
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

		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("service create error:"+err.Error()))
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	userAgentStr := c.Request.UserAgent()
	baseLog := &TerminalLogModel{
		Ip:        c.ClientIP(),
		UserAgent: userAgentStr,
		Place:     place,
		PlaceId:   placeId,
		WorkerId:  workerId,
	}

	baseLog.UserId = request.JWT.UserId
	baseLog.UserName = request.JWT.Name
	baseLog.UserAccount = request.JWT.Account
	baseLog.LoginId = request.JWT.LoginId
	err = this_.Start(key,
		&CreateParam{
			place:    place,
			placeId:  placeId,
			workerId: workerId,
			lastUser: c.Query("lastUser"),
			lastDir:  c.Query("lastDir"),
		},
		&terminal.Size{
			Cols: cols,
			Rows: rows,
		}, ws, baseLog)
	if err != nil {
		_ = ws.WriteMessage(websocket.BinaryMessage, []byte("start error:"+err.Error()))
		this_.Logger.Error("websocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	res = base.HttpNotResponse
	return
}

func (this_ *api) systemInfo(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service := this_.GetService(request.Key)
	if service == nil || service.service == nil {
		return
	}

	data := map[string]interface{}{}
	data["cpuInfo"], _ = service.service.GetCpuInfo()

	res = data
	return
}

func (this_ *api) systemMonitor(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	service := this_.GetService(request.Key)
	if service == nil || service.service == nil {
		return
	}

	data := map[string]interface{}{}
	data["cpuPercent"], _ = service.service.GetCpuPercent()
	data["memInfo"], _ = service.service.GetMemInfo()
	data["diskStats"], _ = service.service.GetDiskStats()

	res = data
	return
}
func (this_ *api) uploadWebsocket(request *base.RequestBean, c *gin.Context) (res interface{}, err error) {

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

	service := this_.GetService(key)
	if service == nil {
		err = errors.New("会话[" + key + "]不存在")

		this_.Logger.Error("uploadWebsocket start error", zap.Error(err))
		_ = ws.Close()
		return
	}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				this_.Logger.Error("uploadWebsocket error", zap.Any("error", e))
			}
		}()
		var isClosed bool
		ws.SetCloseHandler(func(code int, text string) error {
			isClosed = true
			return nil
		})
		var buf []byte
		var readErr error
		var writeErr error
		var n int
		for !isClosed {
			_, buf, readErr = ws.ReadMessage()
			if readErr != nil && readErr != io.EOF {
				break
			}
			//this_.Logger.Info("ws on read", zap.Any("bs", string(buf)))
			n, writeErr = service.service.Write(buf)
			if writeErr != nil {
				break
			}
			if n >= 200 {
				writeErr = ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", n)))
				if writeErr != nil {
					break
				}
			}
			if readErr == io.EOF {
				readErr = nil
				break
			}
		}

		if readErr != nil && !isClosed {
			this_.Logger.Error("uploadWebsocket read error", zap.Error(readErr))
		}

		if writeErr != nil && !isClosed {
			this_.Logger.Error("uploadWebsocket write error", zap.Error(writeErr))
		}
	}()

	res = base.HttpNotResponse
	return
}

type Request struct {
	Place    string `json:"place,omitempty"`
	PlaceId  string `json:"placeId,omitempty"`
	Key      string `json:"key,omitempty"`
	WorkerId string `json:"workerId"`
	LastUser string `json:"lastUser,omitempty"`
	LastDir  string `json:"lastDir,omitempty"`
	*terminal.Size
}

func (this_ *api) check(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &module_toolbox.ToolboxModel{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Option == "" {
		err = errors.New("SSH 配置不存在")
		return
	}

	var config *ssh.Config
	config, err = this_.toolboxService.GetSSHConfig(request.Option)
	if err != nil {
		return
	}

	service := ssh.NewTerminalService(config, "", "")

	err = service.TestClient()
	return
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
		err = service.service.ChangeSize(request.Size)
	}
	return
}

func (this_ *api) getLogs(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = this_.WorkerFactory.getLogs(request.Place, request.PlaceId)
	return
}

func (this_ *api) deleteLog(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	path := this_.WorkerFactory.getLogPath(request.Place, request.PlaceId, request.WorkerId)
	if ex, _ := util.PathExists(path); ex {
		err = os.Remove(path)
	}
	return
}

func (this_ *api) cleanLog(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &Request{}
	if !base.RequestJSON(request, c) {
		return
	}

	path := this_.WorkerFactory.getLogPath(request.Place, request.PlaceId, request.WorkerId)
	if ex, _ := util.PathExists(path); !ex {
		return
	}
	f, err := os.Create(path)
	defer func() { _ = f.Close() }()
	_, err = f.WriteString("")
	return
}

func (this_ *api) downloadLog(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	this_.Logger.Info("下载日志记录 start")
	res = base.HttpNotResponse
	defer func() {
		if err != nil {
			_, _ = c.Writer.WriteString(err.Error())
		}
	}()

	request := map[string]string{}

	err = c.Bind(&request)
	if err != nil {
		return
	}

	fileName := "" + request["fileName"] + ".log"
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(fileName)))

	// 此处不设置 文件大小，如果设置文件大小，将无法终止下载
	//c.Header("Content-Length", fmt.Sprint(fileInfo.Size))
	c.Header("download-file-name", fileName)

	path := this_.WorkerFactory.getLogPath(request["place"], request["placeId"], request["workerId"])

	if ex, _ := util.PathExists(path); ex {
		var f *os.File
		f, err = os.Open(path)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		_, err = io.Copy(c.Writer, f)
	} else {
		_, err = c.Writer.WriteString("暂无日志记录")
	}
	c.Status(http.StatusOK)
	return
}
