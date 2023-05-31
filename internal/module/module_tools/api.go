package module_tools

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/team-ide/go-tool/util"
	"io"
	"net/url"
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
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: base64Power, Do: this_.base64})
	apis = append(apis, &base.ApiWorker{Power: md5Power, Do: this_.md5})
	apis = append(apis, &base.ApiWorker{Power: urlEncode, Do: this_.urlEncode})
	apis = append(apis, &base.ApiWorker{Power: randomNumber, Do: this_.randomNumber})
	apis = append(apis, &base.ApiWorker{Power: randomString, Do: this_.randomString})

	return
}

type BaseRequest struct {
	Value       string `json:"value,omitempty"`
	ValueEncode string `json:"valueEncode,omitempty"`
	Decode      string `json:"decode,omitempty"`
	DecodeValue string `json:"decodeValue,omitempty"`
	Min         int64  `json:"min,omitempty"`
	Max         int64  `json:"max,omitempty"`
	MinLen      int    `json:"minLen,omitempty"`
	MaxLen      int    `json:"maxLen,omitempty"`
	RandString  string `json:"randString,omitempty"`
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
