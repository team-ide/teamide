package module_elasticsearch

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"teamide/internal/base"
	"teamide/internal/module/module_toolbox"
	"teamide/pkg/elasticsearch"
	"teamide/pkg/util"
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
	Power            = base.AppendPower(&base.PowerAction{Action: "elasticsearch", Text: "Elasticsearch", ShouldLogin: true, StandAlone: true})
	infoPower        = base.AppendPower(&base.PowerAction{Action: "info", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	indexesPower     = base.AppendPower(&base.PowerAction{Action: "indexes", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	createIndexPower = base.AppendPower(&base.PowerAction{Action: "createIndex", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteIndexPower = base.AppendPower(&base.PowerAction{Action: "deleteIndex", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	getMappingPower  = base.AppendPower(&base.PowerAction{Action: "getMapping", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	putMappingPower  = base.AppendPower(&base.PowerAction{Action: "putMapping", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	searchPower      = base.AppendPower(&base.PowerAction{Action: "search", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	scrollPower      = base.AppendPower(&base.PowerAction{Action: "scroll", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	insertDataPower  = base.AppendPower(&base.PowerAction{Action: "insertData", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	updateDataPower  = base.AppendPower(&base.PowerAction{Action: "updateData", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	deleteDataPower  = base.AppendPower(&base.PowerAction{Action: "deleteData", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	reindexPower     = base.AppendPower(&base.PowerAction{Action: "reindex", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
	closePower       = base.AppendPower(&base.PowerAction{Action: "close", Text: "Elasticsearch信息", ShouldLogin: true, StandAlone: true, Parent: Power})
)

func (this_ *api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/info"}, Power: infoPower, Do: this_.info})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/indexes"}, Power: indexesPower, Do: this_.indexes})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/createIndex"}, Power: createIndexPower, Do: this_.createIndex})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/deleteIndex"}, Power: deleteIndexPower, Do: this_.deleteIndex})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/getMapping"}, Power: getMappingPower, Do: this_.getMapping})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/putMapping"}, Power: putMappingPower, Do: this_.putMapping})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/search"}, Power: searchPower, Do: this_.search})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/scroll"}, Power: scrollPower, Do: this_.scroll})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/insertData"}, Power: insertDataPower, Do: this_.insertData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/updateData"}, Power: updateDataPower, Do: this_.updateData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/deleteData"}, Power: deleteDataPower, Do: this_.deleteData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/reindex"}, Power: reindexPower, Do: this_.reindex})
	apis = append(apis, &base.ApiWorker{Apis: []string{"elasticsearch/close"}, Power: closePower, Do: this_.close})

	return
}

func (this_ *api) getConfig(requestBean *base.RequestBean, c *gin.Context) (config *elasticsearch.Config, err error) {
	config = &elasticsearch.Config{}
	err = this_.toolboxService.BindConfig(requestBean, c, config)
	if err != nil {
		return
	}
	return
}

func getService(esConfig elasticsearch.Config) (res *elasticsearch.V7Service, err error) {
	key := "elasticsearch-" + esConfig.Url
	if esConfig.Username != "" {
		key += "-" + util.GetMd5String(key+esConfig.Username)
	}
	if esConfig.Password != "" {
		key += "-" + util.GetMd5String(key+esConfig.Password)
	}
	if esConfig.CertPath != "" {
		key += "-" + util.GetMd5String(key+esConfig.CertPath)
	}

	var service util.Service
	service, err = util.GetService(key, func() (res util.Service, err error) {
		var s *elasticsearch.V7Service
		s, err = elasticsearch.CreateESService(esConfig)
		if err != nil {
			util.Logger.Error("getService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		_, err = s.GetClient()
		if err != nil {
			util.Logger.Error("getService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		res = s
		return
	})
	if err != nil {
		return
	}
	res = service.(*elasticsearch.V7Service)
	return
}

type BaseRequest struct {
	IndexName       string                 `json:"indexName"`
	Id              string                 `json:"id"`
	Mapping         map[string]interface{} `json:"mapping"`
	PageIndex       int                    `json:"pageIndex"`
	PageSize        int                    `json:"pageSize"`
	ScrollId        string                 `json:"scrollId"`
	Doc             interface{}            `json:"doc"`
	SourceIndexName string                 `json:"sourceIndexName"`
	DestIndexName   string                 `json:"destIndexName"`
	WhereList       []*elasticsearch.Where `json:"whereList"`
	OrderList       []*elasticsearch.Order `json:"orderList"`
}

func (this_ *api) info(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
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
	service, err := getService(*config)
	if err != nil {
		return
	}

	res, err = service.Indexes()
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
func (this_ *api) insertData(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	config, err := this_.getConfig(requestBean, c)
	if err != nil {
		return
	}
	service, err := getService(*config)
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
	service, err := getService(*config)
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
func (this_ *api) close(requestBean *base.RequestBean, c *gin.Context) (res interface{}, err error) {
	return
}
