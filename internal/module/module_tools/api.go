package module_tools

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	Power       = base.AppendPower(&base.PowerAction{Action: "tools", Text: "小工具", ShouldLogin: true, StandAlone: true})
	base64Power = base.AppendPower(&base.PowerAction{Action: "base64", Text: "Base64", Parent: Power, ShouldLogin: true, StandAlone: true})
	md5Power    = base.AppendPower(&base.PowerAction{Action: "md5", Text: "MD5", Parent: Power, ShouldLogin: true, StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: base64Power, Do: this_.base64})
	apis = append(apis, &base.ApiWorker{Power: md5Power, Do: this_.md5})

	return
}

type BaseRequest struct {
	Value  string `json:"value,omitempty"`
	Base64 string `json:"base64,omitempty"`
}

func (this_ *Api) base64(_ *base.RequestBean, c *gin.Context) (res interface{}, err error) {

	request := &BaseRequest{}
	if !base.RequestJSON(request, c) {
		return
	}

	if request.Value != "" {
		res = base64.StdEncoding.EncodeToString([]byte(request.Value))
	} else if request.Base64 != "" {
		var bs []byte
		bs, err = base64.StdEncoding.DecodeString(request.Base64)
		if err != nil {
			return
		}
		res = string(bs)
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
