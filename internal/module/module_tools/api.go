package module_tools

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"teamide/internal/context"
	"teamide/pkg/base"
)

type Api struct {
	*context.ServerContext
}

func NewApi(ServerContext *context.ServerContext) *Api {
	return &Api{
		ServerContext: ServerContext,
	}
}

var (

	// Power 基本 权限
	Power        = base.AppendPower(&base.PowerAction{Action: "tools", Text: "小工具", ShouldLogin: true, StandAlone: true})
	base64Power  = base.AppendPower(&base.PowerAction{Action: "base64", Text: "Base64", Parent: Power, ShouldLogin: true, StandAlone: true})
	md5Power     = base.AppendPower(&base.PowerAction{Action: "md5", Text: "MD5", Parent: Power, ShouldLogin: true, StandAlone: true})
	urlEncode    = base.AppendPower(&base.PowerAction{Action: "urlEncode", Text: "urlEncode", Parent: Power, ShouldLogin: true, StandAlone: true})
	randomNumber = base.AppendPower(&base.PowerAction{Action: "randomNumber", Text: "randomNumber", Parent: Power, ShouldLogin: true, StandAlone: true})
	randomString = base.AppendPower(&base.PowerAction{Action: "randomString", Text: "randomString", Parent: Power, ShouldLogin: true, StandAlone: true})
	fileSearch   = base.AppendPower(&base.PowerAction{Action: "fileSearch", Text: "fileSearch", Parent: Power, ShouldLogin: true, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: base64Power, Do: this_.base64})
	apis = append(apis, &base.ApiWorker{Power: md5Power, Do: this_.md5})
	apis = append(apis, &base.ApiWorker{Power: urlEncode, Do: this_.urlEncode})
	apis = append(apis, &base.ApiWorker{Power: randomNumber, Do: this_.randomNumber})
	apis = append(apis, &base.ApiWorker{Power: randomString, Do: this_.randomString})
	apis = append(apis, &base.ApiWorker{Power: fileSearch, Do: this_.fileSearch})

	return
}

type BaseRequest struct {
	Value             string  `json:"value,omitempty"`
	ValueEncode       string  `json:"valueEncode,omitempty"`
	Decode            string  `json:"decode,omitempty"`
	DecodeValue       string  `json:"decodeValue,omitempty"`
	Min               int64   `json:"min,omitempty"`
	Max               int64   `json:"max,omitempty"`
	MinLen            int     `json:"minLen,omitempty"`
	MaxLen            int     `json:"maxLen,omitempty"`
	RandString        string  `json:"randString,omitempty"`
	Path              string  `json:"path,omitempty"`
	SearchFile        string  `json:"searchFile,omitempty"`
	SearchFileMinSize float64 `json:"searchFileMinSize,omitempty"` // 搜索 文件大小 大于等于该值 单位 M
	SearchFileMaxSize float64 `json:"searchFileMaxSize,omitempty"` // 搜索 文件大小 小于等于该值 单位 M
	SearchContent     string  `json:"searchContent,omitempty"`
	RecursiveDir      bool    `json:"recursiveDir,omitempty"`    // 是否递归目录
	RecursiveLevel    int     `json:"recursiveLevel,omitempty"`  // 递归目录层级  0 无限制
	FileMaxReadSize   float64 `json:"fileMaxReadSize,omitempty"` // 读取的 最大文件内容 单位 M
}

func (this_ *Api) base64(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	response := &BaseRequest{}
	res = response
	if request.Value != "" {
		response.ValueEncode = base64.StdEncoding.EncodeToString([]byte(request.Value))
	} else if request.Decode != "" {
		var bs []byte
		bs, err = base64.StdEncoding.DecodeString(request.Decode)
		if err != nil {
			return
		}
		response.DecodeValue = string(bs)
	}

	return
}

func (this_ *Api) urlEncode(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	response := &BaseRequest{}
	res = response
	if request.Value != "" {
		response.ValueEncode = url.QueryEscape(request.Value)
	} else if request.Decode != "" {
		response.DecodeValue, err = url.QueryUnescape(request.Decode)
	}

	return
}

func (this_ *Api) md5(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Value != "" {
		m := md5.New()
		_, err = io.WriteString(m, request.Value)
		if err != nil {
			return
		}
		arr := m.Sum(nil)
		res = fmt.Sprintf("%x", arr)
	}

	return
}

func (this_ *Api) randomNumber(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Max >= request.Min {
		res = util.RandomInt64(request.Min, request.Max)
	}

	return
}

func (this_ *Api) randomString(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.MaxLen >= request.MinLen {

		size := request.MinLen
		if request.MaxLen > request.MinLen {
			size = util.RandomInt(request.MinLen, request.MaxLen)
		}

		ss := strings.Split(request.RandString, "")

		ssSize := len(ss)
		var str string
		for i := 0; i < size; i++ {
			randNum := 0
			randNum = util.RandomInt(0, ssSize*3)
			str += ss[randNum%ssSize]
		}
		res = str
	}

	return
}

func (this_ *Api) fileSearch(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Path == "" {
		err = errors.New("请输入搜索目录")
		return
	}
	searchInfo := &SearchInfo{
		BaseRequest: request,
	}
	if searchInfo.SearchFile != "" {
		searchInfo.fileNameMatchRegexp = regexp.MustCompile(searchInfo.SearchFile)
	}
	if searchInfo.SearchContent != "" {
		searchInfo.fileContentMatchRegexp = regexp.MustCompile(searchInfo.SearchContent)
	}
	searchInfo.Dir = util.FormatPath(searchInfo.Path)
	searchInfo.SearchFileMinSize = searchInfo.SearchFileMinSize * 1024 * 1024
	searchInfo.SearchFileMaxSize = searchInfo.SearchFileMaxSize * 1024 * 1024
	searchInfo.FileMaxReadSize = searchInfo.FileMaxReadSize * 1024 * 1024
	searchInfo.searchFileCache = make(map[string]bool)
	searchFile(searchInfo, request.Path, 0)
	res = searchInfo
	return
}

type SearchInfo struct {
	*BaseRequest
	searchFileCache        map[string]bool          // 已搜索
	FileList               []map[string]interface{} `json:"fileList"` // 搜索到的文件
	Dir                    string                   `json:"dir"`
	fileNameMatchRegexp    *regexp.Regexp
	fileContentMatchRegexp *regexp.Regexp
}

func searchFile(searchInfo *SearchInfo, dir string, dirLevel int) {

	dir = util.FormatPath(dir)
	fs, e := os.ReadDir(dir)
	if e != nil {
		return
	}
	for _, f := range fs {
		path := dir + "/" + f.Name()
		if searchInfo.searchFileCache[path] {
			continue
		}
		searchInfo.searchFileCache[path] = true
		if f.IsDir() {
			if searchInfo.RecursiveDir {
				if searchInfo.RecursiveLevel > 0 && dirLevel >= searchInfo.RecursiveLevel {
					continue
				}
				searchFile(searchInfo, path, dirLevel+1)
			}
			continue
		}
		info, _ := f.Info()
		if info == nil {
			continue
		}
		if searchInfo.fileNameMatchRegexp != nil {
			if !searchInfo.fileNameMatchRegexp.MatchString(f.Name()) {
				continue
			}
		}
		fileSize := info.Size()
		if searchInfo.SearchFileMinSize > 0 && float64(fileSize) < searchInfo.SearchFileMinSize {
			continue
		}
		if searchInfo.SearchFileMaxSize > 0 && float64(fileSize) > searchInfo.SearchFileMaxSize {
			continue
		}
		if searchInfo.fileContentMatchRegexp != nil {
			if searchInfo.FileMaxReadSize > 0 && float64(fileSize) > searchInfo.FileMaxReadSize {
				continue
			}
			// 匹配 文件 内容
			if !searchContentInFile(path, searchInfo.SearchContent, searchInfo.fileContentMatchRegexp) {
				continue
			}
		}
		data := map[string]interface{}{}
		data["path"] = path
		data["name"] = f.Name()
		data["size"] = info.Size()
		data["dirLevel"] = dirLevel
		data["modTime"] = info.ModTime().UnixMilli()

		searchInfo.FileList = append(searchInfo.FileList, data)
	}
}
func searchContentInFile(filePath string, contentMatch string, contentMatchRegexp *regexp.Regexp) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, contentMatch) || contentMatchRegexp.MatchString(scanner.Text()) {
			return true
		}
	}

	if err = scanner.Err(); err != nil {
		return false
	}

	return false
}
