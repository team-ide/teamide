package elasticsearch

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"teamide/pkg/util"
)

type Config struct {
	Url      string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	CertPath string `json:"certPath,omitempty"`
}

func CreateESService(config Config) (*V7Service, error) {
	service := &V7Service{
		Url:      config.Url,
		Username: config.Username,
		Password: config.Password,
		CertPath: config.CertPath,
	}
	err := service.Init()
	return service, err
}

//V7Service 注册处理器在线信息等
type V7Service struct {
	Url         string
	Username    string
	Password    string
	CertPath    string
	lastUseTime int64
	client      *elastic.Client
	clientLock  sync.Mutex
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
	this_.clientLock.Lock()
	defer this_.clientLock.Unlock()
	if this_.client != nil && this_.client.IsRunning() {
		client = this_.client
		return
	}
	var urls []string
	if strings.Contains(this_.Url, ",") {
		urls = strings.Split(this_.Url, ",")
	} else if strings.Contains(this_.Url, ";") {
		urls = strings.Split(this_.Url, ";")
	} else {
		urls = []string{this_.Url}
	}
	var isHttps bool
	for _, one := range urls {
		if strings.HasPrefix(one, "https") {
			isHttps = true
		}
	}

	var options []elastic.ClientOptionFunc

	options = append(options, elastic.SetURL(urls...))
	options = append(options, elastic.SetSniff(false))
	if isHttps {
		httpClient := &http.Client{}
		TLSClientConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		if this_.CertPath != "" {
			certPool := x509.NewCertPool()
			var pemCerts []byte
			pemCerts, err = ioutil.ReadFile(this_.CertPath)
			if err != nil {
				return
			}

			if !certPool.AppendCertsFromPEM(pemCerts) {
				err = errors.New("证书[" + this_.CertPath + "]解析失败")
				return
			}
			TLSClientConfig.RootCAs = certPool

			//TLSClientConfig.Certificates = []tls.Certificate{clicrt}

		}
		httpClient.Transport = &http.Transport{
			TLSClientConfig: TLSClientConfig,
		}
		options = append(options, elastic.SetHttpClient(httpClient))
	}
	if this_.Username != "" {
		options = append(options, elastic.SetBasicAuth(this_.Username, this_.Password))
	}
	client, err = elastic.NewClient(options...)
	if err != nil {
		return
	}
	this_.client = client
	return
}

func (this_ *V7Service) Stop() {
	if this_.client != nil {
		this_.client.Stop()
		this_.client = nil
	}
}

func (this_ *V7Service) DeleteIndex(indexName string) (err error) {
	client, err := this_.GetClient()
	if err != nil {
		return
	}
	//defer client.Stop()
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
	//defer client.Stop()
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
	//defer client.Stop()
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
	//defer client.Stop()
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
	//defer client.Stop()
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
	//defer client.Stop()

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
	//defer client.Stop()
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
	//defer client.Stop()

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
	//defer client.Stop()

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
	//defer client.Stop()

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
	//defer client.Stop()

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
