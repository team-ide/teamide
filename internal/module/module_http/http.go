package module_http

import (
	"bytes"
	"crypto/tls"
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
	Username    string   `json:"username,omitempty"`
	Password    string   `json:"password,omitempty"`
	Url         string   `json:"url,omitempty"`
	Path        string   `json:"path,omitempty"`
	Params      []*Field `json:"params,omitempty"`
	Method      string   `json:"method,omitempty"`
	ContentType string   `json:"contentType,omitempty"`
	Headers     []*Field `json:"headers,omitempty"`

	Body     string `json:"body,omitempty"`
	BodyType string `json:"bodyType,omitempty"`

	Text     string `json:"text,omitempty"`
	TextType string `json:"textType,omitempty"`

	FormData []*Field     `json:"formData,omitempty"`
	Files    []*FieldFile `json:"files,omitempty"`
}

type Field struct {
	Key      string       `json:"key,omitempty"`
	Value    string       `json:"value,omitempty"`
	IsFile   bool         `json:"isFile,omitempty"`
	Files    []*FieldFile `json:"files,omitempty"`
	Selected bool         `json:"selected,omitempty"`
}

type FieldFile struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
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
		for _, one := range this_.Params {
			if !one.Selected || one.Key == "" {
				continue
			}
			if strings.Contains(res, "?") {
				res += "&"
			} else {
				res += "?"
			}
			res = res + one.Key + "=" + one.Value
		}
	}
	return res
}

func (this_ *Request) FormBody() url.Values {
	res := url.Values{}
	if this_.FormData != nil {
		for _, one := range this_.FormData {
			if !one.Selected || one.Key == "" {
				continue
			}
			res.Set(one.Key, util.GetStringValue(one.Value))
		}
	}
	return res
}

func (this_ *Request) GetHeader() http.Header {
	res := http.Header{}
	if this_.Headers != nil {
		for _, one := range this_.Headers {
			if !one.Selected || one.Key == "" {
				continue
			}
			res.Set(one.Key, util.GetStringValue(one.Value))
		}
	}
	if this_.ContentType != "" {
		res.Set("Content-Type", this_.ContentType)
	}
	return res
}

func (this_ *Request) BodyReader() (res io.Reader, err error) {
	if this_.BodyType == "form" {
		if this_.FormData != nil {
			var hasFile bool
			for _, one := range this_.FormData {
				if !one.Selected || one.Key == "" || !one.IsFile {
					continue
				}
				hasFile = true
				break
			}
			if hasFile {
				// 创建一个multipart表单
				var b = &bytes.Buffer{}
				w := multipart.NewWriter(b)
				for _, one := range this_.FormData {
					if !one.Selected || one.Key == "" || !one.IsFile {
						continue
					}
					for _, f := range one.Files {
						this_.Body += "form field " + one.Key + "  binary:" + f.Name + "\n"
						err = this_.appendFormFile(w, one.Key, f)
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
				var wb = &bytes.Buffer{}
				ws := multipart.NewWriter(wb)
				for _, one := range this_.FormData {
					if !one.Selected || one.Key == "" || one.IsFile {
						continue
					}
					p, _ := ws.CreateFormField(one.Key)
					if p != nil {
						_, _ = p.Write([]byte(one.Value))
					}

					err = w.WriteField(one.Key, one.Value)
					if err != nil {
						return
					}
				}
				_ = ws.Close()
				bs, _ := io.ReadAll(wb)
				this_.Body += string(bs)
				if this_.ContentType == "" {
					this_.ContentType = w.FormDataContentType()
				}
				res = b
				return
			}
		}
		this_.Body = this_.FormBody().Encode()
		res = strings.NewReader(this_.Body)
		if this_.ContentType == "" {
			this_.ContentType = "application/x-www-form-urlencoded"
		}
		return
	}

	if this_.BodyType == "text" {
		this_.Body = this_.Text
		res = strings.NewReader(this_.Body)
		if this_.ContentType == "" {
			if this_.TextType == "json" {
				this_.ContentType = "application/json"
			}
			if this_.TextType == "xml" {
				this_.ContentType = "application/xml"
			}
		}
		return
	}

	if this_.BodyType == "binary" {
		if this_.ContentType == "" {
			this_.ContentType = "application/octet-stream"
		}
		// 创建一个multipart表单
		var b = &bytes.Buffer{}
		for _, one := range this_.Files {
			this_.Body += "file binary:" + one.Name + "\n"
			err = this_.appendFile(b, one)
			if err != nil {
				return
			}
		}
		res = b
		return
	}
	return
}

func (this_ *Request) appendFile(b *bytes.Buffer, f *FieldFile) (err error) {
	// 打开文件
	file, err := os.Open(f.Path)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	_, err = io.Copy(b, file)
	if err != nil {
		return
	}
	return
}

func (this_ *Request) appendFormFile(w *multipart.Writer, name string, f *FieldFile) (err error) {
	// 打开文件
	file, err := os.Open(f.Path)
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()
	// 添加文件字段
	fw, err := w.CreateFormFile(name, f.Name)
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
	ContentType   string      `json:"contentType,omitempty"`
	Header        http.Header `json:"header,omitempty"`
	Body          string      `json:"body,omitempty"`
	Status        string      `json:"status,omitempty"`
	StatusCode    int         `json:"statusCode,omitempty"`
	ContentLength int64       `json:"contentLength,omitempty"`
}

type Result struct {
	Request  *Request  `json:"request,omitempty"`
	Response *Response `json:"response,omitempty"`
}

func (this_ *api) Execute(request *Request) (res *Result, err error) {
	reader, err := request.BodyReader()
	if err != nil {
		return
	}
	client := NewClient(request.Config)

	request.Url = request.GetUrl()
	r, err := http.NewRequest(request.Method, request.Url, reader)
	if err != nil {
		return
	}
	r.Header = request.GetHeader()

	if request.ContentType != "" && r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", request.ContentType)
	}

	rR, err := client.Do(r)
	if err != nil {
		return
	}
	resp := &Response{}
	res = &Result{
		Request:  request,
		Response: resp,
	}
	resp.Status = rR.Status
	resp.StatusCode = rR.StatusCode
	resp.ContentLength = rR.ContentLength
	resp.Header = rR.Header
	resp.ContentType = rR.Header.Get("Content-Type")
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
	if strings.Contains(resp.ContentType, "application/octet-stream") {
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
		resp.Body = string(bs)
	}
	return
}
