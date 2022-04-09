package toolbox

import (
	"context"

	elastic "github.com/olivere/elastic/v7"
)

func getESService(url string) (res *ESService, err error) {
	key := "elasticsearch-" + url
	var service Service
	service, err = GetService(key, func() (res Service, err error) {
		var s *ESService
		s, err = CreateESService(url)
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
	res = service.(*ESService)
	return
}

func CreateESService(url string) (*ESService, error) {
	service := &ESService{
		url: url,
	}
	err := service.init()
	return service, err
}

//ESService 注册处理器在线信息等
type ESService struct {
	url         string
	lastUseTime int64
}

func (this_ *ESService) init() error {
	var err error
	return err
}
func (this_ *ESService) GetClient() (client *elastic.Client, err error) {
	defer func() {
		this_.lastUseTime = GetNowTime()
	}()
	client, err = elastic.NewClient(
		elastic.SetURL(this_.url),
		//docker
		elastic.SetSniff(false),
	)
	return
}

func (this_ *ESService) GetWaitTime() int64 {
	return 10 * 60 * 1000
}

func (this_ *ESService) GetLastUseTime() int64 {
	return this_.lastUseTime
}

func (this_ *ESService) Stop() {
}

func (this_ *ESService) DeleteIndex(indexName string) (err error) {
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

func (this_ *ESService) CreateIndex(indexName string) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.CreateIndex(indexName).Do(context.Background())
	if err != nil {
		return
	}
	return
}

func (this_ *ESService) IndexNames() (res []string, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	res, err = client.IndexNames()
	if err != nil {
		return
	}
	return
}

func (this_ *ESService) GetMapping(indexName string) (res map[string]interface{}, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	res, err = client.GetMapping().Index(indexName).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return
}

func (this_ *ESService) PutMapping(indexName string, bodyJSON map[string]interface{}) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()
	_, err = client.PutMapping().Index(indexName).BodyJson(bodyJSON).Do(context.Background())
	if err != nil {
		panic(err)
	}
	return
}

func (this_ *ESService) SetFieldType(indexName string, fieldName string, fieldType string) (err error) {
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

func (this_ *ESService) Search(indexName string, pageIndex int, pageSize int) (res *ESQueryResult, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	search := client.Search(indexName)
	query := elastic.NewBoolQuery()
	searchResult, err := search.Query(query).Size(pageSize).From((pageIndex - 1) * pageSize).Do(context.Background())
	if err != nil {
		return
	}
	res = &ESQueryResult{}

	res.Count = searchResult.TotalHits()

	return
}

func (this_ *ESService) Scroll(indexName string, scrollId string, pageSize int) (res *ESQueryResult, err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	defer client.Stop()

	search := client.Scroll(indexName)
	query := elastic.NewBoolQuery()
	searchResult, err := search.Query(query).Size(pageSize).ScrollId(scrollId).Do(context.Background())
	if err != nil {
		return
	}
	res = &ESQueryResult{}

	res.Count = searchResult.TotalHits()
	return
}
