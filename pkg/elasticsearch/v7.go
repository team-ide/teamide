package elasticsearch

import (
	"context"
	"github.com/olivere/elastic/v7"
	"sort"
	"teamide/pkg/util"
)

type Config struct {
	Url string `json:"url"`
}

func CreateESService(config Config) (*V7Service, error) {
	service := &V7Service{
		Url: config.Url,
	}
	err := service.Init()
	return service, err
}

//V7Service 注册处理器在线信息等
type V7Service struct {
	Url         string
	lastUseTime int64
}

func (this_ *V7Service) Init() error {
	var err error
	return err
}

func (this_ *V7Service) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *V7Service) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *V7Service) GetClient() (client *elastic.Client, err error) {
	defer func() {
		this_.lastUseTime = util.GetNowTime()
	}()
	client, err = elastic.NewClient(
		elastic.SetURL(this_.Url),
		//docker
		elastic.SetSniff(false),
	)
	return
}

func (this_ *V7Service) Stop() {
}

func (this_ *V7Service) DeleteIndex(indexName string) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.DeleteIndex(indexName).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *V7Service) CreateIndex(indexName string, bodyJSON map[string]interface{}) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.CreateIndex(indexName).BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *V7Service) IndexNames() (res []string, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	res, err = client.IndexNames()
	if err != nil {
		return
	}

	sort.Strings(res)
	return
}

func (this_ *V7Service) GetMapping(indexName string) (res interface{}, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	mappingMap, err := client.GetMapping().Index(indexName).Do(context.Background())
	if err != nil {
		return
	}
	for key, value := range mappingMap {
		if key == indexName {
			res = value
		}
	}
	return
}

func (this_ *V7Service) PutMapping(indexName string, bodyJSON map[string]interface{}) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.PutMapping().Index(indexName).BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *V7Service) SetFieldType(indexName string, fieldName string, fieldType string) (err error) {
	bodyJSON := map[string]interface{}{}
	bodyJSON["properties"] = map[string]interface{}{
		fieldName: map[string]interface{}{
			"type": fieldType,
		},
	}
	err = this_.PutMapping(indexName, bodyJSON)
	if err != nil {
		return
	}
	return
}

type SearchResult struct {
	*elastic.SearchResult
}

func (this_ *V7Service) Search(indexName string, pageIndex int, pageSize int) (res *SearchResult, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Search(indexName)
	query := elastic.NewBoolQuery()
	searchResult, err := doer.Query(query).Size(pageSize).From((pageIndex - 1) * pageSize).Do(context.Background())
	if err != nil {
		return
	}
	res = &SearchResult{
		SearchResult: searchResult,
	}

	return
}

type IndexResponse struct {
	*elastic.IndexResponse
}

func (this_ *V7Service) Insert(indexName string, id string, doc interface{}) (res *IndexResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	doer := client.Index()
	indexResponse, err := doer.Index(indexName).Id(id).BodyJson(doc).Refresh("wait_for").Do(context.Background())
	if err != nil {
		return
	}
	res = &IndexResponse{
		IndexResponse: indexResponse,
	}
	return
}

type UpdateResponse struct {
	*elastic.UpdateResponse
}

func (this_ *V7Service) Update(indexName string, id string, doc interface{}) (res *UpdateResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Update()
	updateResponse, err := doer.Index(indexName).Id(id).Doc(doc).Refresh("wait_for").Do(context.Background())
	if err != nil {
		return
	}
	res = &UpdateResponse{
		UpdateResponse: updateResponse,
	}

	return
}

type DeleteResponse struct {
	*elastic.DeleteResponse
}

func (this_ *V7Service) Delete(indexName string, id string) (res *DeleteResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Delete()
	deleteResponse, err := doer.Index(indexName).Id(id).Refresh("wait_for").Do(context.Background())
	if err != nil {
		return
	}
	res = &DeleteResponse{
		DeleteResponse: deleteResponse,
	}

	return
}

type BulkIndexByScrollResponse struct {
	*elastic.BulkIndexByScrollResponse
}

func (this_ *V7Service) Reindex(sourceIndexName string, toIndexName string) (res *BulkIndexByScrollResponse, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Reindex()
	bulkIndexByScrollResponse, err := doer.Source(elastic.NewReindexSource().Index(sourceIndexName)).DestinationIndex(toIndexName).Refresh("true").Do(context.Background())
	if err != nil {
		return
	}
	res = &BulkIndexByScrollResponse{
		BulkIndexByScrollResponse: bulkIndexByScrollResponse,
	}

	return
}

func (this_ *V7Service) Scroll(indexName string, scrollId string, pageSize int) (res *SearchResult, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	doer := client.Scroll(indexName)
	query := elastic.NewBoolQuery()
	searchResult, err := doer.Query(query).Size(pageSize).ScrollId(scrollId).Do(context.Background())
	if err != nil {
		return
	}
	res = &SearchResult{
		SearchResult: searchResult,
	}

	return
}
