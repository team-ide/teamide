package node

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	projectDir, _              = os.Getwd()
	fileName                   = projectDir + "/domain.config"
	readFileTime    int64      = 0 //读取文件的时间
	fileChangedTime int64      = 0 //文件修改时间
	domainData      [][]string     //[{***.gq,8080,http://127.0.0.1:8080/}]
	duPeiZhiSuo     sync.Mutex     //读配置锁
)

// 获取反向代理域名
func getProxyUrl(reqDomain string) string {
	checkFile()
	for _, dms := range domainData {
		if strings.Index(reqDomain, dms[0]) >= 0 {
			return dms[2]
		}
	}
	return domainData[0][2]
}

//读取配置文件

//域名:端口号，未知域名默认用第一个

func checkFile() {
	nowTime := time.Now().Unix()
	if nowTime-readFileTime < 300 {
		return
	}
	//每5分钟判断文件是否修改
	domainFile, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0)
	info, _ := domainFile.Stat()
	if info.ModTime().Unix() == fileChangedTime {
		return
	}
	duPeiZhiSuo.Lock()
	defer duPeiZhiSuo.Unlock()
	domainFile, _ = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0) //加锁再来一遍，防止重入
	info, _ = domainFile.Stat()
	changedTime := info.ModTime().Unix()
	if changedTime == fileChangedTime {
		return
	}

	//文件改变
	//重置数据
	readFileTime = nowTime
	fileChangedTime = changedTime
	domainData = [][]string{}
	bytes, _ := ioutil.ReadFile(fileName)
	split := strings.Split(string(bytes), "\n")
	for _, domainInfo := range split {
		dLen := len(domainInfo)
		if dLen < 8 || dLen > 20 { //忽略错误信息
			continue
		}
		domainItems := strings.Split(domainInfo, ":")
		if len(domainItems) != 2 || len(domainItems[0]) < 3 || len(domainItems[1]) < 2 {
			continue
		}

		if strings.HasSuffix(domainItems[1], "/") {
			domainItems = append(domainItems, "http://127.0.0.1:"+domainItems[1])
		} else {
			domainItems = append(domainItems, "http://127.0.0.1:"+domainItems[1]+"/")

		}
		domainData = append(domainData, domainItems)
	}
	domainSt, _ := json.Marshal(domainData)

	println("配置已修改：" + string(domainSt))

}

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

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	host := getHost(req)
	proxyUrl := getProxyUrl(host)
	url2, _ := url.Parse(proxyUrl)
	println("请求域名：" + host + "，转到：" + proxyUrl)
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url2)
	// Update the headers to allow for SSL redirection
	req.URL.Host = url2.Host
	req.URL.Scheme = url2.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url2.Host
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)

}

func TestHttpProxy(t *testing.T) {
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(":80", nil); err != nil {
		println("Proxy监听80端口错误：" + err.Error())
		panic(err)
	}

	var waitGroupForStop sync.WaitGroup
	waitGroupForStop.Add(1)
	waitGroupForStop.Wait()
}
