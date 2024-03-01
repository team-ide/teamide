package module_mongodb

import (
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/elasticsearch"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"sync"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/base"
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
	Power            = base.AppendPower(&base.PowerAction{Action: "mongodb", Text: "Mongodb", ShouldLogin: true, StandAlone: true})
	check            = base.AppendPower(&base.PowerAction{Action: "check", Text: "Mongodb测试", ShouldLogin: true, StandAlone: true, Parent: Power})
	infoPower        = base.AppendPower(&base.PowerAction{Action: "info", Text: "Mongodb信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	indexesPower     = base.AppendPower(&base.PowerAction{Action: "indexes", Text: "ES索引查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	indexStatPower   = base.AppendPower(&base.PowerAction{Action: "indexStat", Text: "ES索引状态", ShouldLogin: true, StandAlone: true, Parent: Power})
	createIndexPower = base.AppendPower(&base.PowerAction{Action: "createIndex", Text: "ES创建索引", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteIndexPower = base.AppendPower(&base.PowerAction{Action: "deleteIndex", Text: "ES删除索引", ShouldLogin: true, StandAlone: true, Parent: Power})
	getMappingPower  = base.AppendPower(&base.PowerAction{Action: "getMapping", Text: "ES索引信息查询", ShouldLogin: true, StandAlone: true, Parent: Power})
	putMappingPower  = base.AppendPower(&base.PowerAction{Action: "putMapping", Text: "ES索引修改", ShouldLogin: true, StandAlone: true, Parent: Power})
	searchPower      = base.AppendPower(&base.PowerAction{Action: "search", Text: "ES搜索", ShouldLogin: true, StandAlone: true, Parent: Power})
	scrollPower      = base.AppendPower(&base.PowerAction{Action: "scroll", Text: "ES滚动搜索", ShouldLogin: true, StandAlone: true, Parent: Power})
	requestPower     = base.AppendPower(&base.PowerAction{Action: "request", Text: "ES HTTP请求", ShouldLogin: true, StandAlone: true, Parent: Power})
	insertDataPower  = base.AppendPower(&base.PowerAction{Action: "insertData", Text: "ES插入数据", ShouldLogin: true, StandAlone: true, Parent: Power})
	updateDataPower  = base.AppendPower(&base.PowerAction{Action: "updateData", Text: "ES修改数据", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteDataPower  = base.AppendPower(&base.PowerAction{Action: "deleteData", Text: "ES删除数据", ShouldLogin: true, StandAlone: true, Parent: Power})
	reindexPower     = base.AppendPower(&base.PowerAction{Action: "reindex", Text: "ES复制索引", ShouldLogin: true, StandAlone: true, Parent: Power})
	indexAliasPower  = base.AppendPower(&base.PowerAction{Action: "indexAlias", Text: "ES索引别名", ShouldLogin: true, StandAlone: true, Parent: Power})
	importPower      = base.AppendPower(&base.PowerAction{Action: "import", Text: "ES导入", ShouldLogin: true, StandAlone: true, Parent: Power})
	exportPower      = base.AppendPower(&base.PowerAction{Action: "export", Text: "ES导出", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskListPower    = base.AppendPower(&base.PowerAction{Action: "taskList", Text: "ES任务列表", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskStatusPower  = base.AppendPower(&base.PowerAction{Action: "taskStatus", Text: "ES任务状态", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskStopPower    = base.AppendPower(&base.PowerAction{Action: "taskStop", Text: "ES任务停止", ShouldLogin: true, StandAlone: true, Parent: Power})
	taskCleanPower   = base.AppendPower(&base.PowerAction{Action: "taskClean", Text: "ES任务清理", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower       = base.AppendPower(&base.PowerAction{Action: "close", Text: "ES关闭", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: check, Do: this_.check})
	apis = append(apis, &base.ApiWorker{Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Power: indexesPower, Do: this_.indexes})
	apis = append(apis, &base.ApiWorker{Power: indexStatPower, Do: this_.indexStat})
	apis = append(apis, &base.ApiWorker{Power: createIndexPower, Do: this_.createIndex})
	apis = append(apis, &base.ApiWorker{Power: deleteIndexPower, Do: this_.deleteIndex})
	apis = append(apis, &base.ApiWorker{Power: getMappingPower, Do: this_.getMapping})
	apis = append(apis, &base.ApiWorker{Power: putMappingPower, Do: this_.putMapping})
	apis = append(apis, &base.ApiWorker{Power: searchPower, Do: this_.search})
	apis = append(apis, &base.ApiWorker{Power: scrollPower, Do: this_.scroll})
	apis = append(apis, &base.ApiWorker{Power: requestPower, Do: this_.request})
	apis = append(apis, &base.ApiWorker{Power: insertDataPower, Do: this_.insertData})
	apis = append(apis, &base.ApiWorker{Power: updateDataPower, Do: this_.updateData})
	apis = append(apis, &base.ApiWorker{Power: deleteDataPower, Do: this_.deleteData})
	apis = append(apis, &base.ApiWorker{Power: reindexPower, Do: this_.reindex})
	apis = append(apis, &base.ApiWorker{Power: indexAliasPower, Do: this_.indexAlias})
	apis = append(apis, &base.ApiWorker{Power: importPower, Do: this_._import})
	apis = append(apis, &base.ApiWorker{Power: exportPower, Do: this_.export})
	apis = append(apis, &base.ApiWorker{Power: taskStatusPower, Do: this_.taskStatus, NotRecodeLog: true})
	apis = append(apis, &base.ApiWorker{Power: taskListPower, Do: this_.taskList})
	apis = append(apis, &base.ApiWorker{Power: taskStopPower, Do: this_.taskStop})
	apis = append(apis, &base.ApiWorker{Power: taskCleanPower, Do: this_.taskClean})
	apis = append(apis, &base.ApiWorker{Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *elasticsearch.Config, err error) {
	config = &elasticsearch.Config{}
	_, err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(esConfig *elasticsearch.Config) (res elasticsearch.IService, err error) {
	key := "elasticsearch-" + esConfig.Url
	if esConfig.Username != "" {
		key += "-" + base.GetMd5String(key+esConfig.Username)
	}
	if esConfig.Password != "" {
		key += "-" + base.GetMd5String(key+esConfig.Password)
	}
	if esConfig.CertPath != "" {
		key += "-" + base.GetMd5String(key+esConfig.CertPath)
	}

	var serviceInfo *base.ServiceInfo
	serviceInfo, err = base.GetService(key, func() (res *base.ServiceInfo, err error) {
		var s elasticsearch.IService
		s, err = elasticsearch.New(esConfig)
		if err != nil {
			util.Logger.Error("getService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		_, err = s.Info()
		if err != nil {
			util.Logger.Error("getService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Close()
			}
			return
		}
		res = &base.ServiceInfo{
			WaitTime:    10 * 60 * 1000,
			LastUseTime: util.GetNowMilli(),
			Service:     s,
			Stop:        s.Close,
		}
		return
	})
	if err != nil {
		return
	}
	res = serviceInfo.Service.(elasticsearch.IService)
	serviceInfo.SetLastUseTime()
	return
}

type BaseRequest struct {
	WorkerId        string                 `json:"workerId"`
	IndexName       string                 `json:"indexName"`
	Id              string                 `json:"id"`
	Mapping         map[string]interface{} `json:"mapping"`
	PageIndex       int                    `json:"pageIndex"`
	PageSize        int                    `json:"pageSize"`
	ScrollId        string                 `json:"scrollId"`
	Doc             interface{}            `json:"doc"`
	SourceIndexName string                 `json:"sourceIndexName"`
	DestIndexName   string                 `json:"destIndexName"`
	AliasName       string                 `json:"aliasName"`
	WhereList       []*elasticsearch.Where `json:"whereList"`
	OrderList       []*elasticsearch.Order `json:"orderList"`
	TaskId          string                 `json:"taskId"`
}

func (this_ *api) check(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	_, err = getService(config)
	if err != nil {
		return
	}

	return
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	res, err = service.Info()
	if err != nil {
		return
	}
	return
}

func (this_ *api) indexes(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	res, err = service.Indexes()
	if err != nil {
		return
	}

	return
}
func (this_ *api) indexStat(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.IndexStat(request.IndexName)
	if err != nil {
		return
	}

	return
}

func (this_ *api) createIndex(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.CreateIndex(request.IndexName, request.Mapping)
	if err != nil {
		return
	}
	return
}
func (this_ *api) deleteIndex(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.DeleteIndex(request.IndexName)
	if err != nil {
		return
	}
	return
}
func (this_ *api) getMapping(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.GetMapping(request.IndexName)
	if err != nil {
		return
	}
	return
}
func (this_ *api) putMapping(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	err = service.PutMapping(request.IndexName, request.Mapping)
	if err != nil {
		return
	}
	return
}
func (this_ *api) search(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Search(request.IndexName, request.PageIndex, request.PageSize, request.WhereList, request.OrderList)
	if err != nil {
		return
	}
	return
}
func (this_ *api) scroll(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Scroll(request.IndexName, request.ScrollId, request.PageSize, request.WhereList, request.OrderList)
	if err != nil {
		return
	}
	return
}
func (this_ *api) request(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := elasticsearch.PerformRequestOptions{}
	if !base.RequestJSON(&request, c) {
		return
	}

	response, err := service.PerformRequest(request)
	if err != nil {
		return
	}
	data := make(map[string]interface{})
	res = data
	data["header"] = response.Header
	data["body"] = string(response.Body)
	return
}
func (this_ *api) insertData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Insert(request.IndexName, request.Id, request.Doc)
	if err != nil {
		return
	}
	return
}
func (this_ *api) updateData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Update(request.IndexName, request.Id, request.Doc)
	if err != nil {
		return
	}
	return
}
func (this_ *api) deleteData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Delete(request.IndexName, request.Id)
	if err != nil {
		return
	}
	return
}
func (this_ *api) reindex(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.Reindex(request.SourceIndexName, request.DestIndexName)
	if err != nil {
		return
	}
	return
}
func (this_ *api) indexAlias(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res, err = service.IndexAlias(request.IndexName, request.DestIndexName)
	if err != nil {
		return
	}
	return
}

func (this_ *api) _import(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(config)
	if err != nil {
		return
	}

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	var task = &elasticsearch.ImportTask{}
	if !base.RequestJSON(task, c) {
		return
	}

	task.Service = service
	elasticsearch.StartImportTask(task)
	addWorkerTask(request.WorkerId, task.TaskId)
	res = task

	return
}

func (this_ *api) export(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	//config, err := this_.getConfig(requestBean, c)
	//if err != nil {
	//	return
	//}
	//service, err := getService(config)
	//if err != nil {
	//	return
	//}
	//
	//var task = &elasticsearch.ImportTask{}
	//if !base.RequestJSON(task, c) {
	//	return
	//}
	//
	//task.Service = service
	//elasticsearch.StartImportTask(task)
	//res = task

	return
}

func (this_ *api) taskStatus(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res = elasticsearch.GetTask(request.TaskId)
	return
}

func (this_ *api) taskStop(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	elasticsearch.StopTask(request.TaskId)
	return
}

func (this_ *api) taskClean(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	removeWorkerTask(request.WorkerId, request.TaskId)
	return
}

func (this_ *api) taskList(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	res = getWorkerTasks(request.WorkerId)
	return
}

func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	var request = &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}
	removeWorkerTasks(request.WorkerId)
	return
}

var (
	workerTasksCache     = map[string][]string{}
	workerTasksCacheLock = &sync.Mutex{}
)

func addWorkerTask(workerId string, taskId string) {
	workerTasksCacheLock.Lock()
	defer workerTasksCacheLock.Unlock()
	taskIds := workerTasksCache[workerId]
	if util.StringIndexOf(taskIds, taskId) < 0 {
		taskIds = append(taskIds, taskId)
		workerTasksCache[workerId] = taskIds
	}
	return
}

func getWorkerTasks(workerId string) (taskList []*elasticsearch.Task) {
	workerTasksCacheLock.Lock()
	defer workerTasksCacheLock.Unlock()
	taskIds := workerTasksCache[workerId]
	for _, id := range taskIds {
		task := elasticsearch.GetTask(id)
		if task != nil {
			taskList = append(taskList, task)
		}
	}
	return
}

func removeWorkerTasks(workerId string) {
	workerTasksCacheLock.Lock()
	defer workerTasksCacheLock.Unlock()
	taskIds := workerTasksCache[workerId]
	for _, taskId := range taskIds {
		elasticsearch.StopTask(taskId)
		elasticsearch.CleanTask(taskId)
	}
	delete(workerTasksCache, workerId)
	return
}

func removeWorkerTask(workerId string, taskId string) {
	workerTasksCacheLock.Lock()
	defer workerTasksCacheLock.Unlock()

	elasticsearch.StopTask(taskId)
	elasticsearch.CleanTask(taskId)

	taskIds := workerTasksCache[workerId]
	var newIds []string
	for _, id := range taskIds {
		if id != taskId {
			newIds = append(newIds, id)
		}
	}
	taskIds = newIds
	if len(taskIds) == 0 {
		delete(workerTasksCache, workerId)
	} else {
		workerTasksCache[workerId] = taskIds
	}
	return
}
