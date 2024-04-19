package base

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type RequestBean struct {
	JWT          *JWTBean
	Path         string
	Action       string
	ClientKey    string
	ClientTabKey string
	extends      map[string]interface{}
	extendsLock  sync.Mutex
	Power        *PowerAction
}

func (this_ *RequestBean) GetExtend(key string) interface{} {
	if this_.extends == nil {
		return nil
	}

	this_.extendsLock.Lock()
	defer this_.extendsLock.Unlock()
	return this_.extends[key]
}

func (this_ *RequestBean) SetExtend(key string, value interface{}) {
	this_.extendsLock.Lock()
	defer this_.extendsLock.Unlock()

	if this_.extends == nil {
		this_.extends = map[string]interface{}{}
	}

	this_.extends[key] = value
}

var (
	HttpNotResponse = &HttpResponse{}
)

type JWTBean struct {
	Sign        string `json:"sign,omitempty"`
	UserId      int64  `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	Account     string `json:"account,omitempty"`
	Time        int64  `json:"time,omitempty"`
	LoginId     int64  `json:"loginId,omitempty"`
	IsAnonymous bool   `json:"isAnonymous,omitempty"`
}

type HttpResponse struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func RequestJSON(data interface{}, c *gin.Context) bool {
	err := c.ShouldBindBodyWith(data, binding.JSON)
	if err != nil {
		ResponseJSON(nil, err, c)
		return false
	}
	return true
}

func ResponseJSON(data interface{}, err error, c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			util.Logger.Error("resource json error", zap.Any("error", e))
		}
	}()
	if _, exists := c.Get("request-data-bind-json-error"); exists {
		return
	}
	response := HttpResponse{
		Code: "0",
		Data: data,
	}
	if err != nil {
		response.Msg = err.Error()
		baseErr := ToBaseError(err)
		if baseErr != nil {
			response.Code = baseErr.Code
		} else {
			response.Code = "-1"
		}
	} else {
		response.Data = data
	}
	c.Set("request-data-bind-json-error", 1)
	c.JSON(http.StatusOK, response)
}
