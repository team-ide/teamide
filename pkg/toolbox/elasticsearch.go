package toolbox

import (
	"encoding/json"
	"teamide/pkg/elasticsearch"
	"teamide/pkg/form"
)

func elasticsearchWorker() *Worker {
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
		OtherForm: map[string]*form.Form{
			"index": {
				Fields: []*form.Field{
					{
						Label: "IndexName（索引）", Name: "indexName", DefaultValue: "index_xxx",
						Rules: []*form.Rule{
							{Required: true, Message: "索引不能为空"},
						},
					},
					{
						Label: "结构", Name: "mapping", Type: "json", DefaultValue: map[string]interface{}{
							"settings": map[string]interface{}{
								"number_of_shards":   1,
								"number_of_replicas": 0,
							},
							"mappings": map[string]interface{}{
								"properties": map[string]interface{}{
									"title": map[string]interface{}{
										"type": "text",
									},
								},
							},
						},
						Rules: []*form.Rule{
							{Required: true, Message: "结构不能为空"},
						},
					},
				},
			},
		},
	}

	return worker_
}

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

func esWork(work string, config map[string]interface{}, data map[string]interface{}) (res map[string]interface{}, err error) {

	var esConfig elasticsearch.Config
	var configBS []byte
	configBS, err = json.Marshal(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(configBS, &esConfig)
	if err != nil {
		return
	}

	var service *elasticsearch.V7Service
	service, err = getESService(esConfig)
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
