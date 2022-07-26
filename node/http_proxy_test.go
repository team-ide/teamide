package node

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
)

//获取主机名

func getHost(req *http.Request) string {
	if req.Host != "" {
		if hostPart, _, err := net.SplitHostPort(req.Host); err == nil {
			return hostPart
		}
		return req.Host
	}
	return "localhost"

}

var (
	testClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{
					tls.TLS_RSA_WITH_RC4_128_SHA,
					tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
					tls.TLS_RSA_WITH_AES_128_CBC_SHA,
					tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				},
			},
		},
	}
	bindUrl   = "http://127.0.0.1:8081"
	bindHost  = "127.0.0.1:8081"
	proxyUrl  = "https://www.baidu.com"
	proxyHost = "www.baidu.com"
)

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	//host := getHost(req)
	println("请求域名：" + bindUrl + "，转到：" + proxyUrl)

	header := req.Header
	requestPath := req.RequestURI
	println("requestPath:", requestPath)
	println("proxyUrl:", proxyUrl+requestPath)
	toReq, err := http.NewRequest(req.Method, proxyUrl+requestPath, nil)
	if err != nil {
		panic(err)
	}
	//toReq.Header.Set(":authority", "192.168.6.160")
	//toReq.Header.Set(":method", req.Method)
	//toReq.Header.Set(":scheme", "https")
	//toReq.Header.Set(":path", req.RequestURI)
	cs := req.Cookies()
	for _, c := range cs {
		toReq.AddCookie(c)
	}
	for key, value := range header {
		for _, v := range value {
			if strings.ToLower(key) == "content-length" {
				continue
			}
			//if key == "Connection" {
			//	continue
			//}
			v = strings.ReplaceAll(v, bindUrl, proxyUrl)
			v = strings.ReplaceAll(v, bindHost, proxyHost)
			req.Header.Set(key, v)
			if requestPath == "/users/sign_in" {
				println("request header set key:", key, ",value:", v)
			}
		}
	}
	req.Header.Add("X-Forwarded-Host", proxyHost)

	_ = req.ParseMultipartForm(1024 * 1024 * 1024)
	if requestPath == "/users/sign_in" {
		println(req.Method)
		println(req.Form)
		println(req.PostForm)
		println(req.MultipartForm)
		println(req.Body)
	}
	if req.PostForm != nil {
		toReq.PostForm = req.PostForm
	}
	if req.Form != nil {
		toReq.Form = req.Form
	}
	toReq.MultipartForm = req.MultipartForm
	toReq.Body = req.Body
	toReq.Host = proxyHost

	req.RequestURI = ""
	req.URL = toReq.URL

	toRes, err := testClient.Do(toReq)

	if err != nil {
		panic(err)
	}
	defer func() {
		_ = toRes.Body.Close()
	}()
	resHeader := res.Header()

	for key, value := range toRes.Header {
		for _, v := range value {
			if strings.ToLower(key) == "content-length" {
				continue
			}
			v = strings.ReplaceAll(v, bindUrl, proxyUrl)
			v = strings.ReplaceAll(v, bindHost, proxyHost)
			resHeader.Set(key, v)
			if requestPath == "/users/sign_in" {
				println("response header set key:", key, ",value:", v)
			}
		}
	}
	res.WriteHeader(toRes.StatusCode)
	_, err = io.Copy(res, toRes.Body)

	if err != nil {
		panic(err)
	}
	//url2, _ := url.Parse(proxyUrl)
	//// create the reverse proxy
	//proxy := httputil.NewSingleHostReverseProxy(url2)
	//// Update the headers to allow for SSL redirection
	//req.URL.Host = url2.Host
	//req.URL.Scheme = url2.Scheme
	//req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	//req.Host = url2.Host
	//// Note that ServeHttp is non blocking and uses a go routine under the hood
	//proxy.ServeHTTP(res, req)

}

func TestHttpProxy(t *testing.T) {
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		println("Proxy监听80端口错误：" + err.Error())
		panic(err)
	}

	var waitGroupForStop sync.WaitGroup
	waitGroupForStop.Add(1)
	waitGroupForStop.Wait()
}
