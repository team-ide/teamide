package toolbox

import (
	"encoding/json"
	"teamide/pkg/elasticsearch"
)

func getESService(esConfig elasticsearch.Config) (res *elasticsearch.V7Service, err error) {
	key := "elasticsearch-" + esConfig.Url
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *elasticsearch.V7Service
		s, err = elasticsearch.CreateESService(esConfig)
		if err != nil {
			return
		}
		_, err = s.GetClient()
		if err != nil {
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
}

func ESWork(work string, config *elasticsearch.Config, data map[string]interface{}) (res map[string]interface{}, err error) {

	var service *elasticsearch.V7Service
	service, err = getESService(*config)
	if err != nil {
		return
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
	case "indexNames":
		var indexNames []string
		indexNames, err = service.IndexNames()
		if err != nil {
			return
		}
		res["indexNames"] = indexNames
	case "createIndex":
		err = service.CreateIndex(request.IndexName, request.Mapping)
		if err != nil {
			return
		}
	case "deleteIndex":
		err = service.DeleteIndex(request.IndexName)
		if err != nil {
			return
		}
	case "getMapping":
		var mapping interface{}
		mapping, err = service.GetMapping(request.IndexName)
		if err != nil {
			return
		}
		res["mapping"] = mapping
	case "putMapping":
		err = service.PutMapping(request.IndexName, request.Mapping)
		if err != nil {
			return
		}
	case "search":
		var queryResult *elasticsearch.SearchResult
		queryResult, err = service.Search(request.IndexName, request.PageIndex, request.PageSize)
		if err != nil {
			return
		}
		res["result"] = queryResult
	case "scroll":
		var result *elasticsearch.SearchResult
		result, err = service.Scroll(request.IndexName, request.ScrollId, request.PageSize)
		if err != nil {
			return
		}
		res["result"] = result
	case "insertData":
		var result *elasticsearch.IndexResponse
		result, err = service.Insert(request.IndexName, request.Id, request.Doc)
		if err != nil {
			return
		}
		res["result"] = result
	case "updateData":
		var result *elasticsearch.UpdateResponse
		result, err = service.Update(request.IndexName, request.Id, request.Doc)
		if err != nil {
			return
		}
		res["result"] = result
	case "deleteData":
		var result *elasticsearch.DeleteResponse
		result, err = service.Delete(request.IndexName, request.Id)
		if err != nil {
			return
		}
		res["result"] = result
	case "reindex":
		var result *elasticsearch.BulkIndexByScrollResponse
		result, err = service.Reindex(request.SourceIndexName, request.DestIndexName)
		if err != nil {
			return
		}
		res["result"] = result
	}
	return
}
