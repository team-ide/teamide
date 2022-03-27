package toolbox

import "encoding/json"

func init() {
	worker_ := &Worker{
		Name: "elasticsearch",
		Text: "Elasticsearch",
		Work: esWork,
	}

	AddWorker(worker_)
}

type ElasticsearchBaseRequest struct {
	IndexName string                 `json:"indexName"`
	Mapping   map[string]interface{} `json:"mapping"`
	PageIndex int                    `json:"pageIndex"`
	PageSize  int                    `json:"pageSize"`
	ScrollId  string                 `json:"scrollId"`
}

func esWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {
	var service *ESService
	service, err = getESService(config["address"].(string))
	if err != nil {
		return
	}

	var bs []byte
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
		err = service.CreateIndex(request.IndexName)
		if err != nil {
			return
		}
	case "deleteIndex":
		err = service.DeleteIndex(request.IndexName)
		if err != nil {
			return
		}
	case "getMapping":
		var mapping map[string]interface{}
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
		var queryResult *ESQueryResult
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
	}
	return
}

type ESQueryResult struct {
	Count int64 `json:"count"`
}
