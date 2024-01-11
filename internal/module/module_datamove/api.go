package module_datamove

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/datamove"
	"github.com/team-ide/go-tool/db"
	"github.com/team-ide/go-tool/elasticsearch"
	"github.com/team-ide/go-tool/kafka"
	"github.com/team-ide/go-tool/redis"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
	"teamide/pkg/ssh"
)

type api struct {
	toolboxService *module_toolbox.ToolboxService
}

func NewApi(toolboxService *module_toolbox.ToolboxService) *api {
	return &api{
		toolboxService: toolboxService,
	}
}

var (
	Power    = base.AppendPower(&base.PowerAction{Action: "datamove", Text: "DataMove", ShouldLogin: true, StandAlone: true})
	start    = base.AppendPower(&base.PowerAction{Action: "start", Text: "启动", ShouldLogin: true, StandAlone: true, Parent: Power})
	stop     = base.AppendPower(&base.PowerAction{Action: "stop", Text: "停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	delete_  = base.AppendPower(&base.PowerAction{Action: "delete", Text: "删除", ShouldLogin: true, StandAlone: true, Parent: Power})
	get      = base.AppendPower(&base.PowerAction{Action: "get", Text: "获取", ShouldLogin: true, StandAlone: true, Parent: Power})
	list     = base.AppendPower(&base.PowerAction{Action: "list", Text: "列表", ShouldLogin: true, StandAlone: true, Parent: Power})
	download = base.AppendPower(&base.PowerAction{Action: "download", Text: "下载", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: start, Do: this_.start})
	apis = append(apis, &base.ApiWorker{Power: stop, Do: this_.stop})
	apis = append(apis, &base.ApiWorker{Power: delete_, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Power: get, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: list, Do: this_.list})
	apis = append(apis, &base.ApiWorker{Power: download, Do: this_.download})

	return
}

func (this_ *api) fullConfig_(toolboxId int64, config *datamove.DataSourceConfig) (err error) {
	var sshConfig *ssh.Config
	switch config.Type {
	case "database":
		config.DbConfig = &db.Config{}
		sshConfig, err = this_.toolboxService.BindConfigById(toolboxId, config.DbConfig)
		if sshConfig != nil {
			config.DbConfig.SSHClient, err = ssh.NewClient(*sshConfig)
			if err != nil {
				util.Logger.Error("fullConfig_ ssh NewClient error", zap.Error(err))
				return
			}
		}
		break
	case "elasticsearch":
		config.EsConfig = &elasticsearch.Config{}
		sshConfig, err = this_.toolboxService.BindConfigById(toolboxId, config.EsConfig)
		break
	case "kafka":
		config.KafkaConfig = &kafka.Config{}
		sshConfig, err = this_.toolboxService.BindConfigById(toolboxId, config.KafkaConfig)
		break
	case "redis":
		config.RedisConfig = &redis.Config{}
		sshConfig, err = this_.toolboxService.BindConfigById(toolboxId, config.RedisConfig)
		if sshConfig != nil {
			config.RedisConfig.SSHClient, err = ssh.NewClient(*sshConfig)
			if err != nil {
				util.Logger.Error("fullConfig_ ssh NewClient error", zap.Error(err))
				return
			}
		}
		break

	}
	return
}
func (this_ *api) fullConfig(request *BaseRequest) (err error) {
	err = this_.fullConfig_(request.FromToolboxId, request.From)
	if err != nil {
		return
	}
	err = this_.fullConfig_(request.ToToolboxId, request.To)
	if err != nil {
		return
	}
	return
}

type BaseRequest struct {
	*datamove.Options
	TaskKey       string `json:"taskKey"`
	FromToolboxId int64  `json:"fromToolboxId"`
	ToToolboxId   int64  `json:"toToolboxId"`
}

func (this_ *api) start(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = this_.fullConfig(request)
	if err != nil {
		return
	}

	request.Key = util.GetUUID()
	request.Dir = this_.getAnnexPath(requestBean, request.Key)

	taskInfo := &TaskInfo{}
	t, err := datamove.New(request.Options)
	if err != nil {
		return
	}
	taskInfo.Task = t
	taskInfo.Request = request

	err = this_.saveInfo(requestBean, t.Key, taskInfo)
	if err != nil {
		return
	}
	addTaskInfo(taskInfo)
	go func() {
		t.Run()
		removeTaskInfo(request.Key)

		ds, _ := os.ReadDir(request.Dir)
		if len(ds) > 0 {
			taskInfo.AnnexInfo, _ = util.LoadDirInfo(request.Dir, true)
			taskInfo.HasAnnex = true
		}

		err = this_.saveInfo(requestBean, t.Key, taskInfo)
		if err != nil {
			return
		}
	}()
	return
}

func (this_ *api) stop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	find := getTaskInfo(request.TaskKey)
	if find != nil {
		find.Stop()
	}
	return
}

func (this_ *api) delete(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	find := getTaskInfo(request.TaskKey)
	if find != nil {
		find.Stop()
	}
	removeTaskInfo(request.TaskKey)
	path := this_.getTaskPath(requestBean, request.TaskKey)
	if e, _ := util.PathExists(path); e {
		err = os.RemoveAll(path)
	}
	return
}

func (this_ *api) get(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	infoPath := this_.getInfoPath(requestBean, request.TaskKey)

	bs, _ := os.ReadFile(infoPath)
	if len(bs) == 0 {
		return
	}
	res = getTaskInfo(request.TaskKey)
	if res != nil {
		return
	}
	res = &TaskInfo{}
	_ = json.Unmarshal(bs, res)

	return
}

func (this_ *api) list(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	tasksDir := this_.getTasksDir(requestBean)
	if e, _ := util.PathExists(tasksDir); !e {
		return
	}
	var listData []*TaskInfo
	ds, err := os.ReadDir(tasksDir)
	if err != nil {
		return
	}
	for _, d := range ds {
		if d.IsDir() {
			find := getTaskInfo(d.Name())
			if find != nil {
				listData = append(listData, find)
				continue
			}
			infoPath := this_.getInfoPath(requestBean, d.Name())

			bs, _ := os.ReadFile(infoPath)
			if len(bs) == 0 {
				continue
			}
			taskInfo := &TaskInfo{}
			_ = json.Unmarshal(bs, taskInfo)
			listData = append(listData, taskInfo)
		}
	}
	res = listData
	return
}

func (this_ *api) getInfoPath(requestBean *base.RequestBean, taskKey string) (path string) {
	path = this_.getTaskPath(requestBean, taskKey) + "info.json"
	return
}

func (this_ *api) getAnnexPath(requestBean *base.RequestBean, taskKey string) (path string) {
	path = this_.getTaskPath(requestBean, taskKey) + "annex/"
	if e, _ := util.PathExists(path); !e {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return
}

func (this_ *api) saveInfo(requestBean *base.RequestBean, taskKey string, taskInfo *TaskInfo) (err error) {

	bs, err := json.Marshal(taskInfo)
	if err != nil {
		this_.toolboxService.Logger.Error("task info to json error", zap.Error(err))
		return
	}
	path := this_.getInfoPath(requestBean, taskKey)
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(bs)
	return
}

func (this_ *api) getTaskPath(requestBean *base.RequestBean, taskKey string) (path string) {
	path = this_.getTasksDir(requestBean) + taskKey + "/"
	return
}
func (this_ *api) getTasksDir(requestBean *base.RequestBean) (dir string) {
	dir = this_.toolboxService.GetFilesDir() + "uses/"
	dir = fmt.Sprintf("%s/%d/data-move/", dir, requestBean.JWT.UserId)
	if e, _ := util.PathExists(dir); !e {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	return
}

func (this_ *api) download(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	data := map[string]string{}
	err = c.Bind(&data)
	if err != nil {
		return
	}

	taskKey := data["taskKey"]
	if taskKey == "" {
		err = errors.New("taskKey获取失败")
		return
	}

	taskInfoPath := this_.getInfoPath(requestBean, taskKey)
	if e, _ := util.PathExists(taskInfoPath); !e {
		err = errors.New("任务不存在")
		return
	}
	bs, _ := os.ReadFile(taskInfoPath)
	taskInfo := &TaskInfo{}
	if len(bs) != 0 {
		_ = json.Unmarshal(bs, taskInfo)
	}
	if taskInfo.Request == nil || taskInfo.Request.Dir == "" {
		err = errors.New("任务导出文件丢失")
		return
	}

	path := taskInfo.Request.Dir
	exists, err := util.PathExists(path)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New("文件不存在")
		return
	}
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-1]
	}
	var fileName string
	var fileSize int64
	ff, err := os.Lstat(path)
	if err != nil {
		return
	}
	var fileInfo *os.File
	if ff.IsDir() {
		fs, _ := os.ReadDir(path)
		if len(fs) == 1 && !fs[0].IsDir() {
			fPath := path + "/" + fs[0].Name()
			ff, _ = os.Lstat(fPath)
			fileInfo, err = os.Open(fPath)
			if err != nil {
				return
			}
		} else {
			exists, err = util.PathExists(path + ".zip")
			if err != nil {
				return
			}
			if !exists {
				err = util.Zip(path, path+".zip")
				if err != nil {
					return
				}
			}
			ff, err = os.Lstat(path + ".zip")
			if err != nil {
				return
			}
			fileInfo, err = os.Open(path + ".zip")
			if err != nil {
				return
			}
		}
	} else {
		fileInfo, err = os.Open(path)
		if err != nil {
			return
		}
	}
	fileName = ff.Name()
	fileSize = ff.Size()

	defer func() {
		_ = fileInfo.Close()
	}()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	_, err = io.Copy(c.Writer, fileInfo)
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	res = base.HttpNotResponse
	return
}

type TaskInfo struct {
	*task.Task
	Request   *BaseRequest  `json:"request"`
	HasAnnex  bool          `json:"hasAnnex"`
	AnnexInfo *util.DirInfo `json:"annexInfo"`
}

var taskInfoCache = map[string]*TaskInfo{}
var taskInfoLocker = &sync.Mutex{}

func getTaskInfo(taskKey string) *TaskInfo {
	taskInfoLocker.Lock()
	defer taskInfoLocker.Unlock()

	return taskInfoCache[taskKey]
}

func addTaskInfo(task *TaskInfo) {
	taskInfoLocker.Lock()
	defer taskInfoLocker.Unlock()

	taskInfoCache[task.Key] = task
}

func removeTaskInfo(taskKey string) {
	taskInfoLocker.Lock()
	defer taskInfoLocker.Unlock()

	delete(taskInfoCache, taskKey)
}
