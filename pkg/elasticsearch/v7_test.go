package elasticsearch

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttps(t *testing.T) {
	var err error
	var options []elastic.ClientOptionFunc

	var CertPath = "C:\\Users\\ZhuLiang\\Downloads\\http.crt"
	CertPath = ""
	options = append(options, elastic.SetURL("https://192.168.0.85:11110/"))
	options = append(options, elastic.SetSniff(false))
	if CertPath != "" {
		httpClient := &http.Client{}
		TLSClientConfig := &tls.Config{}
		certPool := x509.NewCertPool()
		var pemCerts []byte
		pemCerts, err = ioutil.ReadFile(CertPath)
		if err != nil {
			panic(err)
			return
		}

		if !certPool.AppendCertsFromPEM(pemCerts) {
			err = errors.New("证书[" + CertPath + "]解析失败")
			panic(err)
			return
		}
		TLSClientConfig.RootCAs = certPool
		TLSClientConfig.InsecureSkipVerify = true

		//TLSClientConfig.Certificates = []tls.Certificate{clicrt}

		httpClient.Transport = &http.Transport{
			TLSClientConfig: TLSClientConfig,
		}
		options = append(options, elastic.SetHttpClient(httpClient))
	}
	options = append(options, elastic.SetBasicAuth("elastic", "elastic"))

	var client *elastic.Client
	client, err = elastic.NewClient(options...)
	if err != nil {
		panic(err)
		return
	}

	res, err := client.IndexNames()
	if err != nil {
		panic(err)
		return
	}
	bs, _ := json.Marshal(res)
	println("IndexNames:", string(bs))
}
