package toolbox

import (
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"teamide/pkg/form"
)

func init() {
	worker_ := &Worker{
		Name: "elasticsearch",
		Text: "Elasticsearch",
		Work: esWork,
		ConfigForm: &form.Form{
			Fields: []*form.Field{
				{
					Label: "连接地址（http://127.0.0.1:9200）", Name: "url", DefaultValue: "http://127.0.0.1:9200",
					Rules: []*form.Rule{
						{Required: true, Message: "连接地址不能为空"},
					},
				},
			},
		},
	}

	AddWorker(worker_)
}

type ElasticsearchBaseRequest struct {
	IndexName string                 `json:"indexName"`
	Id        string                 `json:"id"`
	Mapping   map[string]interface{} `json:"mapping"`
	PageIndex int                    `json:"pageIndex"`
	PageSize  int                    `json:"pageSize"`
	ScrollId  string                 `json:"scrollId"`
	Doc       interface{}            `json:"doc"`
}

type ESConfig struct {
	Url string `json:"url"`
}

func esWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var esConfig ESConfig
	var bs []byte
	bs, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &esConfig)
	if err != nil {
		return
	}

	var service *ESService
	service, err = getESService(esConfig)
	if err != nil {
		return
	}

	bs, err = json.Marshal(data)
	if err != nil {
		return
	}
	request := &ElasticsearchBaseRequest{}
	err = json.Unmarshal(bs, request)
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
		var queryResult *elastic.SearchResult
		queryResult, err = service.Search(request.IndexName, request.PageIndex, request.PageSize)
		if err != nil {
			return
		}
		res["result"] = queryResult
	case "scroll":
		var queryResult *ESQueryResult
		queryResult, err = service.Scroll(request.IndexName, request.ScrollId, request.PageSize)
		if err != nil {
			return
		}
		res["result"] = queryResult
	case "insertData":
		var result *elastic.IndexResponse
		result, err = service.Insert(request.IndexName, request.Id, request.Doc)
		if err != nil {
			return
		}
		res["result"] = result
	case "updateData":
		var result *elastic.UpdateResponse
		result, err = service.Update(request.IndexName, request.Id, request.Doc)
		if err != nil {
			return
		}
		res["result"] = result
	case "deleteData":
		var result *elastic.DeleteResponse
		result, err = service.Delete(request.IndexName, request.Id)
		if err != nil {
			return
		}
		res["result"] = result
	}
	return
}

type ESQueryResult struct {
	Count int64 `json:"count"`
}

type ESUpdateResult struct {
	Count int64 `json:"count"`
}
