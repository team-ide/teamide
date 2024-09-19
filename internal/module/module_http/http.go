package module_http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"github.com/team-ide/go-tool/util"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	MaxIdleConn        int  `json:"maxIdleConn,omitempty"`     // 控制空闲（保持活动）的最大数量
	MaxConnPerHost     int  `json:"maxConnPerHost,omitempty"`  // 每个主机的连接（包括拨号中的连接），包括活动状态和空闲状态。违反限制时，拨号将被阻止
	IdleConnTimeout    int  `json:"idleConnTimeout,omitempty"` // 空闲的最长时间 秒
	Timeout            int  `json:"timeout,omitempty"`         // 超时时间 秒
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
}

func NewClient(config *Config) (client *http.Client) {
	if config == nil {
		config = &Config{}
		config.MaxIdleConn = 0
		config.MaxConnPerHost = 0
		config.IdleConnTimeout = 0
		config.InsecureSkipVerify = true
		config.Timeout = 30
	}
	// 创建传输对象
	transport := &http.Transport{
		MaxIdleConns:    config.MaxIdleConn,
		MaxConnsPerHost: config.MaxConnPerHost,
		IdleConnTimeout: time.Second * time.Duration(config.IdleConnTimeout),
		TLSClientConfig: &tls.Config{
			// 指定不校验 SSL/TLS 证书
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}
	// 创建 HTTP 客户端
	client = &http.Client{
		Transport: transport,
		Timeout:   time.Second * time.Duration(config.Timeout),
	}
	return
}

type Request struct {
	*Config
	Username    string              `json:"username,omitempty"`
	Password    string              `json:"password,omitempty"`
	Url         string              `json:"url,omitempty"`
	Path        string              `json:"path,omitempty"`
	Params      map[string]string   `json:"params,omitempty"`
	Method      string              `json:"method,omitempty"`
	ContentType string              `json:"contentType,omitempty"`
	Header      http.Header         `json:"header,omitempty"`
	Body        string              `json:"body,omitempty"`
	Data        map[string]any      `json:"data,omitempty"`
	Files       map[string][]string `json:"files,omitempty"`
	IsForm      bool                `json:"isForm,omitempty"`
	IsJson      bool                `json:"isJson,omitempty"`
	IsXml       bool                `json:"isXml,omitempty"`
}

func (this_ *Request) GetUrl() string {
	res := this_.Url
	if this_.Path != "" {
		if res != "" && !strings.HasSuffix(res, "/") && !strings.HasPrefix(this_.Path, "/") {
			res += "/" + this_.Path
		} else {
			res += this_.Path
		}
	}
	if this_.Params != nil {
		for key, value := range this_.Params {
			if strings.Contains(res, "?") {
				res += "&"
			} else {
				res += "?"
			}
			res = res + key + "=" + value
		}
	}
	return res
}

func (this_ *Request) FormBody() url.Values {
	res := url.Values{}
	if this_.Data != nil {
		for key, value := range this_.Data {
			res.Set(key, util.GetStringValue(value))
		}
	}
	return res
}

func (this_ *Request) GetHeader() http.Header {
	res := this_.Header
	if res == nil {
		res = http.Header{}
	}
	res.Set("Content-Type", this_.ContentType)
	return res
}

func (this_ *Request) BodyReader() (res io.Reader, err error) {
	if this_.Files != nil && len(this_.Files) > 0 {
		// 创建一个multipart表单
		var b = &bytes.Buffer{}
		w := multipart.NewWriter(b)
		for name, fs := range this_.Files {
			for _, f := range fs {
				err = this_.appendFile(w, name, f)
				if err != nil {
					return
				}
			}
		}
		// 关闭multipart writer，它将写入结束边界
		err = w.Close()
		if err != nil {
			return
		}
		if this_.Data != nil {
			for key, value := range this_.Data {
				err = w.WriteField(key, util.GetStringValue(value))
				if err != nil {
					return
				}
			}
		}
		if this_.ContentType == "" {
			this_.ContentType = w.FormDataContentType()
		}
		return b, err
	} else if this_.IsForm {
		res = strings.NewReader(this_.FormBody().Encode())
	} else if this_.IsJson {
		if this_.Data != nil {
			var bs []byte
			bs, err = json.Marshal(this_.Data)
			if err != nil {
				return
			}
			res = strings.NewReader(string(bs))
		} else {
			res = strings.NewReader(this_.Body)
		}
	} else if this_.IsXml {
		if this_.Data != nil {
			var bs []byte
			bs, err = xml.Marshal(this_.Data)
			if err != nil {
				return
			}
			res = strings.NewReader(string(bs))
		} else {
			res = strings.NewReader(this_.Body)
		}
	}
	return
}

func (this_ *Request) appendFile(w *multipart.Writer, name string, filePath string) (err error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	// 添加文件字段
	fw, err := w.CreateFormFile(name, file.Name())
	if err != nil {
		return
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return
	}
	return
}

type Response struct {
	ContentType   string              `json:"contentType,omitempty"`
	Header        http.Header         `json:"header,omitempty"`
	Body          string              `json:"body,omitempty"`
	Data          map[string]any      `json:"data,omitempty"`
	Files         map[string][]string `json:"files,omitempty"`
	IsForm        bool                `json:"isForm,omitempty"`
	IsJson        bool                `json:"isJson,omitempty"`
	IsXml         bool                `json:"isXml,omitempty"`
	Status        string              `json:"status,omitempty"`
	StatusCode    int                 `json:"statusCode,omitempty"`
	ContentLength int64               `json:"contentLength,omitempty"`
}

func (this_ *api) Execute(request *Request) (res *Response, err error) {
	reader, err := request.BodyReader()
	if err != nil {
		return
	}
	client := NewClient(request.Config)

	r, err := http.NewRequest(request.Method, request.GetUrl(), reader)
	if err != nil {
		return
	}
	r.Header = request.GetHeader()

	rR, err := client.Do(r)
	if err != nil {
		return
	}
	res = &Response{}
	res.Status = rR.Status
	res.StatusCode = rR.StatusCode
	res.ContentLength = rR.ContentLength
	res.Header = rR.Header
	res.ContentType = rR.Header.Get("Content-Type")
	resBody := rR.Body
	if resBody == nil {
		return
	}
	defer func() { _ = resBody.Close() }()
	var bs []byte
	tempDir, err := util.GetTempDir()
	if err != nil {
		return
	}
	// 是文件流
	if strings.Contains(res.ContentType, "application/octet-stream") {
		fileName := util.GetUUID()

		teamFilePath := tempDir + "/" + fileName
		var f *os.File
		f, err = os.Create(teamFilePath)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		_, err = io.Copy(f, resBody)
		if err != nil {
			return
		}
	} else {
		bs, err = io.ReadAll(resBody)
		if err != nil {
			return
		}
		res.Body = string(bs)
	}
	return
}
