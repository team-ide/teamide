package toolbox

import (
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"teamide/pkg/elasticsearch"
	"teamide/pkg/util"
)

func getESService(esConfig elasticsearch.Config) (res *elasticsearch.V7Service, err error) {
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

	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *elasticsearch.V7Service
		s, err = elasticsearch.CreateESService(esConfig)
		if err != nil {
			util.Logger.Error("getESService error", zap.Any("key", key), zap.Error(err))
			if s != nil {
				s.Stop()
			}
			return
		}
		_, err = s.GetClient()
		if err != nil {
			util.Logger.Error("getESService error", zap.Any("key", key), zap.Error(err))
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

type ElasticsearchBaseRequest struct {
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

func ESWork(work string, config *elasticsearch.Config, data map[string]interface{}) (res map[string]interface{}, err error) {

	var service *elasticsearch.V7Service

	if work != "close" {
		service, err = getESService(*config)
		if err != nil {
			return
		}
	}

	dataBS, err := json.Marshal(data)
	if err != nil {
		return
	}
	request := &ElasticsearchBaseRequest{}
	err = json.Unmarshal(dataBS, request)
	if err != nil {
		return
	}

	res = map[string]interface{}{}
	switch work {
	case "info":
		var info *elastic.NodesInfoResponse
		info, err = service.Info()
		if err != nil {
			return
		}
		res["info"] = info
		break
	case "indexNames":
		var indexNames []string
		indexNames, err = service.IndexNames()
		if err != nil {
			return
		}
		res["indexNames"] = indexNames
		break
	case "createIndex":
		err = service.CreateIndex(request.IndexName, request.Mapping)
		if err != nil {
			return
		}
		break
	case "deleteIndex":
		err = service.DeleteIndex(request.IndexName)
		if err != nil {
			return
		}
		break
	case "getMapping":
		var mapping interface{}
		mapping, err = service.GetMapping(request.IndexName)
		if err != nil {
			return
		}
		res["mapping"] = mapping
		break
	case "putMapping":
		err = service.PutMapping(request.IndexName, request.Mapping)
		if err != nil {
			return
		}
		break
	case "search":
		var queryResult *elasticsearch.SearchResult
		queryResult, err = service.Search(request.IndexName, request.PageIndex, request.PageSize, request.WhereList, request.OrderList)
		if err != nil {
			return
		}
		res["result"] = queryResult
		break
	case "scroll":
		var result *elasticsearch.SearchResult
		result, err = service.Scroll(request.IndexName, request.ScrollId, request.PageSize, request.WhereList, request.OrderList)
		if err != nil {
			return
		}
		res["result"] = result
		break
	case "insertData":
		var result *elasticsearch.IndexResponse
		result, err = service.Insert(request.IndexName, request.Id, request.Doc)
		if err != nil {
			return
		}
		res["result"] = result
		break
	case "updateData":
		var result *elasticsearch.UpdateResponse
		result, err = service.Update(request.IndexName, request.Id, request.Doc)
		if err != nil {
			return
		}
		res["result"] = result
		break
	case "deleteData":
		var result *elasticsearch.DeleteResponse
		result, err = service.Delete(request.IndexName, request.Id)
		if err != nil {
			return
		}
		res["result"] = result
		break
	case "reindex":
		var result *elasticsearch.BulkIndexByScrollResponse
		result, err = service.Reindex(request.SourceIndexName, request.DestIndexName)
		if err != nil {
			return
		}
		res["result"] = result
		break
	case "close":
		break
	}
	return
}
