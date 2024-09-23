package module_http

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/dop251/goja"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/task"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"teamide/internal/module/module_toolbox"
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
	ToolboxId   int64    `json:"toolboxId,omitempty"`
	ExtendId    int64    `json:"extendId,omitempty"`
	ExecuteId   string   `json:"executeId,omitempty"`
	Name        string   `json:"name,omitempty"`
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
	extend   *Extend
	dir      string

	runtime        *goja.Runtime
	scriptContext  map[string]interface{}
	lock           sync.Mutex
	toolboxService *module_toolbox.ToolboxService
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

func (this_ *Request) init() (err error) {
	if this_.runtime != nil {
		return
	}
	this_.lock.Lock()
	defer this_.lock.Unlock()
	if this_.runtime != nil {
		return
	}
	this_.runtime = goja.New()
	this_.scriptContext = javascript.NewContext()
	if len(this_.scriptContext) > 0 {
		for key, value := range this_.scriptContext {
			err = this_.runtime.Set(key, value)
			if err != nil {
				return
			}
		}
	}
	if this_.extend != nil {
		for _, one := range this_.extend.Variables {
			if !one.Selected || one.Key == "" {
				continue
			}
			err = this_.runtime.Set(one.Key, one.Value)
			if err != nil {
				return
			}
		}
		for _, one := range this_.extend.Secrets {
			if !one.Selected || one.Key == "" {
				continue
			}
			err = this_.runtime.Set(one.Key, this_.toolboxService.DecryptOptionAttr(one.Value))
			if err != nil {
				return
			}
		}
	}
	return
}

func (this_ *Request) scriptValue(script string, param *task.ExecutorParam) (res any, err error) {
	if script == "" {
		return
	}

	err = this_.init()
	if err != nil {
		return
	}

	this_.lock.Lock()
	defer this_.lock.Unlock()

	if param == nil {
		param = &task.ExecutorParam{}
	}

	err = this_.runtime.Set("index", param.Index)
	if err != nil {
		return
	}
	err = this_.runtime.Set("workerIndex", param.WorkerIndex)
	if err != nil {
		return
	}

	v, err := this_.runtime.RunString(script)
	if err != nil {
		err = errors.New("get scriptValue error:" + err.Error())
		return
	}
	if v != nil {
		res = v.Export()
	}

	return
}
func (this_ *Request) stringArg(arg string, param *task.ExecutorParam) (res any, err error) {
	if arg == "" {
		res = ""
		return
	}
	text := ""
	var re *regexp.Regexp
	re, _ = regexp.Compile(`[$]+{(.+?)}`)
	indexList := re.FindAllIndex([]byte(arg), -1)
	var lastIndex int = 0
	var v any
	var vS string
	for _, indexes := range indexList {
		text += arg[lastIndex:indexes[0]]

		lastIndex = indexes[1]

		script := arg[indexes[0]+2 : indexes[1]-1]
		script = strings.TrimSpace(script)
		v, err = this_.scriptValue(script, param)
		if err != nil {
			return
		}
		vS = util.GetStringValue(v)
		text += vS
		//fmt.Println("stringArg:"+arg, ",value:", v, ",valueS:", vS, ",text:", text)
	}
	text += arg[lastIndex:]

	res = text
	if v != nil && text == vS {
		return v, nil
	}
	return
}
func (this_ *Request) formatArg(arg interface{}, param *task.ExecutorParam) (res any, err error) {
	if arg == nil {
		return
	}
	switch tV := arg.(type) {
	case string:
		res, err = this_.stringArg(tV, param)
		break
	case []interface{}:
		var list []interface{}
		for _, one := range tV {
			var v interface{}
			v, err = this_.formatArg(one, param)
			if err != nil {
				return
			}
			list = append(list, v)
		}
		res = list
		break
	case map[string]any:
		var data = map[string]any{}
		for key, one := range tV {
			var v interface{}
			v, err = this_.formatArg(one, param)
			if err != nil {
				return
			}
			data[key] = v
		}
		res = data
		break
	default:
		res = tV
		break
	}

	return
}

func (this_ *Request) formatArgs(args []interface{}, param *task.ExecutorParam) (res []interface{}, err error) {
	if len(args) == 0 {
		return
	}
	for _, arg := range args {
		var v interface{}
		v, err = this_.formatArg(arg, param)
		if err != nil {
			return
		}
		res = append(res, v)
	}

	return
}

func (this_ *Request) formatValue(name, value string) (res string) {
	defer func() {
		if name != "" {
			_ = this_.init()
			if this_.runtime != nil {
				_ = this_.runtime.Set(name, res)
			}
		}
	}()
	v, err := this_.stringArg(value, nil)
	if err != nil {
		res = value
		this_.toolboxService.Logger.Error("formatValue error", zap.Any("value", value), zap.Error(err))
		return
	}
	res = util.GetStringValue(v)
	//fmt.Println("format:"+name, ",value:", res)

	return
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
			res = res + one.Key + "=" + this_.formatValue(one.Key, one.Value)
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
			res.Set(one.Key, this_.formatValue(one.Key, one.Value))
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
			res.Set(one.Key, this_.formatValue(one.Key, one.Value))
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
					v := this_.formatValue(one.Key, one.Value)
					if p != nil {
						_, _ = p.Write([]byte(v))
					}

					err = w.WriteField(one.Key, v)
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
	FileName      string      `json:"fileName,omitempty"`
	*ContentInfo
}

type Execute struct {
	StartTime    int64     `json:"startTime,omitempty"`
	EndTime      int64     `json:"endTime,omitempty"`
	Error        string    `json:"error,omitempty"`
	RequestTime  int64     `json:"requestTime,omitempty"`
	ResponseTime int64     `json:"responseTime,omitempty"`
	Request      *Request  `json:"request,omitempty"`
	Response     *Response `json:"response,omitempty"`
}

func (this_ *api) Execute(request *Request) (res *Execute, err error) {
	res = &Execute{
		Request: request,
	}
	res.StartTime = util.GetNowMilli()
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

	res.RequestTime = util.GetNowMilli()
	rR, err := client.Do(r)
	res.ResponseTime = util.GetNowMilli()
	if err != nil {
		return
	}

	defer func() {
		res.EndTime = util.GetNowMilli()
		if err != nil {
			res.Error = err.Error()
			err = nil
		}

		filePath := request.dir + "execute.json"
		var f *os.File
		f, err = os.Create(filePath)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		_, err = f.WriteString(util.GetStringValue(res))
		if err != nil {
			return
		}
	}()
	if e, _ := util.PathExists(request.dir); !e {
		_ = os.MkdirAll(request.dir, fs.ModePerm)
	}
	resp := &Response{}
	res.Response = resp
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
	resp.ContentInfo = GetContentInfo(resp.ContentType)
	// 是文件流
	if resp.ContentInfo.IsFile {
		//  获取Content-Disposition头
		contentDisposition := rR.Header.Get("Content-Disposition")
		if contentDisposition != "" && strings.Contains(contentDisposition, "filename=") {
			// 解析Content-Disposition头来提取文件名
			ss := strings.Split(contentDisposition, "filename=")
			if len(ss) > 1 {
				filename := strings.TrimSpace(ss[1])
				resp.FileName, _ = url.PathUnescape(filename)
				if resp.FileName == "" {
					resp.FileName = filename
				}
			}
		}
		if resp.FileName == "" {
			resp.FileName = GetFileNameFromURL(request.Url)
		}
		if resp.FileName == "" {
			resp.FileName = util.GetUUID()
		}

		filePath := request.dir + resp.FileName
		var f *os.File
		f, err = os.Create(filePath)
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
